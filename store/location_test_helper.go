package store

import (
	"fmt"
	"linebot/model"
	"linebot/util"
	"strconv"

	"googlemaps.github.io/maps"
)

type FakeMapService struct {
	pageSize      int
	pageIndex     int
	pageTokenList []string
	// 控制該fetch最多可以取得多少筆資料
	maxCount int
}

func NewFakeMapService(pageSize int, maxCount int) *FakeMapService {
	return &FakeMapService{
		// 預設頁數1
		pageIndex:     1,
		pageSize:      pageSize,
		pageTokenList: []string{},
		maxCount:      maxCount,
	}
}

func (f *FakeMapService) Search(request *maps.NearbySearchRequest) ([]model.RestaurantInfo, string, error) {
	f.pageIndex = f.getCurrentPageIndex(request.PageToken)

	list := []model.RestaurantInfo{}

	firstIndex := f.pageSize * (f.pageIndex - 1)
	lastIndex := f.pageSize * f.pageIndex

	fmt.Println("----")
	fmt.Printf("firstIndex: %d, lastIndex: %d\n", firstIndex, lastIndex)

	for i := firstIndex; i < lastIndex; i++ {
		// 超過最大筆數就跳出
		if f.maxCount < i+1 {
			return list, "", nil
		}
		list = append(list, model.RestaurantInfo{
			Name:    fmt.Sprintf("Restaurant-%d", i+1),
			PlaceID: strconv.Itoa(i + 1),
		})
	}

	newPageToken := util.RandomID(12)
	f.pageTokenList = append(f.pageTokenList, newPageToken)

	return list, newPageToken, nil
}

func (f *FakeMapService) getCurrentPageIndex(pageToken string) int {
	if pageToken == "" {
		return 1
	}

	existTokenIndex := FindIndex(f.pageTokenList, pageToken)
	if existTokenIndex == -1 {
		fmt.Println("token not found!!")
		return 1
	}

	return existTokenIndex + 2
}

func FindIndex(arr []string, elem string) int {
	for i, v := range arr {
		if v == elem {
			return i
		}
	}
	return -1
}
