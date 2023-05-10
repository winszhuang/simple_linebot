package store

import (
	"fmt"
	"linebot/config"
	"linebot/model"
	"linebot/service"
	"linebot/util"
	"log"
	"sync"

	"googlemaps.github.io/maps"
)

type LocationStore struct {
	// 所有換頁的token
	NextPageTokenCache map[string][]string
	DataCache          map[string][]model.RestaurantInfo
	MapService         service.MapService
	Setting            LocationSetting
	Mu                 sync.RWMutex
}

type LocationSetting struct {
	Radius   uint
	Type     maps.PlaceType
	Language string
	OpenNow  bool
}

const COMPLETE_TOKEN = "complete"

var defaultSetting = LocationSetting{
	Radius:   500,
	Type:     maps.PlaceTypeRestaurant,
	Language: "zh-TW",
	OpenNow:  true,
}

var locationStore *LocationStore

func InitLocationStore() *LocationStore {
	mapService, err := service.InitGoogleMapService(config.C.GoogleMapApiKey)
	if err != nil {
		log.Fatal(err)
	}
	locationStore = NewLocationStore(mapService, defaultSetting)
	return locationStore
}

func NewLocationStore(mapService service.MapService, locationSetting LocationSetting) *LocationStore {
	return &LocationStore{
		MapService:         mapService,
		DataCache:          make(map[string][]model.RestaurantInfo),
		NextPageTokenCache: make(map[string][]string),
		Setting:            coverSetting(locationSetting),
	}
}

func GetLocationStore() *LocationStore {
	return locationStore
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

func (lm *LocationStore) List(params ListParams) ([]model.RestaurantInfo, error) {
	if params.PageIndex <= 0 || params.PageSize <= 0 {
		return []model.RestaurantInfo{}, fmt.Errorf("invalid params")
	}
	// init
	key := getKeyByLatLng(params.Lat, params.Lng)
	lm.Mu.Lock()
	if _, ok := lm.DataCache[key]; !ok {
		lm.DataCache[key] = []model.RestaurantInfo{}
		lm.NextPageTokenCache[key] = []string{}
	}
	lm.Mu.Unlock()

	resultList := util.Paginate(lm.DataCache[key], params.PageIndex, params.PageSize)
	for len(resultList) < params.PageSize {
		// if search is complete, force return cache
		// 確定fetch到底了，就不繼續fetch，並直接回傳
		if include(lm.NextPageTokenCache[key], COMPLETE_TOKEN) {
			return resultList, nil
		}

		var lastPageToken string
		if len(lm.NextPageTokenCache[key]) > 0 {
			lastPageToken = lm.NextPageTokenCache[key][len(lm.NextPageTokenCache[key])-1]
		}
		newList, newPageToken, err := lm.Search(SearchParams{
			Lat:           params.Lat,
			Lng:           params.Lng,
			NextPageToken: lastPageToken,
		})
		if err != nil {
			return []model.RestaurantInfo{}, err
		}

		// 確定沒有下一頁了
		lm.Mu.Lock()
		var shouldBeAddToken string
		if newPageToken == "" {
			shouldBeAddToken = COMPLETE_TOKEN
		} else {
			shouldBeAddToken = newPageToken
		}
		if !include(lm.NextPageTokenCache[key], shouldBeAddToken) {
			lm.NextPageTokenCache[key] = append(lm.NextPageTokenCache[key], shouldBeAddToken)
			lm.DataCache[key] = append(lm.DataCache[key], newList...)
		}
		lm.Mu.Unlock()

		resultList = util.Paginate(lm.DataCache[key], params.PageIndex, params.PageSize)
	}
	return resultList, nil
}

func (lm *LocationStore) Search(params SearchParams) ([]model.RestaurantInfo, string, error) {
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

func include(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
