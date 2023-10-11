package api

import (
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) updateOrder(ctx *gin.Context) {
	var req dto.Order
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorMessage("invalid request"))
		return
	}

	arg := dto.Order(req)

	item, err := s.OrderService.UpdateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, item)
}
