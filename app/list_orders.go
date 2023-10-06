package app

import (
	"github.com/chizidotdev/copia/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) listOrders(ctx *gin.Context) {
	user := server.getUser(ctx)

	items, err := server.OrderService.ListOrders(ctx, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, items)
}
