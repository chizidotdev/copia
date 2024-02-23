package middleware

import (
	"net/http"

	"github.com/chizidotdev/shop/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	ProfileKey = "profile"
	StateKey   = "state"
)

func GetAuthenticatedUser(ctx *gin.Context) repository.User {
	session := sessions.Default(ctx)
	profile := session.Get(ProfileKey)
	user, ok := profile.(repository.User)
	if !ok {
		return repository.User{}
	}

	return user
}

func IsAuthenticated(ctx *gin.Context) {
	user := GetAuthenticatedUser(ctx)
	if user == (repository.User{}) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx.Next()
}
