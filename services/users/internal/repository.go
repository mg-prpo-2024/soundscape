package internal

import (
	"database/sql"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user UserDto) error
	GetUser(id string) (*UserDto, error)
	SetCustomerId(auth0Id string, customerId string) error
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(userDto UserDto) error {
	user := User{
		Auth0Id:          userDto.Id,
		Email:            userDto.Email,
		StripeCustomerId: sql.NullString{String: "", Valid: false},
	}
	result := r.db.Create(&user)
	return result.Error
}

func (r *repository) GetUser(id string) (*UserDto, error) {
	user := User{}
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &UserDto{
		Id:               user.Auth0Id,
		Email:            user.Email,
		StripeCustomerId: &user.StripeCustomerId.String,
	}, nil
}

func (r *repository) SetCustomerId(auth0Id string, customerId string) error {
	result := r.db.Model(&User{}).Where("auth0_id = ?", auth0Id).Update("stripe_customer_id", customerId)
	return result.Error
}
