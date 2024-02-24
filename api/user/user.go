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

var (
	stateKey   = middleware.StateKey
	profileKey = middleware.ProfileKey
)

var userRoles = map[string]repository.UserRole{
	"master":   repository.UserRoleMaster,
	"vendor":   repository.UserRoleVendor,
	"customer": repository.UserRoleCustomer,
}

type UserHandler struct {
	pgStore *repository.Queries
	Config  oauth2.Config
}

func NewUserHandler(pgStore *repository.Queries) *UserHandler {
	gob.Register(repository.User{})

	oauthConfig := oauth2.Config{
		ClientID:     config.EnvVars.GoogleClientID,
		ClientSecret: config.EnvVars.GoogleClientSecret,
		RedirectURL:  config.EnvVars.AuthCallbackURL,
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
	httpUtil.Success(ctx, http.StatusOK, httpUtil.SuccessResponse{
		Data:    user,
		Message: "User retrieved successfully.",
	})
}
