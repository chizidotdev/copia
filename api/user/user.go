package user

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
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
	pgStore *repository.Repository
	Config  oauth2.Config
}

type UserResponse struct {
	*repository.User
	Store *repository.Store `json:"store,omitempty"`
}

func NewUserHandler(pgStore *repository.Repository) *UserHandler {
	gob.Register(repository.User{})
	gob.Register(repository.GetUserRow{})

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
	session := sessions.Default(ctx)
	user := middleware.GetAuthenticatedUser(ctx)

	userProfile, err := u.pgStore.GetUser(ctx, user.ID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusUnauthorized,
			MessageID: "",
			Message:   "User not found.",
			Reason:    err.Error(),
		})
		return
	}

	session.Set(profileKey, userProfile)
	if err := session.Save(); err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to save session.",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    user,
		Message: "User retrieved successfully.",
	})
}
