package apperror

import (
	"errors"
	"net/http"
)

type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Fields  any    `json:"fields,omitempty"`
}

func (e *Error) Error() string { return e.Message }

func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

func BadRequest(msg string) *Error   { return New(http.StatusBadRequest, "BAD_REQUEST", msg) }
func Unauthorized(msg string) *Error { return New(http.StatusUnauthorized, "UNAUTHORIZED", msg) }
func Forbidden(msg string) *Error    { return New(http.StatusForbidden, "FORBIDDEN", msg) }
func NotFound(msg string) *Error     { return New(http.StatusNotFound, "NOT_FOUND", msg) }
func Conflict(msg string) *Error     { return New(http.StatusConflict, "CONFLICT", msg) }
func Internal(msg string) *Error     { return New(http.StatusInternalServerError, "INTERNAL", msg) }

func Validation(msg string, fields any) *Error {
	e := New(http.StatusUnprocessableEntity, "VALIDATION_ERROR", msg)
	e.Fields = fields
	return e
}

func As(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}
