package merchant

import (
	"fmt"
	"math/rand"
)

type Merchant struct {
	Name  string
	Phone string
}

type User struct {
	merchantList []Merchant
}

var (
	userMap map[string]User = make(map[string]User)
)

func AddMerchant(userID string, merchantName string, phone string) (string, bool) {
	if IsMerchantExist(userID, merchantName) {
		return "該商家已經存在摟", false
	}

	newMerchant := Merchant{merchantName, phone}
	user := userMap[userID]
	user.merchantList = append(user.merchantList, newMerchant)
	userMap[userID] = user
	fmt.Println("userMap[userID]: ", userMap[userID])
	return "商家" + merchantName + "新增成功!!", true
}

func IsMerchantExist(userID string, merchantName string) bool {
	for _, v := range userMap[userID].merchantList {
		if v.Name == merchantName {
			return true
		}
	}
	return false
}

func ViewMerchants(userID string) string {
	checkHaveUser(userID)
	var str string
	for _, v := range userMap[userID].merchantList {
		str += "---" + "\n"
		str += v.Name + "\n"
		str += v.Phone + "\n"
	}
	return str
}

func RemoveMerchant(userID string, merchantName string) (string, bool) {
	checkHaveUser(userID)

	user := userMap[userID]
	newList, hasChange := filter(user.merchantList, func(merchant Merchant) bool {
		return merchant.Name != merchantName
	})
	if hasChange {
		user.merchantList = newList
		userMap[userID] = user
		return "刪除店家成功", true
	} else {
		return "找不到該店家名稱，請重新輸入", false
	}
}

func PickMerchant(userID string) string {
	checkHaveUser(userID)

	merchantList := userMap[userID].merchantList
	merchantLen := len(merchantList)

	if merchantLen == 0 {
		return "尚未有店家，請先加入店家再做隨機選店!!"
	}

	randIndex := rand.Intn(merchantLen)

	return merchantList[randIndex].Name + "\n" + merchantList[randIndex].Phone
}

func checkHaveUser(userID string) {
	_, ok := userMap[userID]
	if !ok {
		userMap[userID] = User{merchantList: []Merchant{}}
	}
}

func filter(list []Merchant, filterFunc func(merchant Merchant) bool) ([]Merchant, bool) {
	filtered := []Merchant{}
	for _, e := range list {
		if filterFunc(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered, len(list) != len(filtered)
}
