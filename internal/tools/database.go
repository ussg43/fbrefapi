package tools

import(
	log "github.com/sirupsen/logrus"
)

type LoginDetails struct{
	AuthToken string
	Username string
}

type PlayerDetails struct{
	Name string
	Username string
}

type DatabaseInterface interface{
	GetUserLoginDetails(username string) *LoginDetails
	GetPlayerDetails(username string) *PlayerDetails
	SetUpDB() error
}

func NewDatabase() (*DatabaseInterface, error){
	var database DatabaseInterface = &mockDB{}

	var err error = database.SetUpDB()
	if err != nil{
		log.Error(err)
		return nil, err
	}
	return &database, nil
}