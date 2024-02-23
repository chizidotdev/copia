package user

import (
	"net/http"

	"github.com/chizidotdev/shop/config"
	"github.com/chizidotdev/shop/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (u *UserHandler) GoogleLogin(ctx *gin.Context) {
	errRedirectURL := config.EnvVars.AuthCallbackURL + "?errors="
	state, err := util.GenerateRandString(32)
	if err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}
	session := sessions.Default(ctx)
	session.Set(stateKey, state)
	if err := session.Save(); err != nil {
		ctx.Redirect(http.StatusPermanentRedirect, errRedirectURL)
		return
	}

	url := u.Config.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}
