package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadJSON(t *testing.T) {
	app := Config{}
	sampleJSON := `{"name": "test"}`
	var data struct {
		Name string `json:"name"`
	}
	r := httptest.NewRequest("POST", "/", bytes.NewBufferString(sampleJSON))
	w := httptest.NewRecorder()

	err := app.readJSON(w, r, &data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if data.Name != "test" {
		t.Fatalf("expected name to be 'test', got '%s'", data.Name)
	}
}

func TestWriteJSON(t *testing.T) {
	app := Config{}
	w := httptest.NewRecorder()
	data := jsonResponse{
		Error:   false,
		Message: "success",
		Data:    nil,
	}
	err := app.writeJSON(w, http.StatusOK, data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	res := w.Result()
	defer res.Body.Close()

	var responseData jsonResponse
	json.NewDecoder(res.Body).Decode(&responseData)

	if responseData.Error != data.Error || responseData.Message != data.Message {
		t.Errorf("expected %v, got %v", data, responseData)
	}
}

func TestErrorJSON(t *testing.T) {
	app := Config{}
	w := httptest.NewRecorder()
	expectedErr := errors.New("test error")

	e := app.errorJSON(w, expectedErr, http.StatusInternalServerError)
	if e != nil {
		t.Fatalf("expected no error, got %v", e)
	}

	res := w.Result()
	defer res.Body.Close()

	var responseData jsonResponse
	json.NewDecoder(res.Body).Decode(&responseData)

	if !responseData.Error || responseData.Message != expectedErr.Error() {
		t.Errorf("expected error message '%s', got '%s'", expectedErr.Error(), responseData.Message)
	}
}
