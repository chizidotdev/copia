package user

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type loginUserRequest struct {
	Email string `json:"email" binding:"required"`
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var user loginUserRequest
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
	// randStr, err := util.GenerateRandString(5)
	// if err != nil {
	//   // handle error
	// }
	// authCode := strings.ToUpper(randStr)
	// // send authCode to user email

	userProfile, err := u.pgStore.GetUserByEmail(ctx, user.Email)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid credentials",
			Reason:    err.Error(),
		})
		return
	}

	session := sessions.Default(ctx)
	session.Set(profileKey, userProfile)
	if err := session.Save(); err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to save session",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    userProfile,
		Message: "User logged in successfully",
	})
}
