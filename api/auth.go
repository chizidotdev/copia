package api

import (
	"crypto/rand"
	"encoding/base64"
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
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err.Error()))
		return
	}

	user, err := s.UserService.GetUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *Server) loginWithSSO(ctx *gin.Context) {
	state, err := generateRandomState()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err.Error()))
		return
	}
	session := sessions.Default(ctx)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err.Error()))
		return
	}

	url := s.AuthService.Config.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Server) callback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if ctx.Query("state") != session.Get("state") {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse("Invalid state parameter."))
		return
	}

	code := ctx.Query("code")
	userProfile, err := s.AuthService.GetUserData(code)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse("Failed to exchange an authorization code for a token."))
		return
	}

	session.Set("profile", userProfile)
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err.Error()))
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, util.EnvVars.AuthDomain)
}

func (s *Server) getUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	if profile == nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse("Unauthorized"))
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func (s *Server) isAuthenticated(ctx *gin.Context) {
	if sessions.Default(ctx).Get("profile") == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse("Unauthorized"))
	} else {
		ctx.Next()
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
