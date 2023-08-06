package app

import (
	"net/http"

	"github.com/chizidotdev/copia/internal/dto"
	"github.com/chizidotdev/copia/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) isAuth(ctx *gin.Context) {
	reqToken := ctx.Request.Header.Get("Authorization")
	if reqToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized"))
		return
	}

	user, err := server.TokenManager.Parse(reqToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}

func (server *Server) getUser(ctx *gin.Context) *dto.Claims {
	user := ctx.MustGet("user").(*dto.Claims)

	return user
}
