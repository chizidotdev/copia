package user

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *UserHandler) Logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Delete(profileKey)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to update session",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    nil,
		Message: "Log out successful.",
	})
}
