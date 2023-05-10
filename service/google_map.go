package service

import (
	"context"
	"fmt"
	"linebot/model"
	"log"
	"sync"

	"googlemaps.github.io/maps"
)

type MapService interface {
	Search(nearbySearchRequest *maps.NearbySearchRequest) ([]model.RestaurantInfo, string, error)
}

type GoogleMapService struct {
	client *maps.Client
	apiKey string
}

type PlaceDetailData struct {
	FormattedPhoneNumber     string
	InternationalPhoneNumber string
	URL                      string
	Index                    int
}

func InitGoogleMapService(apiKey string) (MapService, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &GoogleMapService{client, apiKey}, nil
}

// 根據經緯度搜尋附近店家
func (ms *GoogleMapService) Search(nearbySearchRequest *maps.NearbySearchRequest) ([]model.RestaurantInfo, string, error) {
	resp, err := ms.client.NearbySearch(context.Background(), nearbySearchRequest)
	if err != nil {
		return nil, "", err
	}

	resultList := ms.generateRestaurantList(resp.Results)
	resultList = ms.appendDetail(resultList)
	return resultList, resp.NextPageToken, nil
}

func (ms *GoogleMapService) generateRestaurantList(list []maps.PlacesSearchResult) []model.RestaurantInfo {
	result := make([]model.RestaurantInfo, 0)
	for _, place := range list {
		var photoUrl string
		if len(place.Photos) > 0 {
			photoUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", place.Photos[0].PhotoReference, ms.apiKey)
		} else {
			photoUrl = "https://mnapoli.fr/images/posts/null.png"
		}

		result = append(result, model.RestaurantInfo{
			Name:             place.Name,
			Rating:           place.Rating,
			UserRatingsTotal: place.UserRatingsTotal,
			Vicinity:         place.Vicinity,
			BusinessStatus:   place.BusinessStatus,
			Lat:              place.Geometry.Location.Lat,
			Lng:              place.Geometry.Location.Lng,
			PlaceID:          place.PlaceID,
			Photo:            photoUrl,
		})
	}
	return result
}

func (ms *GoogleMapService) appendDetail(list []model.RestaurantInfo) []model.RestaurantInfo {
	wg := sync.WaitGroup{}
	placeDetailChan := make(chan PlaceDetailData)
	for i, p := range list {
		wg.Add(1)
		go func(index int, place model.RestaurantInfo) {
			detail, err := ms.client.PlaceDetails(context.Background(), &maps.PlaceDetailsRequest{PlaceID: place.PlaceID})
			if err != nil {
				log.Printf("got error when get place %s detail: %s", place.Name, err)
				placeDetailChan <- PlaceDetailData{Index: index}
			} else {
				placeDetailChan <- PlaceDetailData{
					FormattedPhoneNumber:     detail.FormattedPhoneNumber,
					InternationalPhoneNumber: detail.InternationalPhoneNumber,
					URL:                      detail.URL,
					Index:                    index,
				}
			}
			wg.Done()
		}(i, p)
	}
	wg.Wait()
	close(placeDetailChan)

	for d := range placeDetailChan {
		list[d.Index].FormattedPhoneNumber = d.FormattedPhoneNumber
		list[d.Index].InternationalPhoneNumber = d.InternationalPhoneNumber
		list[d.Index].URL = d.URL
	}

	return list
}
