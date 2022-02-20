package main

import (
	"time"
)

const TimeLayout = "2006-01-02"

type MypageUserGetModel struct {
	*ResponseBase
	Name  string
	Email string
}

func NewMypageUserGetModel(name, email string, base *ResponseBase) *MypageUserGetModel {
	return &MypageUserGetModel{Name: name, Email: email, ResponseBase: base}
}

type MypageZosGetModel struct {
	*ResponseBase
	Name       string
	Email      string
	Zos        Zos
	Categories Categories
	TodayStr   string
}

func NewMypageZosGetModel(name, email string, zos Zos, categories Categories, base *ResponseBase) *MypageZosGetModel {
	m := &MypageZosGetModel{Name: name, Email: email, Zos: zos, Categories: categories, TodayStr: time.Now().Format("2006-01-02"), ResponseBase: base}
	m.Zos.Sort()
	m.Zos.SetCategoryName(categories)
	return m
}

type requestZo struct {
	AchievementDate time.Time `json:"achievementdate"`
	Exp             int       `json:"exp"`
	CategoryId      int       `json:"category_id"`
	Message         string    `json:"message"`
}

func NewRequestZo(
	achievementDate time.Time,
	exp int,
	categoryId int,
	message string,
) *requestZo {
	return &requestZo{
		AchievementDate: achievementDate,
		Exp:             exp,
		CategoryId:      categoryId,
		Message:         message,
	}
}

type responseZo struct {
	Id              int       `json:"id"`
	AchievementDate time.Time `json:"achievementdate"`
	Exp             int       `json:"exp"`
	CategoryId      int       `json:"category_id"`
	Message         string    `json:"message"`
}

type PostZoResponseData struct {
	*ResponseBase
	*responseZo
}

func (d *PostZoResponseData) GetBaseData() *ResponseBase {
	return d.ResponseBase
}
func (d *PostZoResponseData) SetBaseData(statusCode int, err *myError) {
	d.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}
