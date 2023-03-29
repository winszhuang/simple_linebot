package handler

import (
	"fmt"
	c "linebot/constants"
)

type Merchant struct {
	Name  string
	Phone string
}

type UserInputInfo struct {
	mode     c.Directive
	question c.Question
	data     *Merchant
}

var (
	inputData = make(map[string]*UserInputInfo)
)

func LoadUserInputData(userId string) *UserInputInfo {
	if _, ok := inputData[userId]; !ok {
		inputData[userId] = &UserInputInfo{}
	}
	return inputData[userId]
}

func RemoveUserInputData(userId string) {
	delete(inputData, userId)
}

func ResetUserInputData(userId string) {
	userInputData := inputData[userId]
	userInputData.
		SetMode("").
		SetQuestion("").
		SetData(func(m *Merchant) *Merchant {
			m.Name = ""
			m.Phone = ""
			return m
		})
}

func (inputInfo *UserInputInfo) SetMode(mode c.Directive) *UserInputInfo {
	inputInfo.mode = mode
	return inputInfo
}

func (inputInfo *UserInputInfo) SetQuestion(question c.Question) *UserInputInfo {
	inputInfo.question = question
	return inputInfo
}

func (inputInfo *UserInputInfo) SetData(fn func(*Merchant) *Merchant) *UserInputInfo {
	if inputInfo.data == nil {
		inputInfo.data = &Merchant{}
	}
	inputInfo.data = fn(inputInfo.data)
	return inputInfo
}

func (inputInfo *UserInputInfo) GetMode() c.Directive {
	return inputInfo.mode
}

func (inputInfo *UserInputInfo) GetQuestion() c.Question {
	return inputInfo.question
}

func (inputInfo *UserInputInfo) GetData() *Merchant {
	return inputInfo.data
}

func (inputInfo *UserInputInfo) IsInMode() bool {
	fmt.Println("該使用者當前狀態是", inputInfo)
	return inputInfo.mode != ""
}
