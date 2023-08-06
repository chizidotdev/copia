package app

import (
	"net/http"

	"github.com/chizidotdev/copia/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) listOrders(ctx *gin.Context) {
	user := server.getUser(ctx)

	items, err := server.OrderService.ListOrders(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, items)
}
