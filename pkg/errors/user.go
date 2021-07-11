package errors

import (
	"net/http"
)

// User errors
var (
	ErrUserNotFound            = NewStringError("user not found", http.StatusNotFound)
	ErrUserEmailAlreadyExisted = NewStringError("email is already existed", http.StatusBadRequest)
	ErrUserPhoneAlreadyExisted = NewStringError("phone number is already existed", http.StatusBadRequest)
	ErrUserRoleNotFound        = NewStringError("user role not found", http.StatusNotFound)
)
