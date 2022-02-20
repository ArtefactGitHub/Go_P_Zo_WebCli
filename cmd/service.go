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

	result := NewMypageUserGetModel(res.FamilyName+res.GivenName, res.Email, res.ResponseBase)
	return result, err
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

	resCategory, err := RequestGetAllCategory(userToken)
	if err != nil || resZo.StatusCode != http.StatusOK {
		return nil, err
	}

	result := NewMypageZosGetModel(resUser.FamilyName+resUser.GivenName, resUser.Email, *resZo.Zos, resCategory.Categories, resZo.ResponseBase)
	return result, err
}

func (s *service) PostNewZo(userToken *UserToken, values url.Values) (*MypageGetModel, error) {
	rz, err := convertRequestZo(values.Get("achievementdate"), values.Get("exp"), values.Get("categoryId"), values.Get("message"))
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
