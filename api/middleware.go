package api

import (
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) getAuthenticatedUser(ctx *gin.Context) dto.UserResponse {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	user, ok := profile.(dto.UserResponse)
	if !ok {
		return dto.UserResponse{}
	}

	return user
}

func (s *Server) isAuthenticated(ctx *gin.Context) {
	user := s.getAuthenticatedUser(ctx)
	if user == (dto.UserResponse{}) {
		err := util.Errorf(util.ErrorUnauthorized, "Unauthorized")
		ctx.AbortWithStatusJSON(errorResponse(err))
		return
	}

	ctx.Next()
}

func (s *Server) getUser(ctx *gin.Context) {
	user := s.getAuthenticatedUser(ctx)
	ctx.JSON(http.StatusOK, user)
}
