package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateUser(t *testing.T) {

	postBody := map[string]interface{}{
		"email":    "me@here.com",
		"password": "verysecret",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.CreateUser)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected http.StatusAccepted but got %d", rr.Code)
	}
}

func TestLogin(t *testing.T) {

	payload := `{"email":"test@example.com","password":"password"}`
	req, _ := http.NewRequest("POST", "/v1/user/login", bytes.NewBufferString(payload))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testApp.Login)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}

}
