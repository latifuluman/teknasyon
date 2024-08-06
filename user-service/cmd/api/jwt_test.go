package main

import (
	"testing"
)

func TestCreateToken(t *testing.T) {
	userID := "123"
	token, err := CreateToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected token, got empty string")
	}
}

func TestVerifyToken(t *testing.T) {
	userID := "123"
	token, err := CreateToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	session, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if session == nil {
		t.Fatalf("expected session, got nil")
	}
	if session.UserID != userID {
		t.Fatalf("expected userID %s, got %s", userID, session.UserID)
	}

	// Test invalid token
	invalidToken := "invalid.token.string"
	_, err = VerifyToken(invalidToken)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

}
