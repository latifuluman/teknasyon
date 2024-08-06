package main

import "errors"

type sessionKey string

const (
	SessionKey sessionKey = "session-key"
	webPort    string     = "80"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrUserNotFound = errors.New("user_not_found")
	ErrInvalidUser  = errors.New("invalid_user")
)
