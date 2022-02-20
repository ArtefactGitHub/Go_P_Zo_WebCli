package main

import (
	"database/sql"
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
	User *User `json:"user"`
}

func (tr *PostUserResponseData) GetBaseData() *ResponseBase {
	return tr.ResponseBase
}
func (tr *PostUserResponseData) SetBaseData(statusCode int, err *myError) {
	tr.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}

type PostLoginResponseData struct {
	*ResponseBase
	*UserToken
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
	CategoryId         int       `json:"categoryid"`
	Message            string    `json:"message"`
	UserId             int       `json:"user_id"`
	AchievementDateStr string
}

type Zos struct {
	Zos []*Zo `json:"zos"`
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
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ColorId    int    `json:"color_id"`
	CreateType int    `json:"create_type"`
	UserId     int    `json:"user_id"`
}

type GetAllCategoryResponseData struct {
	*ResponseBase
	Categories []Category
}

func (d *GetAllCategoryResponseData) GetBaseData() *ResponseBase {
	return d.ResponseBase
}
func (d *GetAllCategoryResponseData) SetBaseData(statusCode int, err *myError) {
	d.ResponseBase = &ResponseBase{StatusCode: statusCode, Error: err}
}
