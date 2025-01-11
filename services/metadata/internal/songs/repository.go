package songs

import (
	"soundscape/shared/metadatadb"

	"gorm.io/gorm"
)

type Repository interface {
	GetSongs(songIds []string) ([]*metadatadb.Song, error)
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetSongs(songIds []string) ([]*metadatadb.Song, error) {
	var songs []*metadatadb.Song
	err := r.db.Where("id in (?)", songIds).Find(&songs).Error
	return songs, err
}
