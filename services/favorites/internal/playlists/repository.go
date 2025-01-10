package playlists

import (
	"soundscape/shared/metadatadb"

	"gorm.io/gorm"
)

type Repository interface {
	GetPlaylists() ([]*metadatadb.Album, error)
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetPlaylists() ([]*metadatadb.Album, error) {
	albums := []*metadatadb.Album{}
	err := r.db.Model(&metadatadb.Album{}).Preload("Artist").Limit(10).Find(&albums).Error
	return albums, err
}
