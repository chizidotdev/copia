package user

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Email string `json:"email" binding:"required"`
}

func (u *UserHandler) CreateUser(ctx *gin.Context) {
	var user createUserRequest
	err := ctx.BindJSON(&user)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid email",
			Reason:    err.Error(),
		})
		return
	}

	// TODO: Add validation for email

	userProfile, err := u.pgStore.CreateUser(ctx, repository.CreateUserParams{
		FirstName: "",
		LastName:  "",
		Email:     user.Email,
		GoogleID:  "",
		Image:     "",
		Role:      repository.UserRoleCustomer,
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to create user",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusCreated,
		Data:    userProfile,
		Message: "User created successfully",
	})
}
