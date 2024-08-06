package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Session struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
}

var (
	secretKey = []byte(os.Getenv("JWT_SECRET"))
)

// CreateToken takes user id and generates a jwt token that contains it
func CreateToken(userID string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

// VerifyToken takes a token and verifies it. If verified, the session belongs to the token is returned
func VerifyToken(tokenString string) (*Session, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if userID, ok := claims["user_id"].(string); ok {
			s := &Session{
				SessionID: tokenString,
				UserID:    userID,
			}
			return s, nil
		}

	}
	return nil, fmt.Errorf("invalid token")
}
