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
	"time"
)

type UserService struct {
	Store      core.UserRepository
	emailStore core.EmailRepository
	redisStore core.RedisRepository
	Config     oauth2.Config
}

func NewUserService(
	userRepo core.UserRepository,
	emailRepo core.EmailRepository,
	redisRepo core.RedisRepository,
) *UserService {
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
		Store:      userRepo,
		emailStore: emailRepo,
		redisStore: redisRepo,
		Config:     oauthConfig,
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

func (u *UserService) sendEmailVerificationEmail(ctx context.Context, user core.User) error {
	tokenDuration := time.Hour * 24
	token, err := u.GenerateAuthState()
	if err != nil {
		return err
	}
	err = u.redisStore.Set(ctx, token, user.Email, tokenDuration)

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

	log.Printf("Email verification link: %s", url)
	err = u.emailStore.SendEmail([]string{user.Email}, "Verify Email", emailBody)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetGoogleAuthConfig() oauth2.Config {
	return u.Config
}

func (u *UserService) CreateUser(ctx context.Context, req core.CreateUserRequest) (core.UserResponse, error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to hash password",
			Reason:    err.Error(),
		}
		return core.UserResponse{}, errors.Errorf(errResp)
	}

	user, err := u.Store.CreateUser(ctx, core.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to create user",
			Reason:    err.Error(),
		}
		return core.UserResponse{}, errors.Errorf(errResp)
	}

	go func() {
		emailErr := u.sendEmailVerificationEmail(ctx, user)
		if emailErr != nil {
			log.Println("Failed to send email: " + emailErr.Error())
		}
	}()

	// TODO: Use a message queue to create product settings

	return core.UserResponse{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	}, nil
}

func (u *UserService) GetUser(ctx context.Context, req core.LoginUserRequest) (core.UserResponse, error) {
	user, err := u.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorUnauthorized,
			MessageID: "",
			Message:   "Invalid credentials.",
			Reason:    err.Error(),
		}
		return core.UserResponse{}, errors.Errorf(errResp)
	}

	if user.GoogleID != "" {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorUnauthorized,
			MessageID: "",
			Message:   "Looks like you signed up with Google. Please login with Google",
			Reason:    err.Error(),
		}
		return core.UserResponse{}, errors.Errorf(errResp)
	}

	err = comparePassword(user.Password, req.Password)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorUnauthorized,
			MessageID: "",
			Message:   "Invalid credentials.",
			Reason:    err.Error(),
		}
		return core.UserResponse{}, errors.Errorf(errResp)
	}

	return core.UserResponse{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	}, nil
}

func (u *UserService) SendVerificationEmail(ctx context.Context, email string) error {
	user, err := u.Store.GetUserByEmail(ctx, email)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "User not found",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	if user.EmailVerified == true {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Email already verified",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	go func() {
		emailErr := u.sendEmailVerificationEmail(ctx, user)
		if emailErr != nil {
			log.Println("Failed to send email: " + emailErr.Error())
		}
	}()

	return nil
}

func (u *UserService) VerifyEmail(ctx context.Context, req core.VerifyEmailRequest) error {
	email, err := u.redisStore.Get(ctx, req.Code)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Code is invalid",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	user, err := u.Store.GetUserByEmail(ctx, email)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "User not found",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	_, err = u.Store.UpdateUser(ctx, core.User{
		Email:         user.Email,
		EmailVerified: true,
	})
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to update user",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	go func() {
		_ = u.redisStore.Delete(ctx, req.Code)

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
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Email does not exist.",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	if user.EmailVerified == false {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Email not verified.",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}
	if user.GoogleID != "" {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Looks like you signed up with Google. Please login with Google",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	tokenDuration := time.Minute * 15
	token, err := u.GenerateAuthState()
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to create token",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	err = u.redisStore.Set(ctx, token, user.Email, tokenDuration)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to set token",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
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

	log.Printf("Email verification link: %s", url)
	err = u.emailStore.SendEmail(
		[]string{user.Email},
		"Reset Password",
		emailBody,
	)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to send email",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	return nil
}

func (u *UserService) ChangePassword(ctx context.Context, req core.ChangePasswordRequest) error {
	email, err := u.redisStore.Get(ctx, req.Code)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Code is invalid",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to hash password",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	_, err = u.Store.UpdateUser(ctx, core.User{
		Email:    email,
		Password: hashedPassword,
	})
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to update password",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
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
		errResp := errors.ErrResponse{
			Code:      errors.ErrorForbidden,
			MessageID: "",
			Message:   "Failed to get Google user data",
			Reason:    err.Error(),
		}
		return core.UserResponse{}, errors.Errorf(errResp)
	}

	userExists, err := u.Store.GetUserByEmail(ctx, user.Email)
	if err == nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorForbidden,
			MessageID: "",
			Message:   "Looks like you signed up with email and password. Please login with email and password",
			Reason:    err.Error(),
		}

		if userExists.GoogleID == "" {
			return core.UserResponse{}, errors.Errorf(errResp)
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
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to update user",
			Reason:    err.Error(),
		}
		return core.UserResponse{}, errors.Errorf(errResp)
	}

	return core.UserResponse{
		ID:            userProfile.ID,
		FirstName:     userProfile.FirstName,
		LastName:      userProfile.LastName,
		Email:         userProfile.Email,
		EmailVerified: userProfile.EmailVerified,
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
