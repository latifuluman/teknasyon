package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/v1/users/login", app.Login)
	mux.Post("/v1/users", app.CreateUser)
	mux.Get("/v1/users/{userID}", app.GetUser)
	mux.Delete("/v1/users/{userID}", app.DeleteUser)
	return mux
}
