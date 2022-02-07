package main

type MypageGetModel struct {
	*ResponseBase
	Name  string
	Email string
	Zos   Zos
}

func NewMypageGetModel(name, email string, zos Zos, base *ResponseBase) *MypageGetModel {
	return &MypageGetModel{Name: name, Email: email, Zos: zos, ResponseBase: base}
}
