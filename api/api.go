package api

import( 
	"net/http"
	"encoding/json"
)

type PlayerParams struct {
	Name string
	Club string
	Season string
}

type PlayerResponse struct {
	Code int
	Name string
	Stats map[string]float64
}

type Error struct {
	Code int

	Message string
}

func writeError(w http.ResponseWriter, message string, code int){
	resp := Error{Code : code, Message: message}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

var(
	RequestErrHandler = func(w http.ResponseWriter, err error){
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrHandler = func(w http.ResponseWriter){
		writeError(w, "An Unexpected Error Occured", http.StatusInternalServerError)
	}
)
