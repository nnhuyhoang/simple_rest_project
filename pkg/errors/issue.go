package errors

import "net/http"

// Issue errors
var (
	ErrIssueNotFound        = NewStringError("issue not found", http.StatusNotFound)
	ErrIssueDeleteCompleted = NewStringError("cannot delete completed issue", http.StatusBadRequest)
)
