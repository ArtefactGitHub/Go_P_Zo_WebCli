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

type MyPageZosPostModel struct {
	*ResponseBase
	*responseZo
}

func (d *MyPageZosPostModel) GetBaseData() *ResponseBase {
	return d.ResponseBase
}
func (d *MyPageZosPostModel) SetBaseData(statusCode int, err *myError) {
	d.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}

// UserCategory
type requestUserCategory struct {
	Name    string `json:"name"`
	ColorId int    `json:"color_id"`
	UserId  int    `json:"user_id"`
}

func NewRequestCategory(name string, colorId, userId int) *requestUserCategory {
	return &requestUserCategory{
		Name: name, ColorId: colorId, UserId: userId,
	}
}

type responseCategory struct {
	Id      int    `json:"id"`
	Number  int    `json:"number"`
	Name    string `json:"name"`
	ColorId int    `json:"color_id"`
	UserId  int    `json:"user_id"`
}

type PostUserCategoryResponseData struct {
	*ResponseBase
	*responseCategory
}

func (d *PostUserCategoryResponseData) GetBaseData() *ResponseBase {
	return d.ResponseBase
}
func (d *PostUserCategoryResponseData) SetBaseData(statusCode int, err *myError) {
	d.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}
