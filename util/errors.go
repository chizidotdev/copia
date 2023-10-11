package util

import (
	"errors"
	"net/http"
)

var (
	ErrorBadRequest   = errors.New("bad request")
	ErrorNotFound     = errors.New("not found")
	ErrorInternal     = errors.New("internal error")
	ErrorUnauthorized = errors.New("unauthorized")
	ErrorForbidden    = errors.New("forbidden")
)

func ErrorResponse(err error) (code int, obj any) {
	if errors.Is(err, ErrorUnauthorized) {
		return http.StatusUnauthorized, ErrorMessage(err.Error())
	}
	if errors.Is(err, ErrorNotFound) {
		return http.StatusNotFound, ErrorMessage(err.Error())
	}
	if errors.Is(err, ErrorBadRequest) {
		return http.StatusBadRequest, ErrorMessage(err.Error())
	}
	if errors.Is(err, ErrorForbidden) {
		return http.StatusForbidden, ErrorMessage(err.Error())
	}

	return http.StatusInternalServerError, ErrorMessage(err.Error())
}

func ErrorMessage(message string) string {
	return message
}
