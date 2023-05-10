package model

type RestaurantInfo struct {
	Name   string
	Rating float32
	// 評分的總人數
	UserRatingsTotal int
	// 位置的詳細描述
	Vicinity                 string
	BusinessStatus           string
	Lat                      float64
	Lng                      float64
	PlaceID                  string
	Photo                    string
	FormattedPhoneNumber     string
	InternationalPhoneNumber string
	URL                      string
}
