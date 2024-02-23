package user

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type UserData struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (u *UserHandler) getGoogleUserData(ctx context.Context, code string) (UserData, error) {
	token, err := u.Config.Exchange(ctx, code)
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

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return UserData{}, err
	}

	// gob.Register(UserData{})
	var user UserData
	err = json.Unmarshal(data, &user)
	if err != nil {
		return UserData{}, err
	}

	return user, nil
}
