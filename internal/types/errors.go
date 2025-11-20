package types

import "net/http"

type StatusError struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	HTTPCode int    `json:"-"`
}

func (e *StatusError) Error() string {
	return e.Message
}

func NewValidationError(message string) *StatusError {
	return &StatusError{
		Message:  message,
		Code:     "error_processing_request",
		HTTPCode: http.StatusBadRequest,
	}
}

func NewInternalServerError() *StatusError {
	return &StatusError{
		Message:  "An unknown error has occurred",
		Code:     "internal_server_error",
		HTTPCode: http.StatusInternalServerError,
	}
}
