package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func testVerifyToken(token string) (*Session, error) {
	if token == "valid-token" {
		return &Session{
			UserID: "123",
		}, nil
	}
	return nil, errors.New("invalid_token")
}

func Test_CreateAccount(t *testing.T) {
	postBody := map[string]interface{}{
		"account_name":    "test",
		"account_type":    "tl",
		"initial_balance": 0.0,
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/v1/accounts", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer valid-token") // Add Authorization header

	rr := httptest.NewRecorder()

	handler := chi.NewRouter()

	handler.Use(testApp.Authorization(testVerifyToken))
	handler.Post("/v1/accounts", testApp.CreateAccount)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected http.StatusAccepted but got %d", rr.Code)
	}
}
