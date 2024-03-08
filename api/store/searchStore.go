package store

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type searchStoresReq struct {
	Query string `json:"query" binding:"required"`
}

func (s *StoreHandler) SearchStores(ctx *gin.Context) {
	var req searchStoresReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if err != nil {
			httpUtil.Success(ctx, &httpUtil.SuccessResponse{
				Code:    http.StatusOK,
				Data:    []repository.Store{},
				Message: "No stores found",
			})
			return
		}

		return
	}
	stores, err := s.pgStore.SearchStores(ctx, req.Query)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to search stores",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    stores,
		Message: "Stores found",
	})
}
