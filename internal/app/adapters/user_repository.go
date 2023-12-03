package adapters

import (
	"context"
	"errors"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/internal/app/usecases"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

var _ usecases.UserRepository = (*UserRepositoryImpl)(nil)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Panic("Failed to migrate User", err)
	}
	return &UserRepositoryImpl{DB: db}
}

type User struct {
	Base
	FirstName     string `gorm:"not null" json:"firstName"`
	LastName      string `gorm:"not null" json:"lastName"`
	Email         string `gorm:"not null;uniqueIndex" json:"email"`
	EmailVerified bool   `gorm:"not null;default:false" json:"emailVerified"`
	Password      string `json:"password"`
	GoogleID      string `json:"googleID"`

	Orders          []Order         `json:"orders"`
	Products        []Product       `json:"products"`
	ProductSettings ProductSettings `json:"productSettings"`
}

func (r *UserRepositoryImpl) CreateUser(_ context.Context, arg core.User) (core.User, error) {
	user := User{
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		Password:  arg.Password,
	}
	err := r.DB.Create(&user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return core.User{}, errors.New("email already exists")
	}
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
		FirstName:     arg.FirstName,
		LastName:      arg.LastName,
		Email:         arg.Email,
		EmailVerified: arg.EmailVerified,
		Password:      arg.Password,
		GoogleID:      arg.GoogleID,
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

func (r *UserRepositoryImpl) UpdateUser(_ context.Context, arg core.User) (core.User, error) {
	user := User{
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Email:     arg.Email,
		Password:  arg.Password,
	}
	err := r.DB.Model(&User{}).Where("email = ?", arg.Email).Updates(&user).Error
	return core.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}, err
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
		GoogleID:  user.GoogleID,
	}, err
}
