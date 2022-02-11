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

func (s *service) GetMypage(userToken *UserToken) (*MypageGetModel, error) {
	resUser, err := RequestGetUser(userToken)
	if err != nil || resUser.StatusCode != http.StatusOK {
		return nil, err
	}

	resZo, err := RequestGetAllZo(userToken)
	if err != nil || resZo.StatusCode != http.StatusOK {
		return nil, err
	}

	result := NewMypageGetModel(resUser.FamilyName+resUser.GivenName, resUser.Email, *resZo.Zos, resZo.ResponseBase)
	return result, err
}

func (s *service) PostNewZo(userToken *UserToken, values url.Values) (*MypageGetModel, error) {
	rz, err := convertRequestZo(values.Get("achievementdate"), values.Get("exp"), values.Get("categoryId"), values.Get("content"))
	if err != nil {
		return nil, err
	}

	resZo, err := RequestPostZo(userToken, rz)
	if err != nil || resZo.StatusCode != http.StatusOK {
		return nil, err
	}

	return nil, nil
}

func convertRequestZo(_achievementDate string, _exp string, _categoryId string, message string,
) (*requestZo, error) {
	if achievementDate, err := time.Parse(TimeLayout, _achievementDate); err != nil {
		return nil, fmt.Errorf("達成日が正しい値ではありません。achievementDate: %s", _achievementDate)
	} else if exp, err := strconv.Atoi(_exp); err != nil {
		return nil, fmt.Errorf("獲得経験値が正しい値ではありません。exp: %s", _exp)
	} else if categoryId, err := strconv.Atoi(_categoryId); err != nil {
		return nil, fmt.Errorf("カテゴリーIDが正しい値ではありません。categoryId: %s", _categoryId)
	} else {
		return NewRequestZo(
			achievementDate,
			exp,
			categoryId,
			message,
		), nil
	}
}
