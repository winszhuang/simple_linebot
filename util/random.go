package util

import (
	"linebot/model"
	"math/rand"
	"time"
)

const randomString = "abcdefghijklmnopqrstuvwxyz"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	// 為了保證每次運行程序時，隨機生成的數字序列都是不同的
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomRestaurant() model.RestaurantInfo {
	return model.RestaurantInfo{
		Name:    RandomName(),
		PlaceID: RandomID(12),
	}
}

func RandomName() string {
	var str string

	k := len(randomString)
	for i := 0; i < 12; i++ {
		c := randomString[rand.Intn(k)]
		str += string(c)
	}

	return str
}

func RandomID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
