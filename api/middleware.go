package api

import (
	"github.com/chizidotdev/copia/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) isAuthenticated(ctx *gin.Context) {
	if sessions.Default(ctx).Get("profile") == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorMessage("Unauthorized"))
	} else {
		ctx.Next()
	}
}
