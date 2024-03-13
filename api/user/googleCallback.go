package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chizidotdev/shop/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	upsertedUser, err := u.pgStore.UpsertUser(ctx, repository.UpsertUserParams{
		FirstName: user.GivenName,
		LastName:  user.FamilyName,
		Email:     user.Email,
		GoogleID:  user.Id,
		Image:     user.Picture,
		Role:      repository.UserRoleCustomer,
	})
	if err != nil {
		// TODO: add error code to redirect url
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", errRedirectURL, "failed to upsert user"))
		return
	}

	userProfile := &repository.GetUserRow{
		ID:        upsertedUser.ID,
		FirstName: upsertedUser.FirstName,
		LastName:  upsertedUser.LastName,
		Email:     upsertedUser.Email,
		Image:     upsertedUser.Image,
		Role:      upsertedUser.Role,
		GoogleID:  upsertedUser.GoogleID,
		CreatedAt: upsertedUser.CreatedAt,
		UpdatedAt: upsertedUser.UpdatedAt,
	}
	store, err := u.pgStore.GetStoreByUserId(ctx, upsertedUser.ID)
	if err == nil {
		userProfile.StoreID = uuid.NullUUID{UUID: store.ID, Valid: true}
	}

	session.Set(profileKey, userProfile)
	if err := session.Save(); err != nil {
		log.Println("Error saving session: ", err)
		ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", errRedirectURL, "failed to save session"))
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, successRedirectURL)
}
