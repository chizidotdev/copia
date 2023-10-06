package app

import (
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) updateOrder(ctx *gin.Context) {
	var req dto.Order
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse("invalid request"))
		return
	}

	arg := dto.Order(req)

	item, err := server.OrderService.UpdateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, item)
}
