package main

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) GetMypageViewModel(userToken *UserToken) (*GetUserResponseData, error) {
	res, err := RequestGetUser(userToken)
	return res, err
}
