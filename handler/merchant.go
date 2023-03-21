package merchant

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
	checkHaveUser(userID)
	user, ok := userMap[userID]
	if ok {
		return "該商家已經存在摟"
	} else {
		newMerchant := Merchant{merchantName, phone}
		user.merchantList = append(user.merchantList, newMerchant)
		return "商家" + merchantName + "新增成功!!"
	}
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
