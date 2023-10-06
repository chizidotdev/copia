package app

import (
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) deleteOrder(ctx *gin.Context) {
	idParam := ctx.Params.ByName("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err.Error()))
		return
	}

	user := server.getUser(ctx)

	arg := dto.DeleteOrderParams{
		ID:        orderID,
		UserEmail: user.Email,
	}

	err = server.OrderService.DeleteOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, "Item deleted")
}
