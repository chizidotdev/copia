package app

import (
	"github.com/chizidotdev/copia/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) listOrders(ctx *gin.Context) {
	user := s.getUser(ctx)

	items, err := s.OrderService.ListOrders(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, items)
}
