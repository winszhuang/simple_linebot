package store

import (
	"database/sql"
	"errors"
	"fmt"
	"linebot/model"
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

type DBStore struct {
	gorm *gorm.DB
}

var (
	dbStore *DBStore
	dbUrl   string
)

func InitDB() {
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
		log.Fatal(err)
	}

	// 檢查連接是否正常
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	dbStore = &DBStore{
		gorm: gormDB,
	}

	initTable(&User{})
	initTable(&Restaurant{})
	initTable(&UserRestaurant{})
}

func GetDBStore() *DBStore {
	return dbStore
}

func (ds *DBStore) InitUserInDb(userId string, bot *linebot.Client) error {
	if ds.IsUserExists(userId) {
		return nil
	}

	userData, err := bot.GetProfile(userId).Do()
	if err != nil {
		return err
	}

	return ds.CreateUser(
		userData.DisplayName,
		userData.Language,
		userData.PictureURL,
		userData.UserID,
	)
}

// 初始化table，如果沒有table的話就創建
func initTable[T any](schema T) {
	if !dbStore.gorm.Migrator().HasTable(schema) {
		if err := dbStore.gorm.Migrator().CreateTable(schema); err != nil {
			log.Fatal(err)
		}
	}
}

func (ds *DBStore) CreateUser(name string, language string, picture string, userId string) error {
	user := User{
		Name:     name,
		Language: language,
		Picture:  picture,
		LineID:   userId,
	}
	result := ds.gorm.Create(&user)
	return result.Error
}

func (ds *DBStore) GetUserByLineID(lineID string) (User, error) {
	var user User
	result := ds.gorm.Where("line_id = ?", lineID).First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

func (ds *DBStore) IsUserExists(lineID string) bool {
	var count int64
	ds.gorm.
		Model(&User{}).
		Where("line_id = ?", lineID).
		Count(&count)
	return count > 0
}

func (ds *DBStore) IsRestaurantSaved(lineID string, restaurantName string) bool {
	var count int64
	ds.gorm.
		Model(&User{}).
		Joins("JOIN user_restaurants ON users.id = user_restaurants.user_id").
		Joins("JOIN restaurants ON user_restaurants.restaurant_id = restaurants.id").
		Where("users.line_id = ? AND restaurants.name = ?", lineID, restaurantName).
		Count(&count)
	return count > 0
}

func (ds *DBStore) GetRestaurantListByUser(lineID string) ([]model.RestaurantInfo, error) {
	var restaurants []model.RestaurantInfo
	result := ds.gorm.Table("restaurants").Select("restaurants.*").Joins("inner join user_restaurants on user_restaurants.restaurant_id = restaurants.id").Joins("inner join users on users.id = user_restaurants.user_id").Where("users.line_id = ?", lineID).Find(&restaurants)
	return restaurants, result.Error
}

func (ds *DBStore) CreateRestaurant(name string, phone string, address string) (Restaurant, error) {
	restaurant := Restaurant{
		Name:    name,
		Phone:   phone,
		Address: "",
	}
	result := ds.gorm.Create(&restaurant)
	if result.Error != nil {
		return Restaurant{}, result.Error
	}
	return restaurant, nil
}

func (ds *DBStore) AddRestaurantToUser(lineID string, restaurantID int) error {
	user, err := ds.GetUserByLineID(lineID)
	if err != nil {
		return err
	}
	userRestaurant := UserRestaurant{
		UserID:       user.ID,
		RestaurantID: restaurantID,
	}
	result := ds.gorm.Create(&userRestaurant)
	return result.Error
}

func (ds *DBStore) GetRestaurantByName(restaurantName string) (Restaurant, error) {
	var restaurant Restaurant
	result := ds.gorm.Table("restaurants").
		Where("restaurants.name = ?", restaurantName).
		First(&restaurant)
	return restaurant, result.Error
}

func (ds *DBStore) RemoveRestaurantFromUser(lineID string, restaurantName string) error {
	user, err := ds.GetUserByLineID(lineID)
	if err != nil {
		return err
	}

	restaurant, err := ds.GetRestaurantByName(restaurantName)
	if err != nil {
		return errors.New("GetRestaurantByName error!!")
	}

	result := ds.gorm.Table("user_restaurants").
		Where("user_restaurants.user_id = ? AND user_restaurants.restaurant_id = ?", user.ID, restaurant.ID).
		Delete(&UserRestaurant{})
	return result.Error
}

func (ds *DBStore) PickRestaurantFromUser(lineID string) (Restaurant, error) {
	var restaurant Restaurant
	result := ds.gorm.Table("restaurants").Select("restaurants.*").Joins("inner join user_restaurants on user_restaurants.restaurant_id = restaurants.id").Joins("inner join users on users.id = user_restaurants.user_id").Where("users.line_id = ?", lineID).Order("RANDOM()").Limit(1).First(&restaurant)
	return restaurant, result.Error
}

func (ds *DBStore) IsUserRestaurantEmpty(lineID string) bool {
	user, err := ds.GetUserByLineID(lineID)
	if err != nil {
		return true
	}
	var count int64
	ds.gorm.Model(&UserRestaurant{}).Where("user_id = ?", user.ID).Count(&count)
	return count == 0
}
