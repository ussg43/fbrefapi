package handlers

import (
	"encoding/json"
	"net/http"

	"fbrefapi/api"
	"fbrefapi/internal/tools/scraper"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

func GetPlayerP90(w http.ResponseWriter, r *http.Request) {
	var params = api.PlayerParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error(err)
		api.InternalErrHandler(w)
		return
	}

	// var database *tools.DatabaseInterface
	// database, err = tools.NewDatabase()
	// if err != nil{
	// 	api.InternalErrHandler(w)
	// 	return
	// }

	var tokenDetails *scraper.Player
	tokenDetails = scraper.NewPlayer(params.Name, params.Club)
	if tokenDetails == nil {
		api.InternalErrHandler(w)
		return
	}

	var response = api.PlayerResponse{
		Code:  http.StatusOK,
		Name:  (*tokenDetails).Name,
		Stats: (*tokenDetails).GetP90(),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error(err)
		api.InternalErrHandler(w)
		return
	}
}

func GetPlayerSeasonal(w http.ResponseWriter, r *http.Request) {
	var params = api.PlayerParams{}
	var decoder *schema.Decoder = schema.NewDecoder()
	var err error

	err = decoder.Decode(&params, r.URL.Query())

	if err != nil {
		log.Error(err)
		api.InternalErrHandler(w)
		return
	}

	// var database *tools.DatabaseInterface
	// database, err = tools.NewDatabase()
	// if err != nil{
	// 	api.InternalErrHandler(w)
	// 	return
	// }

	var tokenDetails *scraper.Player
	tokenDetails = scraper.NewPlayer(params.Name, params.Club)
	if tokenDetails == nil {
		api.InternalErrHandler(w)
		return
	}

	var response = api.PlayerResponse{
		Code:  http.StatusOK,
		Name:  (*tokenDetails).Name,
		Stats: (*tokenDetails).GetSeasonal(params.Season),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panic("lol ur cooked")
		api.InternalErrHandler(w)
		return
	}
}
