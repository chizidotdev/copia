package user

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/chizidotdev/shop/config"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	InvalidStateError     = "invalid_state"
	FailedToExchangeError = "failed_to_exchange"
)

func (u *UserHandler) GoogleCallback(ctx *gin.Context) {
	errRedirectURL := config.EnvVars.AuthCallbackURL + "?errors="
	successRedirectURL := config.EnvVars.AuthCallbackURL

	session := sessions.Default(ctx)
	if ctx.Query(stateKey) != session.Get(stateKey) {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s%s", errRedirectURL, InvalidStateError))
		return
	}

	code := ctx.Query("code")
	user, err := u.getGoogleUserData(ctx, code)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s%s", errRedirectURL, FailedToExchangeError))
		return
	}

	userProfile, err := u.pgStore.UpsertUser(ctx, repository.UpsertUserParams{
		FirstName: user.GivenName,
		LastName:  user.FamilyName,
		Email:     user.Email,
		GoogleID:  sql.NullString{String: user.Id, Valid: true},
		Image:     user.Picture,
	})
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s%s", errRedirectURL, FailedToExchangeError))
		return
	}

	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s%s", errRedirectURL, FailedToExchangeError))
		return
	}

	session.Set(profileKey, userProfile)
	if err := session.Save(); err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}

	ctx.Redirect(http.StatusPermanentRedirect, successRedirectURL)
}
