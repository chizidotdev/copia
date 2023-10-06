package app

import (
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) createOrder(ctx *gin.Context) {
	var req dto.Order
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err.Error()))
		return
	}

	user := s.getUser(ctx)
	item, err := s.OrderService.CreateOrder(ctx, dto.Order{
		UserEmail:             user.Email,
		CustomerID:            req.CustomerID,
		Status:                req.Status,
		ShippingDetails:       req.ShippingDetails,
		EstimatedDeliveryDate: req.EstimatedDeliveryDate,
		OrderDate:             req.OrderDate,
		TotalAmount:           req.TotalAmount,
		PaymentStatus:         req.PaymentStatus,
		PaymentMethod:         req.PaymentMethod,
		BillingAddress:        req.BillingAddress,
		ShippingAddress:       req.ShippingAddress,
		Notes:                 req.Notes,
		OrderItems:            req.OrderItems,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, item)
}

/*
 Example request body:
 {
	"user_email": "user@email.com",
	"customer_id": "b3f9a7a0-9f9a-4e4a-8b1a-9b2b1c9c0a0a",
	"status": "pending",
	"shipping_details": "UPS",
	"estimated_delivery_date": "2021-08-23T00:00:00Z",
	"order_date": "2021-07-23T00:00:00Z",
	"total_amount": 100,
	"payment_status": "pending",
	"payment_method": "credit",
	"billing_address": "123 Main St.",
	"shipping_address": "456 Main St.",
	"notes": "some notes",
	"order_items": [
		{
			"product_id": "b3f9a7a0-9f9a-4e4a-8b1a-9b2b1c9c0a0a",
			"quantity": 1,
			"unit_price": 100,
			"sub_total": 100
		}
	]
 }
*/
