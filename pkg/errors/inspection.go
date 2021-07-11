package errors

import (
	"net/http"
)

// Inspection errors
var (
	ErrInspectionNotFound = NewStringError("inspection not found", http.StatusNotFound)
)
