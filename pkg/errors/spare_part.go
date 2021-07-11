package errors

import (
	"net/http"
)

// SparePart errors
var (
	ErrSparePartNotFound = NewStringError("spare part not found", http.StatusNotFound)
)
