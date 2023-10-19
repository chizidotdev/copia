package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) listOrders(ctx *gin.Context) {
	user := s.getAuthenticatedUser(ctx)

	orders, err := s.OrderService.ListOrders(ctx, user.ID)
	if err != nil {
		errorResponse(err)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
