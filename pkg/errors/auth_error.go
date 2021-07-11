package errors

import (
	"net/http"
)

// Product errors
var (
	ErrInvalidEmail             = NewStringError("email is invalid", http.StatusBadRequest)
	ErrInvalidPhoneNumber       = NewStringError("phone number is invalid", http.StatusBadRequest)
	ErrEmptyPassword            = NewStringError("password is empty", http.StatusBadRequest)
	ErrInvalidPassword          = NewStringError("password is invalid", http.StatusBadRequest)
	ErrPasswordNotMatch         = NewStringError("password is not match with re password", http.StatusBadRequest)
	ErrIncorrectEmailOrPassword = NewStringError("email or password is incorrect", http.StatusBadRequest)
	ErrGenJWTFailed             = NewStringError("cannot geenrate JWT token", http.StatusInternalServerError)
	ErrInvalidAccessToken       = NewStringError("invalid access token", http.StatusUnauthorized)
	ErrPermissionDenied         = NewStringError("permission denied", http.StatusForbidden)
)
