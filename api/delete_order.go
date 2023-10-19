package api

import (
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) deleteOrder(ctx *gin.Context) {
	idParam := ctx.Params.ByName("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorMessage(err.Error()))
		return
	}

	//user := s.getUser(ctx)

	arg := dto.DeleteOrderParams{
		ID: orderID,
		//UserID: user.Email,
	}

	err = s.OrderService.DeleteOrder(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, "Item deleted")
}
