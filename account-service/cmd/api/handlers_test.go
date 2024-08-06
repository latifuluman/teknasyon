package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateAccount(t *testing.T) {

	postBody := map[string]interface{}{
		"user_id":      2,
		"account_name": "test",
		"account_type": "tl",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.CreateAccount)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected http.StatusAccepted but got %d", rr.Code)
	}
}
