package core

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg User) (User, error)
	UpsertUser(ctx context.Context, arg User) (User, error)
	UpdateUser(ctx context.Context, arg User) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
}

type User struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
	Password      string    `json:"password"`
	GoogleID      string    `json:"googleID"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type VerifyEmailRequest struct {
	Code string `json:"code" binding:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

type ChangePasswordRequest struct {
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"emailVerified"`
}
