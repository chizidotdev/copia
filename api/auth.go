package api

import (
	"fmt"
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) login(ctx *gin.Context) {
	var req dto.LoginUserParams
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorMessage(err.Error()))
		return
	}

	user, err := s.UserService.GetUser(ctx, req)
	if err != nil {
		ctx.JSON(errorResponse(err))
		return
	}

	session := sessions.Default(ctx)
	session.Set("profile", user)
	if err := session.Save(); err != nil {
		ctx.JSON(errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Successfully logged in.")
}

func (s *Server) loginWithSSO(ctx *gin.Context) {
	errRedirectURL := util.EnvVars.AuthDomain + "/u/login/error"
	state, err := util.GenerateRandomState()
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

	googleAuthConfig := s.UserService.GetGoogleAuthConfig()
	url := googleAuthConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Server) ssoCallback(ctx *gin.Context) {
	errRedirectURL := util.EnvVars.AuthDomain + "/u/login/error"
	successRedirectURL := util.EnvVars.AuthDomain + "/u/login/success"

	session := sessions.Default(ctx)
	if ctx.Query("state") != session.Get("state") {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s?error=invalid_state", errRedirectURL))
		return
	}

	code := ctx.Query("code")
	userProfile, err := s.UserService.GoogleCallback(ctx, code)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s?error=failed_to_exchange", errRedirectURL))
		return
	}

	session.Set("profile", userProfile)
	if err := session.Save(); err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}

	ctx.Redirect(http.StatusPermanentRedirect, successRedirectURL)
}

func (s *Server) logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		ctx.JSON(errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Successfully logged out.")
}

