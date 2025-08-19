package main

import (
	"fbrefapi/internal/handlers"
	"fmt"
	"net/http"
	"os"


	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	log.SetReportCaller(true)
	r := chi.NewRouter()
	handlers.Handler(r)
	
	fmt.Println("Starting API Service...")

	err2 := http.ListenAndServe(port, r)
	if err2 != nil {
		log.Error(err2)
	}
}
