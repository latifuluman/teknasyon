package main

import (
	"database/sql"
	"errors"
	"net/http"

	"user-service/data"

	"github.com/go-chi/chi/v5"
)

// CreateUser handles the creation of a new user.
func (app *Config) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
		Active    int    `json:"active"`
	}

	// Read and decode JSON request payload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Create user object from request payload
	user := data.User{
		Email:     requestPayload.Email,
		FirstName: requestPayload.FirstName,
		LastName:  requestPayload.LastName,
		Password:  requestPayload.Password,
		Active:    requestPayload.Active,
	}

	// Insert user into the repository
	id, err := app.Repo.InsertUser(user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.ID = id
	payload := jsonResponse{
		Error:   false,
		Message: "",
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Repo.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	jwtToken, err := CreateToken(user.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "",
		Data: struct {
			Token string `json:"token"`
		}{Token: jwtToken},
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		app.errorJSON(w, ErrInvalidUser, http.StatusBadRequest)
	}

	user, err := app.Repo.GetUserByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	if user == nil {
		app.errorJSON(w, ErrUserNotFound)
	}

	// Prepare the response payload with the retrieved accounts.
	payload := jsonResponse{
		Error:   false,
		Message: "",
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		app.errorJSON(w, ErrInvalidUser, http.StatusBadRequest)
		return
	}

	user, err := app.Repo.GetUserByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			app.errorJSON(w, ErrUserNotFound, http.StatusAccepted)

		} else {

			app.errorJSON(w, err)
		}
		return
	}
	if user == nil {
		app.errorJSON(w, ErrUserNotFound)
	}

	err = app.Repo.DeleteUserByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Prepare the response payload with the retrieved accounts.
	payload := jsonResponse{
		Error:   false,
		Message: "",
		Data:    "Silme işlemi başarılı",
	}
	app.writeJSON(w, http.StatusAccepted, payload)

}
