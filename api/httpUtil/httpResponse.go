package httpUtil

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Error   *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Reason    string `json:"reason"`
}

func Error(ctx *gin.Context, err error) {
	code := http.StatusInternalServerError

	var customErr *HttpError
	if !errors.As(err, &customErr) {
		resp := &ErrorResponse{
			Code:      code,
			MessageID: "",
			Message:   "Internal Server Error",
			Reason:    err.Error(),
		}

		ctx.AbortWithStatusJSON(code, &Response{
			Data:    nil,
			Error:   resp,
			Message: resp.Message,
		})
		return
	}

	switch customErr.Code {
	case ErrorInternal:
		code = http.StatusInternalServerError
	case ErrorUnauthorized:
		code = http.StatusUnauthorized
	case ErrorNotFound:
		code = http.StatusNotFound
	case ErrorBadRequest:
		code = http.StatusBadRequest
	case ErrorForbidden:
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
		Data:    nil,
		Error:   resp,
		Message: resp.Message,
	})
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success(ctx *gin.Context, code int, succResp SuccessResponse) {
	ctx.JSON(code, &Response{
		Data:    succResp.Data,
		Error:   nil,
		Message: succResp.Message,
	})
}

func invalidRequestError(err error) *HttpError {
	return Errorf(HttpError{
		Code:      ErrorBadRequest,
		MessageID: "",
		Message:   "Invalid request payload.",
		Reason:    err.Error(),
	})
}
