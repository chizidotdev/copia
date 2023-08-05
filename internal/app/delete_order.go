package app

import (
	"net/http"

	"github.com/chizidotdev/copia/internal/dto"
	"github.com/chizidotdev/copia/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) deleteOrder(ctx *gin.Context) {
	idParam := ctx.Params.ByName("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	userEmail := ctx.Query("email")
	if userEmail == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse("email is required"))
		return
	}

	arg := dto.DeleteOrderParams{
		ID:        orderID,
		UserEmail: userEmail,
	}

	err = server.OrderService.DeleteOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, "Item deleted")
}
