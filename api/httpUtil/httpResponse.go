package httpUtil

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Error   *ErrorResponse `json:"error"`
	Code    int            `json:"code"`
}

type ErrorResponse struct {
	MessageID string `json:"message_id"`
	Message   string `json:"message"`
	Code      int    `json:"code"`
	Reason    string `json:"reason"`
}

func Error(ctx *gin.Context, errResp *ErrorResponse) {
	errorLog := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	errorLog.Println(errResp)
	ctx.AbortWithStatusJSON(errResp.Code, &Response{
		Data:    nil,
		Error:   errResp,
		Message: errResp.Message,
		Code:    errResp.Code,
	})
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

func Success(ctx *gin.Context, succResp *SuccessResponse) {
	ctx.JSON(succResp.Code, &Response{
		Data:    succResp.Data,
		Error:   nil,
		Message: succResp.Message,
		Code:    succResp.Code,
	})
}
