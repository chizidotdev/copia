package api

import (
	"github.com/chizidotdev/copia/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) getOrderByID(ctx *gin.Context) {
	idParam := ctx.Params.ByName("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorMessage(err.Error()))
		return
	}

	item, err := s.OrderService.GetOrderByID(ctx, orderID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, util.ErrorMessage(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, item)
}
