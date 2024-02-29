package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chizidotdev/shop/config"
	"github.com/chizidotdev/shop/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *UserHandler) GoogleLogin(ctx *gin.Context) {
	redirectURL := config.EnvVars.AuthCallbackURI

	log.Println("redirectURL: ", redirectURL)
	if redirectURL == "" {
		ctx.String(http.StatusForbidden, "Invalid redirect_uri")
		return
	}

	errRedirectURL := redirectURL + "?error="
	state, err := util.GenerateRandString(32)
	if err != nil {
		log.Println("Error generating state: ", err)
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s%s", errRedirectURL, severError))
		return
	}
	session := sessions.Default(ctx)
	session.Set(stateKey, state)
	session.Set(redirectURIKey, redirectURL)
	if err := session.Save(); err != nil {
		log.Println("Error saving session: ", err)
		ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s%s", errRedirectURL, severError))
		return
	}

	url := u.Config.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
