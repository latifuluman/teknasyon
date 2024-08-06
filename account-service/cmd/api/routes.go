package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var (
	SessionKey = "session-key"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(app.Authorization(), middleware.Heartbeat("/ping"))

	mux.Post("/v1/accounts/transfer/money", app.TransferMoney)
	mux.Post("/v1/accounts", app.CreateAccount)
	mux.Get("/v1/accounts", app.ListAccounts)
	mux.Get("/v1/accounts/{accountID}", app.GetAccount)
	return mux
}

// Authorization checks the token and decides that is authorized or not
func (app *Config) Authorization() func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				authToken  string
				splitToken []string
				token      string
				err        error
			)

			authToken = r.Header.Get("Authorization")
			if authToken == "" {
				app.errorJSON(w, errors.New("invalid_token"), http.StatusUnauthorized)
				return
			}

			splitToken = strings.Split(authToken, "Bearer ")
			if len(splitToken) == 2 {

				token = splitToken[1]
			} else {

				app.errorJSON(w, errors.New("invalid_token"), http.StatusUnauthorized)
				return
			}
			session, err := VerifyToken(token)
			if err != nil {

				app.errorJSON(w, err, http.StatusUnauthorized)
				return
			}

			// Store session in context
			ctx := context.WithValue(r.Context(), SessionKey, session)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
