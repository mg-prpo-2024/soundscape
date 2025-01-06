package internal

import (
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
	user := Artist{
		Name: userDto.Email,
	}
	result := r.db.Create(&user)
	return result.Error
}
