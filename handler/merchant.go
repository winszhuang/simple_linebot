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

func AddMerchant(userID string, merchantName string, phone string) string {
	for _, v := range userMap[userID].merchantList {
		if v.Name == merchantName {
			return "該商家已經存在摟"
		}
	}

	newMerchant := Merchant{merchantName, phone}
	user := userMap[userID]
	user.merchantList = append(user.merchantList, newMerchant)
	userMap[userID] = user
	fmt.Println("userMap[userID]: ", userMap[userID])
	return "商家" + merchantName + "新增成功!!"
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

func RemoveMerchant(userID string, merchantName string) {
	checkHaveUser(userID)

	user := userMap[userID]
	user.merchantList = filter(user.merchantList, func(merchant Merchant) bool {
		return merchant.Name == merchantName
	})
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

func filter(list []Merchant, filterFunc func(merchant Merchant) bool) []Merchant {
	filtered := []Merchant{}
	for _, e := range list {
		if filterFunc(e) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
