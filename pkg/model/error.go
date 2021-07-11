package model

import "bytes"

// Code describe an Code the error.
type Code int

// Message human can read message
type Message string

//Error ...
type Error struct {
	Code    Code    `json:"Code,omitempty"`
	Message Message `json:"Message,omitempty"`
}

func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
func (e Error) Error() string {
	b := new(bytes.Buffer)

	if e.Message != "" {
		pad(b, ": ")
		b.WriteString(string(e.Message))
	}

	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// StatusCode status code
func (e *Error) StatusCode() int {
	return int(e.Code)
}
