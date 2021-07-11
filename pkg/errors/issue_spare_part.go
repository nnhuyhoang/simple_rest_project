package errors

import "net/http"

// IssueSparePart errors
var (
	ErrIssueSparePartNotFound = NewStringError("issue spare part not found", http.StatusNotFound)
)
