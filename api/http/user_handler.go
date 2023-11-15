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
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user, err := u.UserService.CreateUser(ctx, req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (u *UserHandler) login(ctx *gin.Context) {
	var req core.LoginUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload"))
		return
	}

	user, err := u.UserService.GetUser(ctx, req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	session := sessions.Default(ctx)
	session.Set("profile", user)
	if err := session.Save(); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorInternal, "Failed to update session"))
		return
	}

	ctx.JSON(http.StatusOK, "Successfully logged in.")
}

func (u *UserHandler) loginWithSSO(ctx *gin.Context) {
	errRedirectURL := config.EnvVars.AuthDomain + "/u/login/errors"
	state, err := u.UserService.GenerateAuthState()
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}
	session := sessions.Default(ctx)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}

	googleAuthConfig := u.UserService.GetGoogleAuthConfig()
	url := googleAuthConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (u *UserHandler) ssoCallback(ctx *gin.Context) {
	errRedirectURL := config.EnvVars.AuthDomain + "/u/login/errors"
	successRedirectURL := config.EnvVars.AuthDomain + "/u/login/success"

	session := sessions.Default(ctx)
	if ctx.Query("state") != session.Get("state") {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s?errors=invalid_state", errRedirectURL))
		return
	}

	code := ctx.Query("code")
	userProfile, err := u.UserService.GoogleCallback(ctx, code)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s?errors=failed_to_exchange", errRedirectURL))
		return
	}

	session.Set("profile", userProfile)
	if err := session.Save(); err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}

	ctx.Redirect(http.StatusPermanentRedirect, successRedirectURL)
}

func (u *UserHandler) logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorInternal, "Failed to update session"))
		return
	}

	ctx.JSON(http.StatusOK, "Successfully logged out.")
}

func (u *UserHandler) getUser(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)
	ctx.JSON(http.StatusOK, user)
}

func (u *UserHandler) forgotPassword(ctx *gin.Context) {
	var req core.ResetPasswordRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload"))
		return
	}

	err = u.UserService.ForgotPassword(ctx, req)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, "Successfully sent reset password email.")
}