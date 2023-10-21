package api

import (
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleCreateOrder(ctx *gin.Context) {
	var req dto.Order
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errMsg := util.Errorf(util.ErrorBadRequest, "Invalid request payload")
		ctx.JSON(errorResponse(errMsg))
		return
	}

	user := s.getAuthenticatedUser(ctx)
	order, err := s.OrderService.CreateOrder(ctx, dto.Order{
		UserID:                user.ID,
		CustomerID:            req.CustomerID,
		Status:                req.Status,
		EstimatedDeliveryDate: req.EstimatedDeliveryDate,
		OrderDate:             req.OrderDate,
		OrderItems:            req.OrderItems,
	})
	if err != nil {
		ctx.JSON(errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (s *Server) handleGetOrderByID(ctx *gin.Context) {
	idParam := ctx.Params.ByName("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		errMsg := util.Errorf(util.ErrorBadRequest, "Invalid order ID")
		ctx.JSON(errorResponse(errMsg))
		return
	}

	order, err := s.OrderService.GetOrderByID(ctx, orderID)
	if err != nil {
		ctx.JSON(errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (s *Server) updateOrder(ctx *gin.Context) {
	idParam := ctx.Params.ByName("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		errMsg := util.Errorf(util.ErrorBadRequest, "Invalid order ID")
		ctx.JSON(errorResponse(errMsg))
		return
	}

	var req dto.Order
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errMsg := util.Errorf(util.ErrorBadRequest, "Invalid request payload")
		ctx.JSON(errorResponse(errMsg))
		return
	}

	user := s.getAuthenticatedUser(ctx)
	order, err := s.OrderService.UpdateOrder(ctx, dto.Order{
		ID:                    orderID,
		UserID:                user.ID,
		CustomerID:            req.CustomerID,
		Status:                req.Status,
		EstimatedDeliveryDate: req.EstimatedDeliveryDate,
		OrderDate:             req.OrderDate,
		OrderItems:            req.OrderItems,
	})
	if err != nil {
		ctx.JSON(errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}
