package service

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/token_manager"
	"github.com/chizidotdev/copia/util"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserParams) (dto.UserResponse, error)
	GetUser(ctx context.Context, req dto.LoginUserParams) (dto.UserResponse, error)
	GoogleCallback(ctx context.Context, code string) (dto.UserResponse, error)
	GetGoogleAuthConfig() oauth2.Config
}

type userService struct {
	Store        *repository.Repository
	TokenManager token_manager.TokenManager
	Config       oauth2.Config
}

func NewUserService(store *repository.Repository) UserService {
	gob.Register(dto.UserResponse{})
	tokenManager, err := token_manager.NewJWTTokenManager(util.EnvVars.AuthSecret)
	if err != nil {
		log.Fatal(err)
	}

	config := oauth2.Config{
		ClientID:     util.EnvVars.GoogleClientID,
		ClientSecret: util.EnvVars.GoogleClientSecret,
		RedirectURL:  util.EnvVars.AuthCallbackURL,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	return &userService{
		Store:        store,
		TokenManager: tokenManager,
		Config:       config,
	}
}

func (u *userService) GetGoogleAuthConfig() oauth2.Config {
	return u.Config
}

func (u *userService) CreateUser(ctx context.Context, req dto.CreateUserParams) (dto.UserResponse, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return dto.UserResponse{}, util.Errorf(util.ErrorInternal, "Failed to hash password")
	}

	user, err := u.Store.CreateUser(ctx, repository.CreateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return dto.UserResponse{}, util.Errorf(util.ErrorInternal, "Failed to create user: "+err.Error())
	}

	return dto.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *userService) GetUser(ctx context.Context, req dto.LoginUserParams) (dto.UserResponse, error) {
	user, err := u.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.UserResponse{}, util.Errorf(util.ErrorUnauthorized, "Invalid credentials")
	}

	err = util.ComparePassword(user.Password, req.Password)
	if err != nil {
		return dto.UserResponse{}, util.Errorf(util.ErrorUnauthorized, "Invalid credentials")
	}

	return dto.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

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

func (u *userService) GoogleCallback(ctx context.Context, code string) (dto.UserResponse, error) {
	user, err := u.getGoogleUserData(code)
	if err != nil {
		return dto.UserResponse{}, util.Errorf(util.ErrorForbidden, "Failed to get user data")
	}

	userProfile, err := u.Store.UpsertUser(ctx, repository.CreateUserParams{
		FirstName: user.GivenName,
		LastName:  user.FamilyName,
		Email:     user.Email,
		Password:  "",
		GoogleID:  user.Id,
	})
	if err != nil {
		return dto.UserResponse{}, util.Errorf(util.ErrorInternal, "Failed to create new user")
	}

	return dto.UserResponse{
		ID:        userProfile.ID,
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
		Email:     userProfile.Email,
		CreatedAt: userProfile.CreatedAt,
		UpdatedAt: userProfile.UpdatedAt,
	}, nil
}

func (u *userService) getGoogleUserData(code string) (UserData, error) {
	config := u.GetGoogleAuthConfig()
	token, err := config.Exchange(context.Background(), code)
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

	gob.Register(UserData{})
	var user UserData
	err = json.Unmarshal(data, &user)
	if err != nil {
		return UserData{}, err
	}

	return user, nil
}
