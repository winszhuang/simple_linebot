package service

import (
	"context"
	"fmt"
	"linebot/constants"

	"googlemaps.github.io/maps"
)

type MapService interface {
	Search(nearbySearchRequest *maps.NearbySearchRequest) ([]constants.RestaurantInfo, string, error)
}

type GoogleMapService struct {
	client *maps.Client
	apiKey string
}

func InitGoogleMapService(apiKey string) (MapService, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &GoogleMapService{client, apiKey}, nil
}

// 根據經緯度搜尋附近店家
func (ms *GoogleMapService) Search(nearbySearchRequest *maps.NearbySearchRequest) ([]constants.RestaurantInfo, string, error) {
	resp, err := ms.client.NearbySearch(context.Background(), nearbySearchRequest)
	if err != nil {
		return nil, "", err
	}

	result := make([]constants.RestaurantInfo, 0)
	for _, place := range resp.Results {
		var photoUrl string
		if len(place.Photos) > 0 {
			photoUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", place.Photos[0].PhotoReference, ms.apiKey)
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
