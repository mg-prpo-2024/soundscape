package internal

import (
	"database/sql"

	"gorm.io/gorm"
)

type Repository interface {
	CreateArtist(user ArtistDto) error
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateArtist(userDto ArtistDto) error {
	user := User{
		Auth0Id:          userDto.Id,
		Email:            userDto.Email,
		StripeCustomerId: sql.NullString{String: "", Valid: false},
	}
	result := r.db.Create(&user)
	return result.Error
}
