package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

// ユーザー情報
func (s *service) GetMypageUser(userToken *UserToken) (*MypageUserGetModel, error) {
	res, err := RequestGetUser(userToken)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}

	resCategory, err := RequestGetAllCategory(userToken)
	if err != nil || resCategory.StatusCode != http.StatusOK {
		return nil, err
	}

	result := NewMypageUserGetModel(res.FamilyName+res.GivenName, res.Email, resCategory.Categories, res.ResponseBase)
	return result, err
}

func (s *service) PostMypageUser(userToken *UserToken, values url.Values) error {
	req, err := validationAddCategory(values)
	if err != nil {
		return err
	}

	resUser, err := RequestGetUser(userToken)
	if err != nil || resUser.StatusCode != http.StatusOK {
		return err
	}

	_, err = RequestPostUserCategory(userToken, NewRequestCategory(req.Name, req.ColorId))
	if err != nil {
		return err
	}

	return nil
}

// Zo
func (s *service) GetMypageZos(userToken *UserToken) (*MypageZosGetModel, error) {
	resUser, err := RequestGetUser(userToken)
	if err != nil || resUser.StatusCode != http.StatusOK {
		return nil, err
	}

	resZo, err := RequestGetAllZo(userToken)
	if err != nil || resZo.StatusCode != http.StatusOK {
		return nil, err
	}

	// TODO 後で調整
	resCategory := Categories{}
	result := NewMypageZosGetModel(resUser.FamilyName+resUser.GivenName, resUser.Email, *resZo.Zos, resCategory, resZo.ResponseBase)
	return result, err
}

func (s *service) PostNewZo(userToken *UserToken, values url.Values) (*MyPageZosPostModel, error) {
	rz, err := validationAddZo(values)
	if err != nil {
		return nil, err
	}

	resZo, err := RequestPostZo(userToken, rz)
	if err != nil {
		return nil, err
	} else if resZo.ResponseBase.Error != nil {
		return nil, resZo.ResponseBase.Error
	}

	return resZo, nil
}

func validationAddCategory(values url.Values) (*requestUserCategory, error) {
	categoryName := values.Get("categoryName")
	if len(categoryName) > 20 {
		return nil, fmt.Errorf("メッセージは20文字以内で指定してください。")
	}

	return NewRequestCategory(categoryName, 0), nil
}

func validationAddZo(values url.Values) (*requestZo, error) {
	achievementDateStr := values.Get("achievementdate")
	achievementDate, err := time.Parse(TimeLayout, achievementDateStr)
	if err != nil {
		return nil, fmt.Errorf("達成日が適切な日付情報ではありません。achievementDate: %s", achievementDateStr)
	} else if achievementDate.Unix() > time.Now().Unix() {
		return nil, fmt.Errorf("達成日に未来が設定されています。achievementDate: %s", achievementDateStr)
	}

	expStr := values.Get("exp")
	exp, err := strconv.Atoi(expStr)
	if err != nil {
		return nil, fmt.Errorf("獲得経験値が正しい値ではありません。exp: %s", expStr)
	} else if exp < 0 || exp > 1000 {
		return nil, fmt.Errorf("獲得経験値は0〜1000の範囲で指定してください。exp: %s", expStr)
	}

	categoryIdStr := values.Get("categoryId")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return nil, fmt.Errorf("カテゴリー番号が正しい値ではありません。categoryNumber: %s", categoryIdStr)
	}

	message := values.Get("message")
	if len(message) > 30 {
		return nil, fmt.Errorf("メッセージは30文字以内で指定してください。")
	}

	return NewRequestZo(
		achievementDate,
		exp,
		categoryId,
		message,
	), nil
}
