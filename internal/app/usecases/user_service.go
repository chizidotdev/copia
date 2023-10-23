package usecases

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/chizidotdev/copia/config"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg core.User) (core.User, error)
	UpsertUser(ctx context.Context, arg core.User) (core.User, error)
	GetUserByEmail(ctx context.Context, email string) (core.User, error)
}

type UserService struct {
	Store  UserRepository
	Config oauth2.Config
}

func NewUserService(userRepo UserRepository) *UserService {
	gob.Register(core.UserResponse{})

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

	return &UserService{
		Store:  userRepo,
		Config: oauthConfig,
	}
}

// hashPassword returns the bcrypt hash of the password.
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// comparePassword compares the hashed password with the password.
func comparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match: %w", err)
	}
	return nil
}

func (u *UserService) GetGoogleAuthConfig() oauth2.Config {
	return u.Config
}

func (u *UserService) CreateUser(ctx context.Context, req core.CreateUserRequest) (core.UserResponse, error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return core.UserResponse{}, errors.Errorf(errors.ErrorInternal, "Failed to hash password")
	}

	user, err := u.Store.CreateUser(ctx, core.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return core.UserResponse{}, errors.Errorf(errors.ErrorBadRequest, "Failed to create User: "+err.Error())
	}

	return core.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

func (u *UserService) GetUser(ctx context.Context, req core.LoginUserRequest) (core.UserResponse, error) {
	user, err := u.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return core.UserResponse{}, errors.Errorf(errors.ErrorUnauthorized, "Invalid credentials. User Email not found")
	}

	log.Println(user)
	err = comparePassword(user.Password, req.Password)
	if err != nil {
		return core.UserResponse{}, errors.Errorf(errors.ErrorUnauthorized, err.Error())
	}

	return core.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
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

func (u *UserService) GoogleCallback(ctx context.Context, code string) (core.UserResponse, error) {
	user, err := u.getGoogleUserData(code)
	if err != nil {
		return core.UserResponse{}, errors.Errorf(errors.ErrorForbidden, "Failed to get Store data")
	}

	userProfile, err := u.Store.UpsertUser(ctx, core.User{
		FirstName: user.GivenName,
		LastName:  user.FamilyName,
		Email:     user.Email,
		Password:  "",
		GoogleID:  user.Id,
	})
	if err != nil {
		return core.UserResponse{}, errors.Errorf(errors.ErrorInternal, "Failed to create new Store")
	}

	return core.UserResponse{
		ID:        userProfile.ID,
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
		Email:     userProfile.Email,
	}, nil
}

func (u *UserService) getGoogleUserData(code string) (UserData, error) {
	authConfig := u.GetGoogleAuthConfig()
	token, err := authConfig.Exchange(context.Background(), code)
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

func (u *UserService) GenerateAuthState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
