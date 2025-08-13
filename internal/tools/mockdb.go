package tools

import (
	"time"
)

type mockDB struct{}

var mockLogindetails = map[string]LoginDetails{
	"unc" :{
		AuthToken: "123ABC",
		Username: "unc",
	},
	"joe":{
		AuthToken: "456DEF",
		Username: "joe",
	},
}

var mockPlayerDetails = map[string]PlayerDetails{
	"unc":{
		Name:"Saka",
		Username: "unc",
	},
	"joe":{
		Name: "Dembele",
		Username: "joe",
	},
}

func(d *mockDB) GetUserLoginDetails(username string) *LoginDetails{
	time.Sleep(time.Second * 1)

	var clientData = LoginDetails{}
	clientData, ok := mockLogindetails[username]
	if !ok{
		return nil
	}

	return &clientData
}

func( d *mockDB) GetPlayerDetails(username string) *PlayerDetails{
	time.Sleep(time.Second *1)

	var clientData = PlayerDetails{}

	clientData, ok := mockPlayerDetails[username]

	if !ok{
		return nil
	}

	return &clientData
}

func(d *mockDB) SetUpDB() error{
	return nil
}