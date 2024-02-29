package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/chizidotdev/shop/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *UserHandler) GoogleCallback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	redirectURI := session.Get(redirectURIKey).(string)
	errRedirectURL := redirectURI + "?error="
	successRedirectURL := redirectURI + "?success=true"

	if ctx.Query(stateKey) != session.Get(stateKey) {
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", errRedirectURL, invalidStateError))
		return
	}

	code := ctx.Query("code")
	user, err := u.getGoogleUserData(ctx, code)
	if err != nil {
		log.Println("Error getting user data from google: ", err)
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", errRedirectURL, failedToExchangeError))
		return
	}

	userProfile, err := u.pgStore.UpsertUser(ctx, repository.UpsertUserParams{
		FirstName: user.GivenName,
		LastName:  user.FamilyName,
		Email:     user.Email,
		GoogleID:  sql.NullString{String: user.Id, Valid: true},
		Image:     user.Picture,
		Role:      repository.UserRoleCustomer,
	})
	if err != nil {
		log.Println("Error upserting user: ", err)
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", errRedirectURL, failedToExchangeError))
		return
	}

	session.Set(profileKey, userProfile)
	if err := session.Save(); err != nil {
		log.Println("Error saving session: ", err)
		ctx.Redirect(http.StatusTemporaryRedirect, errRedirectURL)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, successRedirectURL)
}
