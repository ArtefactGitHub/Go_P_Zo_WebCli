package main

import (
	"sort"
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

type MypageGetModel struct {
	*ResponseBase
	Name       string
	Email      string
	Zos        Zos
	Categories []Category
}

func NewMypageGetModel(name, email string, zos Zos, categories []Category, base *ResponseBase) *MypageGetModel {
	m := &MypageGetModel{Name: name, Email: email, Zos: zos, Categories: categories, ResponseBase: base}
	m.Sort()
	return m
}

func (m *MypageGetModel) Sort() {
	sort.SliceStable(m.Zos.Zos, func(i, j int) bool { return m.Zos.Zos[i].Id > m.Zos.Zos[j].Id })
	sort.SliceStable(m.Zos.Zos, func(i, j int) bool { return m.Zos.Zos[i].AchievementDate.Unix() > m.Zos.Zos[j].AchievementDate.Unix() })
}

type requestZo struct {
	AchievementDate time.Time `json:"achievementdate"`
	Exp             int       `json:"exp"`
	CategoryId      int       `json:"categoryid"`
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
	CategoryId      int       `json:"categoryid"`
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
