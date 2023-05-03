package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type User struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Language string
	Picture  string
	LineID   string `gorm:"column:line_id"`
}

type Restaurant struct {
	ID      int `gorm:"primaryKey"`
	Name    string
	Phone   string
	Address string
}

type UserRestaurant struct {
	ID           int `gorm:"primaryKey"`
	UserID       int
	RestaurantID int
	CreatedAt    time.Time `gorm:"default:now()"`
}

var (
	dbController *gorm.DB = nil
	dbUrl        string
)

func InitDB() error {
	userName := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbDomain := os.Getenv("DB_DOMAIN")
	dbName := os.Getenv("DB_NAME")

	// check is dev
	if os.Getenv("ISPROD") == "" {
		dbUrl = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", userName, password, dbDomain, dbName)
	} else {
		dbUrl = fmt.Sprintf("postgres://%s:%s@%s/%s", userName, password, dbDomain, dbName)
	}

	sqlDB, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}

	// 檢查連接是否正常
	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	dbController = gormDB

	initTable(&User{})
	initTable(&Restaurant{})
	initTable(&UserRestaurant{})

	return nil
}

func InitUserInDb(userId string, bot *linebot.Client) error {
	if IsUserExists(userId) {
		return nil
	}

	userData, err := bot.GetProfile(userId).Do()
	if err != nil {
		return err
	}

	return CreateUser(
		userData.DisplayName,
		userData.Language,
		userData.PictureURL,
		userData.UserID,
	)
}

// 初始化table，如果沒有table的話就創建
func initTable[T any](schema T) {
	if !dbController.Migrator().HasTable(schema) {
		if err := dbController.Migrator().CreateTable(schema); err != nil {
			log.Fatal(err)
		}
	}
}

func CreateUser(name string, language string, picture string, userId string) error {
	user := User{
		Name:     name,
		Language: language,
		Picture:  picture,
		LineID:   userId,
	}
	result := dbController.Create(&user)
	return result.Error
}

func GetUserByLineID(lineID string) (User, error) {
	var user User
	result := dbController.Where("line_id = ?", lineID).First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func IsUserExists(lineID string) bool {
	var count int64
	dbController.
		Model(&User{}).
		Where("line_id = ?", lineID).
		Count(&count)
	return count > 0
}

func IsRestaurantSaved(lineID string, restaurantName string) bool {
	var count int64
	dbController.
		Model(&User{}).
		Joins("JOIN user_restaurants ON users.id = user_restaurants.user_id").
		Joins("JOIN restaurants ON user_restaurants.restaurant_id = restaurants.id").
		Where("users.line_id = ? AND restaurants.name = ?", lineID, restaurantName).
		Count(&count)
	return count > 0
}

func GetRestaurantListByUser(lineID string) ([]Restaurant, error) {
	var restaurants []Restaurant
	result := dbController.Table("restaurants").Select("restaurants.*").Joins("inner join user_restaurants on user_restaurants.restaurant_id = restaurants.id").Joins("inner join users on users.id = user_restaurants.user_id").Where("users.line_id = ?", lineID).Find(&restaurants)
	return restaurants, result.Error
}

func CreateRestaurant(name string, phone string, address string) (Restaurant, error) {
	restaurant := Restaurant{
		Name:    name,
		Phone:   phone,
		Address: "",
	}
	result := dbController.Create(&restaurant)
	if result.Error != nil {
		return Restaurant{}, result.Error
	}
	return restaurant, nil
}

func AddRestaurantToUser(lineID string, restaurantID int) error {
	user, err := GetUserByLineID(lineID)
	if err != nil {
		return err
	}
	userRestaurant := UserRestaurant{
		UserID:       user.ID,
		RestaurantID: restaurantID,
	}
	result := dbController.Create(&userRestaurant)
	return result.Error
}

func GetRestaurantByName(restaurantName string) (Restaurant, error) {
	var restaurant Restaurant
	result := dbController.Table("restaurants").
		Where("restaurants.name = ?", restaurantName).
		First(&restaurant)
	return restaurant, result.Error
}

func RemoveRestaurantFromUser(lineID string, restaurantName string) error {
	user, err := GetUserByLineID(lineID)
	if err != nil {
		return err
	}

	restaurant, err := GetRestaurantByName(restaurantName)
	if err != nil {
		return errors.New("GetRestaurantByName error!!")
	}

	result := dbController.Table("user_restaurants").
		Where("user_restaurants.user_id = ? AND user_restaurants.restaurant_id = ?", user.ID, restaurant.ID).
		Delete(&UserRestaurant{})
	return result.Error
}

func PickRestaurantFromUser(lineID string) (Restaurant, error) {
	var restaurant Restaurant
	result := dbController.Table("restaurants").Select("restaurants.*").Joins("inner join user_restaurants on user_restaurants.restaurant_id = restaurants.id").Joins("inner join users on users.id = user_restaurants.user_id").Where("users.line_id = ?", lineID).Order("RANDOM()").Limit(1).First(&restaurant)
	return restaurant, result.Error
}

func IsUserRestaurantEmpty(lineID string) bool {
	user, err := GetUserByLineID(lineID)
	if err != nil {
		return true
	}
	var count int64
	dbController.Model(&UserRestaurant{}).Where("user_id = ?", user.ID).Count(&count)
	return count == 0
}
