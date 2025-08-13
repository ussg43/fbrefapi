package handlers

import (

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	// "fbrefapi/internal/middleware"
)

func Handler(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)

	r.Route("/player", func(r chi.Router) {
		// r.Use(middleware.Authorization)

		r.Get("/p90", GetPlayerP90)
		r.Get("/seasonal", GetPlayerSeasonal)
		
	})
}
