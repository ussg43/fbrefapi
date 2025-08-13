package middleware

import (
	"errors"
	"net/http"

	"fbrefapi/api"

	"fbrefapi/internal/tools"
	log "github.com/sirupsen/logrus"
)

var UnauthError = errors.New("invalid username or token")

func Authorization(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		var username = r.URL.Query().Get("username")
		var token = r.Header.Get("Authorization")
		var err error

		if username == "" || token == ""{
			log.Error(UnauthError)
			api.RequestErrHandler(w, UnauthError)
			return 
		}
		
		var database *tools.DatabaseInterface
		database, err = tools.NewDatabase()
		if err != nil{
			api.InternalErrHandler(w)
			return
		}
		var loginDetails *tools.LoginDetails
		loginDetails = (*database).GetUserLoginDetails(username)

		if (loginDetails == nil ||(token != (*loginDetails).AuthToken)){
			log.Error(UnauthError)
			api.RequestErrHandler(w, UnauthError)
			return 
		}

		next.ServeHTTP(w, r)
	})
}

