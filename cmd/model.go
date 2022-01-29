package main

import (
	"database/sql"
	"time"
)

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

type PostLoginResponseData struct {
	*ResponseBase
	*SessionData
}

type User struct {
	Id         int          `json:"id"`
	GivenName  string       `json:"given_name"`
	FamilyName string       `json:"family_name"`
	Email      string       `json:"email"`
	Password   string       `json:"password"`
	CreatedAt  time.Time    `json:"createdat"`
	UpdatedAt  sql.NullTime `json:"updatedat"`
}

type SessionData struct {
	SessionId  string
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Expire     time.Time
}
