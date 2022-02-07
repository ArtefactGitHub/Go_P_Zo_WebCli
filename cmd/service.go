package main

import "net/http"

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) GetMypageViewModel(userToken *UserToken) (*MypageGetModel, error) {
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
