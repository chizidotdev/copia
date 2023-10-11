package api

import (
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) createUser(ctx *gin.Context) {
	var req dto.CreateUserParams
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err.Error()))
		return
	}

	user, err := s.UserService.CreateUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
