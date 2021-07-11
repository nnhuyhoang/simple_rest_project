package errors

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
)

// E builds an error value from its arguments.
func E(args ...interface{}) *model.Error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}

	e := &model.Error{}
	for _, arg := range args {
		switch arg := arg.(type) {

		case model.Message:
			e.Message = arg
		case model.Code:
			e.Code = arg

		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			e.Message = model.Message(fmt.Sprintf("unknown type %T, value %v in error call", arg, arg))
		}
	}
	return e
}

// NewStringError new a new string err
func NewStringError(err string, statusCode int) error {

	return E(model.Message(err), model.Code(statusCode))
}

var (
	ErrInternalServerError   = NewStringError("internal server error", http.StatusInternalServerError)
	ErrCannotConnectDatabase = NewStringError("cannot connect to databse", http.StatusInternalServerError)
	ErrInvalidIdFormat       = NewStringError("invalid id format", http.StatusBadRequest)
	ErrUnexpectedHeader      = NewStringError("unexpected headers", http.StatusUnauthorized)
	ErrInvalidAuthenType     = NewStringError("header Authentication Type is invalid", http.StatusUnauthorized)
	ErrInvalidDateFormat     = NewStringError("invalid date format", http.StatusBadRequest)
	ErrInvalidTimeZone       = NewStringError("invalid timezone", http.StatusInternalServerError)
)
