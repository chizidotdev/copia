package http

import (
	"github.com/chizidotdev/copia/api/http/middleware"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/internal/app/usecases"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type OrderHandler struct {
	OrderService *usecases.OrderService
}

func NewOrderHandler(orderService *usecases.OrderService) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}

func (o *OrderHandler) createOrder(ctx *gin.Context) {
	var req core.OrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	order, err := o.OrderService.CreateOrder(ctx, core.OrderRequest{
		UserID:                user.ID,
		CustomerID:            req.CustomerID,
		Status:                req.Status,
		EstimatedDeliveryDate: req.EstimatedDeliveryDate,
		OrderDate:             req.OrderDate,
		OrderItems:            req.OrderItems,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, order)
}

func (o *OrderHandler) listOrders(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)
	orders, err := o.OrderService.ListOrders(ctx, user.ID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (o *OrderHandler) getOrder(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	orderID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid order ID"))
		return
	}

	order, err := o.OrderService.GetOrderByID(ctx, orderID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) deleteOrder(ctx *gin.Context) {
	var req core.DeleteOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	req.UserID = user.ID
	err := o.OrderService.DeleteOrder(ctx, req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, "Successfully deleted order.")
}