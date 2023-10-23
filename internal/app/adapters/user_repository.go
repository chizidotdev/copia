package adapters

import (
	"context"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/internal/app/usecases"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ usecases.UserRepository = (*UserRepositoryImpl)(nil)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{DB: db}
}

type User struct {
	Base
	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"not null;uniqueIndex" json:"email"`
	Password  string `json:"password"`
	GoogleID  string `gorm:"unique" json:"google_id"`

	Orders    []Order `gorm:"foreignKey:UserID" json:"orders"`
}

func (r *UserRepositoryImpl) CreateUser(_ context.Context, arg core.User) (core.User, error) {
	user := User{
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		Password:  arg.Password,
	}
	err := r.DB.Create(&user).Error
	return core.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}, err
}

func (r *UserRepositoryImpl) UpsertUser(_ context.Context, arg core.User) (core.User, error) {
	user := User{
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		Password:  arg.Password,
		GoogleID:  arg.GoogleID,
	}
	result := r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		UpdateAll: true,
	}).Create(&user)
	return core.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}, result.Error
}

func (r *UserRepositoryImpl) GetUserByEmail(_ context.Context, email string) (core.User, error) {
	var user User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return core.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}, err
}
