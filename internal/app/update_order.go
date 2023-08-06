package app

import (
	"net/http"

	"github.com/chizidotdev/copia/internal/dto"
	"github.com/chizidotdev/copia/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) updateOrder(ctx *gin.Context) {
	var req dto.Order
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request"))
		return
	}

	arg := dto.Order(req)

	item, err := server.OrderService.UpdateOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, item)
}
