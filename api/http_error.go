package api

import (
	"errors"
	"github.com/chizidotdev/copia/util"
	"net/http"
)

func errorResponse(err error) (code int, obj any) {
	code = http.StatusInternalServerError
	obj = err.Error()

	var customErr *util.ErrResponse
	if !errors.As(err, &customErr) {
		return code, obj
	}

	switch customErr.Code {
	case util.ErrorInternal:
		code = http.StatusInternalServerError
	case util.ErrorUnauthorized:
		code = http.StatusUnauthorized
	case util.ErrorNotFound:
		code = http.StatusNotFound
	case util.ErrorBadRequest:
		code = http.StatusBadRequest
	case util.ErrorForbidden:
		code = http.StatusForbidden
	}

	return code, customErr.Message
}
