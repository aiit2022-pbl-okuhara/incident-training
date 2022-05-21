package responses

import (
	"github.com/google/go-safeweb/safehttp"
	"github.com/google/safehtml"
)

// Error is a safe error response.
type Error struct {
	StatusCode   safehttp.StatusCode
	ErrorMessage safehtml.HTML
}

// NewError creates a new error response.
func NewError(code safehttp.StatusCode, message safehtml.HTML) Error {
	return Error{
		StatusCode:   code,
		ErrorMessage: message,
	}
}

// Code returns the HTTP response code.
func (e Error) Code() safehttp.StatusCode {
	return e.StatusCode
}

// Message returns the HTTP response message.
func (e Error) Message() safehtml.HTML {
	return e.ErrorMessage
}
