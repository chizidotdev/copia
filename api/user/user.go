package user

import (
	"encoding/gob"
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/config"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	profileKey = middleware.ProfileKey

	stateKey       = "state"
	redirectURIKey = "redirect_uri"

	severError            = "internal_server_error"
	invalidStateError     = "invalid_state"
	failedToExchangeError = "failed_to_exchange"
)

type UserHandler struct {
	pgStore *repository.Queries
	Config  oauth2.Config
}

type UserResponse struct {
	*repository.User
	Store *repository.Store `json:"store,omitempty"`
}

func NewUserHandler(pgStore *repository.Queries) *UserHandler {
	gob.Register(repository.User{})

	oauthConfig := oauth2.Config{
		ClientID:     config.EnvVars.GoogleClientID,
		ClientSecret: config.EnvVars.GoogleClientSecret,
		RedirectURL:  config.EnvVars.GoogleRedirectURI,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	return &UserHandler{
		pgStore: pgStore,
		Config:  oauthConfig,
	}
}

func (u *UserHandler) GetUser(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)
	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    user,
		Message: "User retrieved successfully.",
	})
}
