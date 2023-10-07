package service

import (
	"context"
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/util"
)

type UserService interface {
	CreateUser(ctx context.Context, req dto.CreateUserParams) (repository.User, error)
}

type userService struct {
	Store *repository.Repository
}

func NewUserService(store *repository.Repository) UserService {
	return &userService{
		Store: store,
	}
}

func (u *userService) CreateUser(ctx context.Context, req dto.CreateUserParams) (repository.User, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return repository.User{}, err
	}

	user, err := u.Store.CreateUser(ctx, repository.CreateUserParams{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return repository.User{}, err
	}
	return user, nil
}
