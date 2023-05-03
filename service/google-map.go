package service

import (
	"context"
	"fmt"
	"linebot/constants"
	"os"

	"googlemaps.github.io/maps"
)

var (
	apiKey string
	client *maps.Client
)

func InitMapsClient() error {
	var err error
	apiKey = os.Getenv("GOOGLE_MAP_API_KEY")
	client, err = maps.NewClient(maps.WithAPIKey(apiKey))
	return err
}

// 根據經緯度搜尋附近餐館
func SearchRestaurantByLatLng(lat, lng float64, radius uint, openNow bool, nextPageToken string) ([]constants.RestaurantInfo, string, error) {
	r := &maps.NearbySearchRequest{
		Radius:    radius,
		Type:      maps.PlaceTypeRestaurant,
		Language:  "zh-TW",
		OpenNow:   openNow,
		PageToken: nextPageToken,
	}

	r.Location = &maps.LatLng{
		Lat: lat,
		Lng: lng,
	}

	resp, err := client.NearbySearch(context.Background(), r)
	if err != nil {
		return nil, "", err
	}

	result := make([]constants.RestaurantInfo, 0)
	for i, place := range resp.Results {
		// linebot限定最多一次只能12個，鮮寫死
		if i > 11 {
			break
		}
		var photoUrl string
		if len(place.Photos) > 0 {
			photoUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", place.Photos[0].PhotoReference, apiKey)
		} else {
			photoUrl = "https://mnapoli.fr/images/posts/null.png"
		}

		result = append(result, constants.RestaurantInfo{
			Name:             place.Name,
			Rating:           place.Rating,
			UserRatingsTotal: place.UserRatingsTotal,
			Vicinity:         place.Vicinity,
			BusinessStatus:   place.BusinessStatus,
			Lat:              place.Geometry.Location.Lat,
			Lng:              place.Geometry.Location.Lng,
			ID:               place.PlaceID,
			Photo:            photoUrl,
		})
	}

	return result, resp.NextPageToken, nil
}
