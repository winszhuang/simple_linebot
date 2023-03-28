package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
)

func InitDB() error {
	userName := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	source := fmt.Sprintf("postgresql://%s:%s@localhost/%s?sslmode=disable", userName, password, dbName)
	sqlDB, err := sql.Open("postgres", source)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

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
	return nil
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
	dbController.Model(&User{}).Where("line_id = ?", lineID).Count(&count)
	return count > 0
}

func GetRestaurants(lineID string) ([]Restaurant, error) {
	var restaurants []Restaurant
	result := dbController.Table("restaurants").Select("restaurants.*").Joins("inner join user_restaurants on user_restaurants.restaurant_id = restaurants.id").Joins("inner join users on users.id = user_restaurants.user_id").Where("users.line_id = ?", lineID).Find(&restaurants)
	return restaurants, result.Error
}

func AddRestaurant(lineID string, restaurantID int) error {
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

func RemoveUserRestaurant(lineID string, restaurantID int) error {
	user, err := GetUserByLineID(lineID)
	if err != nil {
		return err
	}
	result := dbController.Where("user_id = ? AND restaurant_id = ?", user.ID, restaurantID).Delete(&UserRestaurant{})
	return result.Error
}
