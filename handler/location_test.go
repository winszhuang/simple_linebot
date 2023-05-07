package handler

import (
	"fmt"
	"linebot/constants"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"googlemaps.github.io/maps"
)

func generateRestaurantList(firstIndex int, lastIndex int) []constants.RestaurantInfo {
	list := []constants.RestaurantInfo{}
	for i := firstIndex; i < lastIndex; i++ {
		list = append(list, constants.RestaurantInfo{
			Name: fmt.Sprintf("Restaurant-%d", i+1),
			ID:   strconv.Itoa(i + 1),
		})
	}
	return list
}

func TestLocationManager_List(t *testing.T) {
	t.Deadline()
	timeout := time.After(3 * time.Second)
	done := make(chan bool)

	go func() {
		fakeMapService := InitFakeMapService(20, 43)
		locationManager := NewLocationManager(fakeMapService, LocationSetting{
			Radius:   500,
			Type:     maps.PlaceTypeRestaurant,
			Language: "zh-TW",
			OpenNow:  true,
		})

		// LatLng不是重點，重點是測試PageIndex和PageSize是否符合回傳結果
		tests := []struct {
			name    string
			args    ListParams
			want    []constants.RestaurantInfo
			wantErr error
		}{
			{
				name: "test1",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 1,
					PageSize:  10,
				},
				want:    generateRestaurantList(0, 10),
				wantErr: nil,
			},
			{
				name: "test2",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 3,
					PageSize:  15,
				},
				want:    generateRestaurantList(30, 43),
				wantErr: nil,
			},
			{
				name: "test3",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 4,
					PageSize:  12,
				},
				want:    generateRestaurantList(36, 43),
				wantErr: nil,
			},
			{
				name: "test3",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 3,
					PageSize:  40,
				},
				want:    []constants.RestaurantInfo{},
				wantErr: nil,
			},
			{
				name: "negative page index",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: -1,
					PageSize:  12,
				},
				want:    []constants.RestaurantInfo{},
				wantErr: fmt.Errorf("invalid params"),
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := locationManager.List(tt.args)

				// err有可能是nil，如果是nil，就無法呼叫Error method
				if err == nil && tt.wantErr == nil {
					// both err and tt.wantErr are nil, so they are equal
				} else if err != nil && tt.wantErr != nil {
					// both err and tt.wantErr are not nil, so compare their error messages
					if err.Error() != tt.wantErr.Error() {
						t.Errorf("got %v, want %v", err, tt.wantErr)
					}
				} else {
					// one of them is nil and the other is not, so they are not equal
					t.Errorf("got %v, want %v", err, tt.wantErr)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("locationManager.List(tt.args) = %v, want %v", got, tt.want)
				}
			})
		}
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}

// 測試多個人同時取同一個定位點的情況
func TestLocationManager_List_DeadLock(t *testing.T) {
	timeout := time.After(3 * time.Second)
	wg := sync.WaitGroup{}
	done := make(chan bool)

	go func() {
		fakeMapService := InitFakeMapService(20, 43)
		locationManager := NewLocationManager(fakeMapService, LocationSetting{
			Radius:   500,
			Type:     maps.PlaceTypeRestaurant,
			Language: "zh-TW",
			OpenNow:  true,
		})

		tests := []struct {
			name    string
			args    ListParams
			want    []constants.RestaurantInfo
			wantErr error
		}{
			{
				name: "test1",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 1,
					PageSize:  10,
				},
				want:    generateRestaurantList(0, 10),
				wantErr: nil,
			},
			{
				name: "test2",
				args: ListParams{
					Lat:       23.4,
					Lng:       120.6,
					PageIndex: 1,
					PageSize:  14,
				},
				want:    generateRestaurantList(0, 14),
				wantErr: nil,
			},
			{
				name: "test3",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 3,
					PageSize:  15,
				},
				want:    generateRestaurantList(30, 43),
				wantErr: nil,
			},
			{
				name: "test4",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 2,
					PageSize:  18,
				},
				want:    generateRestaurantList(18, 36),
				wantErr: nil,
			},
			{
				name: "test5",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 4,
					PageSize:  12,
				},
				want:    generateRestaurantList(36, 43),
				wantErr: nil,
			},
			{
				name: "test6",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: 3,
					PageSize:  40,
				},
				want:    []constants.RestaurantInfo{},
				wantErr: nil,
			},
			{
				name: "negative page index",
				args: ListParams{
					Lat:       23.7,
					Lng:       120.5,
					PageIndex: -1,
					PageSize:  12,
				},
				want:    []constants.RestaurantInfo{},
				wantErr: fmt.Errorf("invalid params"),
			},
		}
		for _, ttt := range tests {
			wg.Add(1)
			go func(tt struct {
				name    string
				args    ListParams
				want    []constants.RestaurantInfo
				wantErr error
			}) {
				t.Run(tt.name, func(t *testing.T) {
					got, err := locationManager.List(tt.args)

					// err有可能是nil，如果是nil，就無法呼叫Error method
					if err == nil && tt.wantErr == nil {
						// both err and tt.wantErr are nil, so they are equal
					} else if err != nil && tt.wantErr != nil {
						// both err and tt.wantErr are not nil, so compare their error messages
						if err.Error() != tt.wantErr.Error() {
							t.Errorf("got %v, want %v", err, tt.wantErr)
						}
					} else {
						// one of them is nil and the other is not, so they are not equal
						t.Errorf("got %v, want %v", err, tt.wantErr)
					}

					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("locationManager.List(tt.args) = %v, want %v", got, tt.want)
					}
					wg.Done()
				})
			}(ttt)
		}
		wg.Wait()
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}

}
