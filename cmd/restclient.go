package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

var accessToken *AccessToken

// ユーザー登録
func RequestSignup(familyName, givenName, email, password string) (*PostUserResponseData, error) {
	reqData := &PostUserReqestData{FamilyName: familyName, GivenName: givenName, Email: email, Password: password}
	resData := &PostUserResponseData{}
	err := Post("users", reqData, &resData)
	return resData, err
}

// ログイン
func RequestSignin(email, password string) (*PostLoginResponseData, error) {
	reqData := &PostLoginRequestData{Identifier: email, Secret: password}
	resData := &PostLoginResponseData{}
	err := Post("usertokens", reqData, &resData)
	return resData, err
}

// POST処理
func Post(api string, reqBody interface{}, resData interface{}) error {
	// Expire token
	if isInvalidToken() {
		_, err := getAccessToken()
		if err != nil {
			return err
		}
	}

	reqHeader := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", accessToken.Jwt)}
	return request(api, http.MethodPost, reqHeader, reqBody, resData)
}

// リクエスト処理
func request(api string, method string, reqHeader map[string]string, reqBody interface{}, resData interface{}) error {
	endpoint := "http://localhost:8000/api/v1/"
	log.Printf("[Post] api: %s, method: %s, reqHeader: %v, reqData: %v", api, method, reqHeader, reqBody)

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	reqURL.Path = path.Join(reqURL.Path, api)
	log.Printf("[Post] reqURL.Path: %s", reqURL.Path)

	var body *bytes.Reader
	if reqBody != nil {
		json, _ := json.MarshalIndent(reqBody, "", "\t")
		body = bytes.NewReader(json)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL.String(), body)
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
	resData := TokenResponse{}
	err := request("client/token", http.MethodPost, nil, &TokenRequest{Id: Cfg.ClientId, Secret: Cfg.ClientSecret}, &resData)
	return &resData, err
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
