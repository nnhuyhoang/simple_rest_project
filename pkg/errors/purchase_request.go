package errors

import "net/http"

// PurchaseRequest errors
var (
	ErrNoOrderFound          = NewStringError("no spare part order found", http.StatusBadRequest)
	ErrOrderAlreadyCompleted = NewStringError("there is an order already completed", http.StatusBadRequest)
)
