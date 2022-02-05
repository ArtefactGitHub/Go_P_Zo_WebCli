package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	UserTokenHeaderName     = "X-Go_P_Zo_UserToken"
	AuthorizationHeaderName = "Authorization"
)

var accessToken *AccessToken

// ユーザー登録
func RequestSignup(familyName, givenName, email, password string) (*PostUserResponseData, error) {
	reqData := &PostUserReqestData{FamilyName: familyName, GivenName: givenName, Email: email, Password: password}
	resData := &PostUserResponseData{}
	err := Post("users", nil, reqData, resData, nil)
	return resData, err
}

// ログイン
func RequestSignin(email, password string) (*PostLoginResponseData, error) {
	reqData := &PostLoginRequestData{Identifier: email, Secret: password}
	resData := &PostLoginResponseData{}
	err := Post("usertokens", nil, reqData, resData, nil)
	return resData, err
}

// ユーザー情報取得
func RequestGetUser(p *UserToken) (*GetUserResponseData, error) {
	resData := &GetUserResponseData{}
	err := Get(fmt.Sprintf("users/%d", p.UserId), nil, nil, resData, p)
	return resData, err
}

// GET処理
func Get(api string, reqHeader map[string]string, reqBody interface{}, resData IResponse, p *UserToken) error {
	return request(api, http.MethodGet, nil, nil, resData, p)
}

// POST処理
func Post(api string, reqHeader map[string]string, reqBody interface{}, resData IResponse, p *UserToken) error {
	return request(api, http.MethodPost, nil, reqBody, resData, p)
}

func request(api string, method string, reqHeader map[string]string, reqBody interface{}, resData IResponse, p *UserToken) error {
	h := map[string]string{}

	// Expire API Access token
	if isInvalidToken() {
		_, err := getAccessToken()
		if err != nil {
			return err
		}
	}
	h[AuthorizationHeaderName] = fmt.Sprintf("Bearer %s", accessToken.Jwt)

	if p != nil {
		// Expire User token
		if p.IsExpired() {
			resData.SetBaseData(http.StatusUnauthorized, &myError{"UserToken expired"})
			return nil
		}
		h[UserTokenHeaderName] = p.Token
	}

	// add header
	for k, v := range reqHeader {
		if _, ok := h[k]; !ok {
			h[k] = v
		}
	}

	// send body
	if reqBody != nil {
		json, _ := json.MarshalIndent(reqBody, "", "\t")
		body := bytes.NewReader(json)
		return requestCore(api, method, h, body, resData)
	} else {
		return requestCore(api, method, h, nil, resData)
	}
}

// リクエスト処理
func requestCore(api string, method string, reqHeader map[string]string, body io.Reader, resData interface{}) error {
	endpoint := "http://localhost:8000/api/v1/"
	log.Printf("[Post] api: %s, method: %s, reqHeader: %v", api, method, reqHeader)

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	reqURL.Path = path.Join(reqURL.Path, api)
	log.Printf("[Post] reqURL.Path: %s", reqURL.Path)

	req, err := http.NewRequest(method, reqURL.String(), body)
	if err != nil {
		return err
	}

	for key, value := range reqHeader {
		req.Header.Set(key, value)
	}

	_, err = doRequest(req, resData)
	return err
}

// トークンが有効か
func isInvalidToken() bool {
	return accessToken == nil ||
		(accessToken != nil && accessToken.ExpiresAt < time.Now().Unix())
}

// トークン取得処理
func getAccessToken() (*AccessToken, error) {
	tokenRes, err := getTokenRequest()
	if err != nil {
		return nil, err
	}
	if tokenRes.ResponseBase.Error != nil {
		return nil, tokenRes.ResponseBase.Error
	}

	accessToken = tokenRes.AccessToken
	return accessToken, nil
}

// トークン取得リクエスト
func getTokenRequest() (*TokenResponse, error) {
	resData := &TokenResponse{}

	var body *bytes.Reader
	json, _ := json.MarshalIndent(TokenRequest{Id: Cfg.ClientId, Secret: Cfg.ClientSecret}, "", "\t")
	body = bytes.NewReader(json)

	err := requestCore("client/token", http.MethodPost, nil, body, resData)
	return resData, err
}

// httpリクエスト処理
func doRequest(req *http.Request, respBody interface{}) (int, error) {
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	log.Printf("[do] body: \n%v", string(body))
	if err != nil {
		return 0, err
	}

	if err := json.Unmarshal(body, respBody); err != nil {
		return 0, err
	}

	return res.StatusCode, nil
}
