package http

import (
	errors2 "errors"
	"fmt"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Data  interface{}    `json:"data"`
	Error *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Reason    string `json:"reason"`
}

func errorResponse(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError

	var customErr *errors.ErrResponse
	if !errors2.As(err, &customErr) {
		resp := &ErrorResponse{
			Code:      code,
			MessageID: "",
			Message:   "Internal Server Error",
			Reason:    err.Error(),
		}

		ctx.AbortWithStatusJSON(code, &Response{
			Data:  nil,
			Error: resp,
		})
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

	var message string
	if customErr.MessageID != "" {
		message = fmt.Sprintf("%s: %s", customErr.MessageID, customErr.Message)
	} else {
		message = customErr.Message
	}

	resp := &ErrorResponse{
		Code:      code,
		MessageID: customErr.MessageID,
		Message:   message,
		Reason:    customErr.Reason,
	}
	ctx.AbortWithStatusJSON(code, &Response{
		Data:  nil,
		Error: resp,
	})
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func successResponse(ctx *gin.Context, code int, succResp SuccessResponse) {
	ctx.JSON(code, &Response{
		Data:  succResp,
		Error: nil,
	})
}

func invalidRequestError(err error) *errors.ErrResponse {
	return errors.Errorf(errors.ErrResponse{
		Code:      errors.ErrorBadRequest,
		MessageID: "",
		Message:   "Invalid request payload.",
		Reason:    err.Error(),
	})
}
