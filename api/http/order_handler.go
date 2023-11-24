package http

import (
	"net/http"

	"github.com/chizidotdev/copia/api/http/middleware"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/internal/app/usecases"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	IDParam := ctx.Param("id")
	orderID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid order ID"))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	err = o.OrderService.DeleteOrder(ctx, core.DeleteOrderRequest{
		ID:     orderID,
		UserID: user.ID,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, "Successfully deleted order.")
}

func (o *OrderHandler) updateOrder(ctx *gin.Context) {
	var req core.UpdateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	order, err := o.OrderService.UpdateOrder(ctx, core.UpdateOrderRequest{
		ID: req.ID,
		OrderRequest: core.OrderRequest{
			UserID:                user.ID,
			CustomerID:            req.CustomerID,
			Status:                req.Status,
			EstimatedDeliveryDate: req.EstimatedDeliveryDate,
			OrderDate:             req.OrderDate,
			OrderItems:            req.OrderItems,
		},
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) updateOrderStatus(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	orderID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid order ID"))
		return
	}

	var req core.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	order, err := o.OrderService.UpdateOrderStatus(ctx, core.UpdateOrderStatusRequest{
		ID:     orderID,
		Status: req.Status,
	})

	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) updateOrderItems(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	orderID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid order ID"))
		return
	}

	var req core.UpdateOrderItemsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload. "+err.Error()))
		return
	}

	order, err := o.OrderService.UpdateOrderItems(ctx, core.UpdateOrderItemsRequest{
		OrderID:    orderID,
		OrderItems: req.OrderItems,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (o *OrderHandler) deleteOrderItems(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	orderID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid order ID"))
		return
	}

	var req core.DeleteOrderItemsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	order, err := o.OrderService.DeleteOrderItems(ctx, core.DeleteOrderItemsRequest{
		OrderID:    orderID,
		OrderItems: req.OrderItems,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, order)
}
