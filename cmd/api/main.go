package main

import (
	"fbrefapi/internal/handlers"
	"fmt"
	"net/http"
	
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	r := chi.NewRouter()
	handlers.Handler(r)
	
	fmt.Println("Starting API Service...")

	err := http.ListenAndServe(":5656", r)
	if err != nil {
		log.Error(err)
	}
}
