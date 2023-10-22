package middleware

import (
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAuthenticatedUser(ctx *gin.Context) core.UserResponse {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user, ok := profile.(core.UserResponse)
	if !ok {
		return core.UserResponse{}
	}

	return user
}

func IsAuthenticated(ctx *gin.Context) {
	user := GetAuthenticatedUser(ctx)
	if user == (core.UserResponse{}) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx.Next()
}
