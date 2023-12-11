package http

import (
	"fmt"
	"github.com/chizidotdev/copia/api/http/middleware"
	"github.com/chizidotdev/copia/config"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/internal/app/usecases"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	profileKey = middleware.ProfileKey
	stateKey   = middleware.StateKey
)

type UserHandler struct {
	UserService *usecases.UserService
}

func NewUserHandler(userService *usecases.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (u *UserHandler) createUser(ctx *gin.Context) {
	var req core.CreateUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errResp := invalidRequestError(err)
		errorResponse(ctx, errResp)
		return
	}

	user, err := u.UserService.CreateUser(ctx, req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	successResponse(ctx, http.StatusCreated, SuccessResponse{
		Data:    user,
		Message: "User created successfully.",
	})
}

func (u *UserHandler) login(ctx *gin.Context) {
	var req core.LoginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errResp := invalidRequestError(err)
		errorResponse(ctx, errResp)
		return
	}

	user, err := u.UserService.GetUser(ctx, req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	session := sessions.Default(ctx)
	session.Set(profileKey, user)
	if err := session.Save(); err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to update session",
			Reason:    err.Error(),
		}
		errorResponse(ctx, errors.Errorf(errResp))
		return
	}

	successResponse(ctx, http.StatusOK, SuccessResponse{
		Data:    user,
		Message: "Login successful.",
	})
}

func (u *UserHandler) loginWithSSO(ctx *gin.Context) {
	errRedirectURL := config.EnvVars.AuthDomain + "/u/login/errors"
	state, err := u.UserService.GenerateAuthState()
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}
	session := sessions.Default(ctx)
	session.Set(stateKey, state)
	if err := session.Save(); err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}

	googleAuthConfig := u.UserService.GetGoogleAuthConfig()
	url := googleAuthConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (u *UserHandler) ssoCallback(ctx *gin.Context) {
	errRedirectURL := config.EnvVars.AuthDomain + "/u/login/error"
	successRedirectURL := config.EnvVars.AuthDomain + "/u/login/success"

	session := sessions.Default(ctx)
	if ctx.Query(stateKey) != session.Get(stateKey) {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s?errors=invalid_state", errRedirectURL))
		return
	}

	code := ctx.Query("code")
	userProfile, err := u.UserService.GoogleCallback(ctx, code)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s?errors=failed_to_exchange", errRedirectURL))
		return
	}

	session.Set(profileKey, userProfile)
	if err := session.Save(); err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}

	ctx.Redirect(http.StatusPermanentRedirect, successRedirectURL)
}

func (u *UserHandler) logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Set(profileKey, nil)
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to update session",
			Reason:    err.Error(),
		}
		errorResponse(ctx, errors.Errorf(errResp))
		return
	}

	successResponse(ctx, http.StatusOK, SuccessResponse{
		Data:    nil,
		Message: "Log out successful.",
	})
}

func (u *UserHandler) getUser(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)
	successResponse(ctx, http.StatusOK, SuccessResponse{
		Data:    user,
		Message: "User retrieved successfully.",
	})
}

func (u *UserHandler) sendVerificationEmail(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errResp := invalidRequestError(err)
		errorResponse(ctx, errResp)
		return
	}

	err = u.UserService.SendVerificationEmail(ctx, req.Email)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	successResponse(ctx, http.StatusOK, SuccessResponse{
		Data:    nil,
		Message: "Verification email sent.",
	})
}

func (u *UserHandler) verifyEmail(ctx *gin.Context) {
	var req core.VerifyEmailRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errResp := invalidRequestError(err)
		errorResponse(ctx, errResp)
		return
	}

	verifiedUser, err := u.UserService.VerifyEmail(ctx, req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	authenticatedUser := middleware.GetAuthenticatedUser(ctx)
	// if the user is authenticated, update the session
	if authenticatedUser.Email == verifiedUser.Email {
		session := sessions.Default(ctx)
		session.Set(profileKey, core.UserResponse{
			ID:            authenticatedUser.ID,
			FirstName:     authenticatedUser.FirstName,
			LastName:      authenticatedUser.LastName,
			Email:         authenticatedUser.Email,
			EmailVerified: verifiedUser.EmailVerified,
		})
		if err := session.Save(); err != nil {
			errResp := errors.ErrResponse{
				Code:      errors.ErrorInternal,
				MessageID: "",
				Message:   "Failed to update session",
				Reason:    err.Error(),
			}
			errorResponse(ctx, errors.Errorf(errResp))
			return
		}
	}

	successResponse(ctx, http.StatusOK, SuccessResponse{
		Data:    authenticatedUser,
		Message: "Email verification successful.",
	})
}

func (u *UserHandler) resetPassword(ctx *gin.Context) {
	var req core.ResetPasswordRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errResp := invalidRequestError(err)
		errorResponse(ctx, errResp)
		return
	}

	err = u.UserService.ResetPassword(ctx, req.Email)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	successResponse(ctx, http.StatusOK, SuccessResponse{
		Data:    nil,
		Message: "Password reset email sent.",
	})
}

func (u *UserHandler) changePassword(ctx *gin.Context) {
	var req core.ChangePasswordRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errResp := invalidRequestError(err)
		errorResponse(ctx, errResp)
		return
	}

	err = u.UserService.ChangePassword(ctx, req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	successResponse(ctx, http.StatusOK, SuccessResponse{
		Data:    nil,
		Message: "Password changed successfully.",
	})
}
