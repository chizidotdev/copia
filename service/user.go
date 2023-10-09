package service

import (
	"context"
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/token_manager"
	"github.com/chizidotdev/copia/util"
	"log"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserParams) (dto.UserResponse, error)
	GetUser(ctx context.Context, req dto.LoginUserParams) (dto.LoginUserResponse, error)
}

type userService struct {
	Store        *repository.Repository
	TokenManager token_manager.TokenManager
}

func NewUserService(store *repository.Repository) UserService {
	tokenManager, err := token_manager.NewJWTTokenManager(util.EnvVars.AuthSecret)
	if err != nil {
		log.Fatal(err)
	}

	return &userService{
		Store:        store,
		TokenManager: tokenManager,
	}
}

func (u *userService) CreateUser(ctx context.Context, req dto.CreateUserParams) (dto.UserResponse, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user, err := u.Store.CreateUser(ctx, repository.CreateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return dto.UserResponse{}, err
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

func (u *userService) GetUser(ctx context.Context, req dto.LoginUserParams) (dto.LoginUserResponse, error) {
	user, err := u.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	err = util.ComparePassword(user.Password, req.Password)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	accessToken, err := u.TokenManager.CreateToken(req.Email, time.Minute*15)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	return dto.LoginUserResponse{
		AccessToken: accessToken,
		User: dto.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
