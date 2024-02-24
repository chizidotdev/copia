package user

import (
	"database/sql"
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
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorBadRequest,
			MessageID: "",
			Message:   "Invalid email",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	// TODO: Add validation for email

	userProfile, err := u.pgStore.CreateUser(ctx, repository.CreateUserParams{
		FirstName: "",
		LastName:  "",
		Email:     user.Email,
		GoogleID:  sql.NullString{String: "", Valid: true},
		Image:     "",
		Role:      userRoles["customer"],
	})
	if err != nil {
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorInternal,
			MessageID: "",
			Message:   "Failed to create user",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	httpUtil.Success(ctx, http.StatusCreated, httpUtil.SuccessResponse{
		Data:    userProfile,
		Message: "User created successfully",
	})
}
