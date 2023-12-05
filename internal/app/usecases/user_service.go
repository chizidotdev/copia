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
	"github.com/chizidotdev/copia/token_manager"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg core.User) (core.User, error)
	UpsertUser(ctx context.Context, arg core.User) (core.User, error)
	UpdateUser(ctx context.Context, arg core.User) (core.User, error)
	GetUserByEmail(ctx context.Context, email string) (core.User, error)
}

type UserService struct {
	Store        UserRepository
	emailStore   core.EmailRepository
	tokenManager token_manager.TokenManager
	Config       oauth2.Config
}

func NewUserService(userRepo UserRepository, emailRepo core.EmailRepository) *UserService {
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

	tokenManager, err := token_manager.NewJWTTokenManager(config.EnvVars.AuthSecret)
	if err != nil {
		log.Fatal("Error initializing JWT token manager")
	}

	return &UserService{
		Store:        userRepo,
		emailStore:   emailRepo,
		Config:       oauthConfig,
		tokenManager: tokenManager,
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

func (u *UserService) sendEmailVerificationEmail(user core.User) {
	tokenDuration := time.Hour * 24
	token, err := u.tokenManager.CreateToken(user.Email, tokenDuration)
	if err != nil {
		log.Println("Failed to create token: " + err.Error())
		return
	}

	// TODO: Use a message queue to send email
	url := fmt.Sprintf("%s/%s?code=%s", config.EnvVars.AuthDomain, "u/verify-email", token)
	emailBody := fmt.Sprintf(`
			<p>Hi %s,</p>
			<p>Click the link below to verify your email.</p>
			<p><a href="%s">Verify Email</a></p>
			<p>This link will expire in 24 hours.</p>
			<br />
			<p>If you did not create an account, please ignore this email.</p>
			<p>Thanks,</p>
			<p>Copia Team.</p>
		`, user.FirstName, url)

	err = u.emailStore.SendEmail(
		[]string{user.Email},
		"Verify Email",
		emailBody,
	)
	if err != nil {
		log.Println("Failed to send email: " + err.Error())
	}
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

	go u.sendEmailVerificationEmail(user)

	// TODO: Use a message queue to create product settings

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
		return core.UserResponse{}, errors.Errorf(errors.ErrorUnauthorized, "Invalid credentials.")
	}

	if user.GoogleID != "" {
		return core.UserResponse{}, errors.Errorf(
			errors.ErrorUnauthorized,
			"Looks like you signed up with Google. Please login with Google",
		)
	}

	err = comparePassword(user.Password, req.Password)
	if err != nil {
		return core.UserResponse{}, errors.Errorf(errors.ErrorUnauthorized, "Invalid credentials.")
	}

	return core.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}, nil
}

func (u *UserService) SendVerificationEmail(ctx context.Context, email string) error {
	user, err := u.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.Errorf(errors.ErrorBadRequest, "User not found")
	}

	if user.EmailVerified == true {
		return errors.Errorf(errors.ErrorBadRequest, "Email already verified.")
	}

	go u.sendEmailVerificationEmail(user)

	return nil
}

func (u *UserService) VerifyEmail(ctx context.Context, req core.VerifyEmailRequest) error {
	payload, err := u.tokenManager.VerifyToken(req.Token)
	if err != nil {
		return errors.Errorf(errors.ErrorUnauthorized, "Token is invalid")
	}
	if err := payload.Valid(); err != nil {
		return errors.Errorf(errors.ErrorUnauthorized, "Token expired")
	}

	user, err := u.Store.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return errors.Errorf(errors.ErrorBadRequest, "User not found")
	}

	_, err = u.Store.UpdateUser(ctx, core.User{
		Email:         user.Email,
		EmailVerified: true,
	})
	if err != nil {
		return errors.Errorf(errors.ErrorInternal, "Failed to update user")
	}

	go func() {
		emailBody := fmt.Sprintf(`
			<p>Hi %s,</p>
			<p>Your Copia email has been verified successfully.</p>
			<br />
			<p>Thanks,</p>
			<p>Copia Team.</p>
		`, user.FirstName)

		err = u.emailStore.SendEmail(
			[]string{user.Email},
			"Email Verified",
			emailBody,
		)
		if err != nil {
			log.Println("Failed to send email: " + err.Error())
		}
	}()

	return nil

}

func (u *UserService) ResetPassword(ctx context.Context, email string) error {
	user, err := u.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.Errorf(errors.ErrorBadRequest, "User not found")
	}

	if user.EmailVerified == false {
		return errors.Errorf(errors.ErrorBadRequest, "Email not verified.")
	}
	if user.GoogleID != "" {
		return errors.Errorf(errors.ErrorBadRequest, "Looks like you signed up with Google. Please login with Google")
	}

	tokenDuration := time.Minute * 15
	token, err := u.tokenManager.CreateToken(email, tokenDuration)
	if err != nil {
		return errors.Errorf(errors.ErrorInternal, "Failed to create token")
	}

	url := fmt.Sprintf("%s/%s?code=%s", config.EnvVars.AuthDomain, "u/change-password", token)
	emailBody := fmt.Sprintf(`
		<p>Hi %s,</p>
		<p>Click the link below to reset your password.</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in 15minutes.</p>
		<br />
		<p>If you did not request a password reset, please ignore this email.</p>
		<p>Thanks,</p>
		<p>Copia Team.</p>
	`, user.FirstName, url)

	err = u.emailStore.SendEmail(
		[]string{user.Email},
		"Reset Password",
		emailBody,
	)
	if err != nil {
		return errors.Errorf(errors.ErrorInternal, "Failed to send email: "+err.Error())
	}

	return nil
}

func (u *UserService) ChangePassword(ctx context.Context, req core.ChangePasswordRequest) error {
	payload, err := u.tokenManager.VerifyToken(req.Token)
	if err != nil {
		return errors.Errorf(errors.ErrorUnauthorized, "Token is invalid")
	}
	if err := payload.Valid(); err != nil {
		return errors.Errorf(errors.ErrorUnauthorized, "Token expired")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return errors.Errorf(errors.ErrorInternal, "Failed to hash password")
	}

	_, err = u.Store.UpdateUser(ctx, core.User{
		Email:    payload.Email,
		Password: hashedPassword,
	})
	if err != nil {
		return errors.Errorf(errors.ErrorInternal, "Failed to update password")
	}

	return nil
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

	userExists, err := u.Store.GetUserByEmail(ctx, user.Email)
	if err == nil {
		if userExists.GoogleID == "" {
			return core.UserResponse{}, errors.Errorf(
				errors.ErrorForbidden,
				"Looks like you signed up with email and password. Please login with email and password",
			)
		}
	}

	userProfile, err := u.Store.UpsertUser(ctx, core.User{
		FirstName:     user.GivenName,
		LastName:      user.FamilyName,
		Email:         user.Email,
		EmailVerified: user.VerifiedEmail,
		Password:      "",
		GoogleID:      user.Id,
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
