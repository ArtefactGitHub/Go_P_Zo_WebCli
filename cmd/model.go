package main

import (
	"database/sql"
	"fmt"
	"sort"
	"time"
)

type IResponse interface {
	GetBaseData() *ResponseBase
	SetBaseData(statusCode int, error *myError)
}

type ResponseBase struct {
	StatusCode int      `json:"statuscode"`
	Error      *myError `json:"error"`
}

type myError struct {
	Message string `json:"message"`
}

func (e *myError) Error() string {
	return e.Message
}

type AccessToken struct {
	Jwt       string `json:"jwt"`
	ExpiresAt int64  `json:"expires_at"`
}

type TokenRequest struct {
	Id     int    `json:"id"`
	Secret string `json:"secret"`
}

type TokenResponse struct {
	*ResponseBase
	AccessToken *AccessToken `json:"access_token"`
}

func (tr *TokenResponse) GetBaseData() *ResponseBase {
	return tr.ResponseBase
}
func (tr *TokenResponse) SetBaseData(statusCode int, err *myError) {
	tr.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}

type PostUserReqestData struct {
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type PostLoginRequestData struct {
	Identifier string `json:"identifier"`
	Secret     string `json:"secret"`
}

type PostUserResponseData struct {
	*ResponseBase
	User      *User      `json:"user"`
	UserToken *UserToken `json:"usertoken"`
}

func (v *PostUserResponseData) GetBaseData() *ResponseBase {
	return v.ResponseBase
}
func (v *PostUserResponseData) SetBaseData(statusCode int, err *myError) {
	v.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}

func (v *PostUserResponseData) String() string {
	str := ""
	if v.ResponseBase == nil {
		str = str + "Base is nil \n"
	} else {
		str = str + fmt.Sprintf("Base: %#v \n", v.ResponseBase)
	}
	if v.User == nil {
		str = str + "User is nil \n"
	} else {
		str = str + fmt.Sprintf("User: %#v \n", v.User)
	}
	return str
}

type PostLoginResponseData struct {
	*ResponseBase
	UserToken *UserToken `json:"user_token"`
}

func (tr *PostLoginResponseData) GetBaseData() *ResponseBase {
	return tr.ResponseBase
}
func (tr *PostLoginResponseData) SetBaseData(statusCode int, err *myError) {
	tr.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}

type GetUserResponseData struct {
	*ResponseBase
	*User `json:"user"`
}

func (tr *GetUserResponseData) GetBaseData() *ResponseBase {
	return tr.ResponseBase
}
func (tr *GetUserResponseData) SetBaseData(statusCode int, err *myError) {
	tr.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}

type User struct {
	Id         int    `json:"id"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}

type UserToken struct {
	Id        int          `json:"id"`
	UserId    int          `json:"user_id"`
	Token     string       `json:"token"`
	ExpiredAt time.Time    `json:"expiredat"`
	CreatedAt time.Time    `json:"createdat"`
	UpdatedAt sql.NullTime `json:"updatedat"`
}

func (t *UserToken) IsExpired() bool {
	return t.ExpiredAt.Unix() < time.Now().Unix()
}

type SessionData struct {
	SessionId string
	UserToken
}

type Zo struct {
	Id                 int       `json:"id"`
	AchievementDate    time.Time `json:"achievementdate"`
	Exp                int       `json:"exp"`
	CategoryId         int       `json:"category_id"`
	Message            string    `json:"message"`
	UserId             int       `json:"user_id"`
	AchievementDateStr string
	CategoryName       string
}

type Zos struct {
	Zos []*Zo `json:"zos"`
}

func (z *Zos) Sort() {
	sort.SliceStable(z.Zos, func(i, j int) bool { return z.Zos[i].Id > z.Zos[j].Id })
	sort.SliceStable(z.Zos, func(i, j int) bool { return z.Zos[i].AchievementDate.Unix() > z.Zos[j].AchievementDate.Unix() })
}

func (z *Zos) SetCategoryName(categories Categories) {
	for _, v := range z.Zos {
		v.CategoryName = categories.GetName(v.CategoryId)
	}
}

type GetAllZoResponseData struct {
	*ResponseBase
	*Zos
}

func (d *GetAllZoResponseData) GetBaseData() *ResponseBase {
	return d.ResponseBase
}
func (d *GetAllZoResponseData) SetBaseData(statusCode int, err *myError) {
	d.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}

type Category struct {
	Id      int    `json:"id"`
	Number  int    `json:"number"`
	Name    string `json:"name"`
	ColorId int    `json:"color_id"`
	UserId  int    `json:"user_id"`
}

type Categories struct {
	Categories []Category `json:"categories"`
}

func (m *Categories) GetName(id int) string {
	for _, v := range m.Categories {
		if id == v.Id {
			return v.Name
		}
	}
	return ""
}

type GetAllCategoryResponseData struct {
	*ResponseBase
	Categories
}

func (d *GetAllCategoryResponseData) GetBaseData() *ResponseBase {
	return d.ResponseBase
}
func (d *GetAllCategoryResponseData) SetBaseData(statusCode int, err *myError) {
	d.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}
