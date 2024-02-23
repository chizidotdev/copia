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
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorInternal,
			MessageID: "",
			Message:   "Failed to update session",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	httpUtil.Success(ctx, http.StatusOK, httpUtil.SuccessResponse{
		Data:    nil,
		Message: "Log out successful.",
	})
}
