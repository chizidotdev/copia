package seed

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

const (
	storeIDParam = "storeID"
)

type SeedHandler struct {
	pgStore *repository.Repository
}

func NewSeedHandler(pgStore *repository.Repository) *SeedHandler {
	return &SeedHandler{
		pgStore: pgStore,
	}
}

func (s *SeedHandler) SeedStore(ctx *gin.Context) {
	storeID, err := repository.ParseUUID(ctx.Param(storeIDParam))
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid store id",
			Reason:    err.Error(),
		})
		return
	}

	store, err := s.pgStore.GetStore(ctx, storeID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusNotFound,
			MessageID: "",
			Message:   "Failed to retrieve store",
			Reason:    err.Error(),
		})
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	if user.ID != store.UserID {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusForbidden,
			MessageID: "",
			Message:   "User does not have permission to seed store",
			Reason:    "",
		})
		return
	}

	err = s.pgStore.ExecTx(ctx, func(tx *repository.Queries) error {
		for i := 0; i < 10; i++ {
			product := seedProduct()
			var txErr error
			_, txErr = s.pgStore.CreateProduct(ctx, repository.CreateProductParams{
				StoreID:     storeID,
				Title:       product.Title,
				Description: product.Description,
				Price:       product.Price,
				OutOfStock:  product.OutOfStock,
			})
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to create product",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Database seeded successfully",
	})
}

type seedUserResult struct {
	Email     string `fake:"{email}"`
	FirstName string `fake:"{firstName}"`
	LastName  string `fake:"{lastName}"`
	Image     string `fake:"{imageURL}"`
	GoogleID  string `fake:"{str}"`
}

func seedUser() seedUserResult {
	var u seedUserResult
	gofakeit.Struct(&u)
	return u
}

type seedStoreResult struct {
	Name        string `fake:"{company}"`
	Description string `fake:"{sentence:3}"`
	Image       string `fake:"{imageURL}"`
}

func seedStore() seedStoreResult {
	var s seedStoreResult
	gofakeit.Struct(&s)
	return s
}

type seedProductResult struct {
	Title       string  `fake:"{sentence:2}"`
	Description string  `fake:"{sentence:5}"`
	Price       float64 `fake:"{price}"`
	OutOfStock  bool    `fake:"{bool}"`
}

func seedProduct() seedProductResult {
	var p seedProductResult
	gofakeit.Struct(&p)
	return p
}
