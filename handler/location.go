package handler

import (
	"fmt"
	"linebot/constants"
	"linebot/service"
	"linebot/util"

	"googlemaps.github.io/maps"
)

type LocationManager struct {
	MapService         service.MapService
	DataCache          map[string][]constants.RestaurantInfo
	NextPageTokenCache map[string]string
	Setting            LocationSetting
}

type LocationSetting struct {
	Radius   uint
	Type     maps.PlaceType
	Language string
	OpenNow  bool
}

type ListParams struct {
	Lat       float64
	Lng       float64
	PageIndex int
	PageSize  int
}

type SearchParams struct {
	Lat           float64
	Lng           float64
	NextPageToken string
}

const completeToken = "complete"

var defaultSetting = LocationSetting{
	Radius:   500,
	Type:     maps.PlaceTypeRestaurant,
	Language: "zh-TW",
	OpenNow:  true,
}

func NewLocationManager(mapService service.MapService, searchSetting LocationSetting) *LocationManager {
	return &LocationManager{
		MapService:         mapService,
		DataCache:          make(map[string][]constants.RestaurantInfo),
		NextPageTokenCache: make(map[string]string),
		Setting:            coverSetting(searchSetting),
	}
}

func (lm *LocationManager) List(params ListParams) ([]constants.RestaurantInfo, error) {
	if params.PageIndex <= 0 || params.PageSize <= 0 {
		return []constants.RestaurantInfo{}, fmt.Errorf("invalid params")
	}
	// init
	key := getKeyByLatLng(params.Lat, params.Lng)
	if _, ok := lm.DataCache[key]; !ok {
		lm.DataCache[key] = []constants.RestaurantInfo{}
		lm.NextPageTokenCache[key] = ""
	}

	resultList := util.Paginate(lm.DataCache[key], params.PageIndex, params.PageSize)
	for len(resultList) < params.PageSize {
		// if search is complete, force return cache
		// 確定fetch到底了，就不繼續fetch，並直接回傳
		if lm.NextPageTokenCache[key] == completeToken {
			return resultList, nil
		}

		newList, newPageToken, err := lm.Search(SearchParams{
			Lat:           params.Lat,
			Lng:           params.Lng,
			NextPageToken: lm.NextPageTokenCache[key],
		})
		if err != nil {
			return []constants.RestaurantInfo{}, err
		}

		// 確定沒有下一頁了
		if newPageToken == "" {
			lm.NextPageTokenCache[key] = completeToken
		} else {
			lm.NextPageTokenCache[key] = newPageToken
		}

		lm.DataCache[key] = append(lm.DataCache[key], newList...)
		resultList = util.Paginate(lm.DataCache[key], params.PageIndex, params.PageSize)
	}
	return resultList, nil
}

func (lm *LocationManager) Search(params SearchParams) ([]constants.RestaurantInfo, string, error) {
	request := &maps.NearbySearchRequest{
		Radius:   lm.Setting.Radius,
		Type:     lm.Setting.Type,
		Language: lm.Setting.Language,
		OpenNow:  lm.Setting.OpenNow,
		Location: &maps.LatLng{
			Lat: params.Lat,
			Lng: params.Lng,
		},
		PageToken: params.NextPageToken,
	}

	return lm.MapService.Search(request)
}

func getKeyByLatLng(lat float64, lng float64) string {
	return fmt.Sprintf("%f,%f", lat, lng)
}

func coverSetting(newSetting LocationSetting) LocationSetting {
	setting := defaultSetting
	if newSetting.Language != "" {
		setting.Language = newSetting.Language
	}
	if newSetting.Type != "" {
		setting.Type = newSetting.Type
	}
	if newSetting.Radius != 0 {
		setting.Radius = newSetting.Radius
	}
	setting.OpenNow = newSetting.OpenNow
	return setting
}
