package main

import "errors"

var (
	errAuthFailed         = errors.New("auth_failed")
	errInsufficentBalance = errors.New("unsufficient_balance")
	errAccountNotFound    = errors.New("account_not_found")
)
