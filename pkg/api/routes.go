package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router() *chi.Mux {
	handler := new(Handlers)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/login", func(r chi.Router) {
		r.Use(handler.VerifyCredentials)
		r.Post("/", handler.GenerateTokens)
	})
	r.Put("/", handler.Register)
	return r
}
