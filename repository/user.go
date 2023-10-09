package repository

import "context"

type CreateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"user_email"`
	Password  string `json:"password"`
}

func (r *Repository) CreateUser(_ context.Context, arg CreateUserParams) (User, error) {
	user := User{
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		Password:  arg.Password,
	}
	result := r.DB.Create(&user)
	return user, result.Error
}

func (r *Repository) GetUserByEmail(_ context.Context, email string) (User, error) {
	var user User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return user, err
}
