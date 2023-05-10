package store

import (
	"fmt"
	"linebot/constant"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const MAX_COUNT = 10

func getKey(index int) string {
	return fmt.Sprintf("ID-%d", index)
}

// 測試同時間訪問map會不會有問題
func TestTempUserInputStore(t *testing.T) {
	// init store
	us := GetUserInputStore()

	t.Run("load multi user in same time", func(t *testing.T) {
		wg := sync.WaitGroup{}
		// 同時讀取
		for i := 0; i < MAX_COUNT; i++ {
			wg.Add(1)
			go func(index int) {
				userId := getKey(index)
				us.LoadUserInputData(userId)
				wg.Done()
			}(i)
		}
		wg.Wait()

		for i := 0; i < MAX_COUNT; i++ {
			fmt.Println(*us)
			if _, ok := us.Load(getKey(i)); !ok {
				t.Errorf("load user %s data fail", getKey(i))
			}
		}
	})

	t.Run("update multi user data in same time", func(t *testing.T) {
		wg := sync.WaitGroup{}
		// 同時讀取
		for i := 0; i < MAX_COUNT; i++ {
			wg.Add(1)
			go func(index int) {
				userId := getKey(index)
				us.LoadUserInputData(userId)
				wg.Done()
			}(i)
		}
		wg.Wait()

		for i := 0; i < MAX_COUNT; i++ {
			wg.Add(1)
			go func(index int) {
				userId := getKey(index)
				userData := us.LoadUserInputData(userId)
				if index%3 == 0 {
					userData.SetMode(constant.Add)
				} else if index%3 == 1 {
					userData.SetMode(constant.List)
				} else if index%3 == 2 {
					userData.SetMode(constant.Remove)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()

		for i := 0; i < MAX_COUNT; i++ {
			userId := getKey(i)
			userData := us.LoadUserInputData(userId)
			if i%3 == 0 {
				require.Equal(t, userData, &InputCache{userId: userId, Mode: constant.Add})
			} else if i%3 == 1 {
				require.Equal(t, userData, &InputCache{userId: userId, Mode: constant.List})
			} else if i%3 == 2 {
				require.Equal(t, userData, &InputCache{userId: userId, Mode: constant.Remove})
			}

		}
	})

	t.Run("write same user in same time", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < MAX_COUNT; i++ {
			wg.Add(3)
			go func(index int) {
				userId := getKey(index)
				userData := us.LoadUserInputData(userId)
				go func() {
					userData.SetMode(constant.Add)
					wg.Done()
				}()
				go func() {
					userData.SetQuestion(constant.Name)
					wg.Done()
				}()
				go func() {
					userData.SetData(func(m *Merchant) *Merchant {
						m.Name = fmt.Sprintf("Name-%d", index)
						m.Phone = fmt.Sprintf("Phone-%d", index)
						return m
					})
					wg.Done()
				}()
			}(i)
		}
		wg.Wait()

		for i := 0; i < MAX_COUNT; i++ {
			userData := us.LoadUserInputData(getKey(i))
			require.Equal(t, userData, &InputCache{
				userId:   getKey(i),
				Mode:     constant.Add,
				Question: constant.Name,
				Data: &Merchant{
					Name:  fmt.Sprintf("Name-%d", i),
					Phone: fmt.Sprintf("Phone-%d", i),
				},
			})
		}
	})
}
