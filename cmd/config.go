package main

import (
	"os"
	"strconv"
)

type Config struct {
	ClientId     int
	ClientSecret string
	CsrfAuthKey  string
}

var Cfg *Config

func init() {
	Cfg = &Config{}
	id, err := strconv.Atoi(os.Getenv("Go_P_Zo_WebCli_ClientId"))
	if err != nil {
		panic("Go_P_Zo_WebCli_ClientId not found")
	}
	secret := os.Getenv("Go_P_Zo_WebCli_ClientSecret")
	if err != nil {
		panic("Go_P_Zo_WebCli_ClientSecret not found")
	}
	csrfAuthKey := os.Getenv("Go_P_Zo_WebCli_CsrfAuthKey")
	if err != nil {
		panic("Go_P_Zo_WebCli_CsrfAuthKey not found")
	}

	Cfg.ClientId = id
	Cfg.ClientSecret = secret
	Cfg.CsrfAuthKey = csrfAuthKey
}
