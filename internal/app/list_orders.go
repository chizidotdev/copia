package app

import (
	"net/http"

	"github.com/chizidotdev/copia/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (server *Server) listOrders(ctx *gin.Context) {
	userEmail := ctx.Query("email")
	if userEmail == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("email is required"))
		return
	}

	items, err := server.OrderService.ListOrders(ctx, userEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, items)
}
