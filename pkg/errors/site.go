package errors

import (
	"net/http"
)

// Site errors
var (
	ErrSiteNotFound = NewStringError("site not found", http.StatusNotFound)
)
