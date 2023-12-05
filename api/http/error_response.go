package http

import (
	errors2 "errors"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
)

func errorResponse(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError
	obj := err.Error()

	log.Println(obj)
	var customErr *errors.ErrResponse
	if !errors2.As(err, &customErr) {
		re := regexp.MustCompile(`not found.?`)
		if re.FindString(err.Error()) != "" {
			ctx.AbortWithStatusJSON(http.StatusNotFound, err.Error())
			return
		}

		ctx.AbortWithStatusJSON(code, obj)
		return
	}

	switch customErr.Code {
	case errors.ErrorInternal:
		code = http.StatusInternalServerError
	case errors.ErrorUnauthorized:
		code = http.StatusUnauthorized
	case errors.ErrorNotFound:
		code = http.StatusNotFound
	case errors.ErrorBadRequest:
		code = http.StatusBadRequest
	case errors.ErrorForbidden:
		code = http.StatusForbidden
	}

	ctx.AbortWithStatusJSON(code, customErr.Message)
}
