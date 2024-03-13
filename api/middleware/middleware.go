package middleware

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	ProfileKey = "profile"
)

func GetAuthenticatedUser(ctx *gin.Context) repository.GetUserRow {
	session := sessions.Default(ctx)
	profile := session.Get(ProfileKey)
	user, ok := profile.(repository.GetUserRow)
	if !ok {
		return repository.GetUserRow{}
	}

	return user
}

func IsAuthenticated(ctx *gin.Context) {
	user := GetAuthenticatedUser(ctx)
	if user == (repository.GetUserRow{}) {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusForbidden,
			MessageID: "",
			Message:   "User is not authenticated",
			Reason:    "User is not authenticated",
		})
		return
	}

	ctx.Next()
}

func IsVendor(ctx *gin.Context) {
	user := GetAuthenticatedUser(ctx)
	if user == (repository.GetUserRow{}) {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusForbidden,
			MessageID: "",
			Message:   "User is not authenticated",
			Reason:    "User is not authenticated",
		})
		return
	}

	if user.Role != repository.UserRoleVendor {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusForbidden,
			MessageID: "",
			Message:   "User is not a vendor",
			Reason:    "User is not a vendor",
		})
		return
	}

	ctx.Next()
}
