package store

import (
	"fmt"
	"linebot/constant"
	"sync"
)

type Merchant struct {
	Name  string
	Phone string
}

type TempUserInputStore struct {
	sync.Map
}

type InputCache struct {
	userId   string
	Mode     constant.Directive
	Question constant.Question
	Data     *Merchant
}

var userInputStore = &TempUserInputStore{}
var mu sync.Mutex

func GetUserInputStore() *TempUserInputStore {
	return userInputStore
}

func (us *TempUserInputStore) LoadUserInputData(userId string) *InputCache {
	value, _ := us.LoadOrStore(userId, &InputCache{userId: userId})
	return value.(*InputCache)
}

func (us *TempUserInputStore) RemoveUserInputData(userId string) {
	us.Delete(userId)
}

func (us *TempUserInputStore) ResetUserInputData(userId string) {
	us.Store(userId, &InputCache{Data: &Merchant{}})
}

func (inputInfo *InputCache) SetMode(mode constant.Directive) *InputCache {
	inputInfo.Mode = mode
	return inputInfo
}

func (inputInfo *InputCache) SetQuestion(question constant.Question) *InputCache {
	inputInfo.Question = question
	return inputInfo
}

func (inputInfo *InputCache) SetData(fn func(*Merchant) *Merchant) *InputCache {
	if inputInfo.Data == nil {
		inputInfo.Data = &Merchant{}
	}
	inputInfo.Data = fn(inputInfo.Data)
	return inputInfo
}

func (inputInfo *InputCache) Reset() {
	userInputStore.Store(inputInfo.userId, &InputCache{Data: &Merchant{}})
}

func (inputInfo *InputCache) IsInMode() bool {
	fmt.Println("該使用者當前狀態是", inputInfo)
	return inputInfo.Mode != ""
}
