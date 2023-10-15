package service

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"github.com/chizidotdev/copia/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// AuthService is used to authenticate our users.
type AuthService struct {
	Config oauth2.Config
}

// NewAuthenticator instantiates the *authService.
func NewAuthenticator() *AuthService {
	config := oauth2.Config{
		ClientID:     util.EnvVars.GoogleClientID,
		ClientSecret: util.EnvVars.GoogleClientSecret,
		RedirectURL:  util.EnvVars.AuthCallbackURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}

	return &AuthService{Config: config}
}

type UserData struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func (a *AuthService) GetUserData(code string) (UserData, error) {
	token, err := a.Config.Exchange(context.Background(), code)
	if err != nil {
		return UserData{}, err
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return UserData{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return UserData{}, err
	}

	gob.Register(UserData{})
	var user UserData
	err = json.Unmarshal(data, &user)
	if err != nil {
		return UserData{}, err
	}

	return user, nil
}

func (a *AuthService) GoogleCallback(code string) (UserData, error) {
	user, err := a.GetUserData(code)
	if err != nil {
		return UserData{}, util.Errorf(util.ErrorForbidden, "Failed to get user data")
	}

	return user, nil
}
