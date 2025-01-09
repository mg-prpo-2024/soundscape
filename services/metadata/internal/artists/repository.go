package artists

import (
	"soundscape/shared/metadatadb"

	"gorm.io/gorm"
)

type Repository interface {
	GetArtist(userId string) (*metadatadb.Artist, error)
	GetArtistAlbums(artistId string) ([]*metadatadb.Album, error)
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetArtist(userId string) (*metadatadb.Artist, error) {
	artist := &metadatadb.Artist{}
	err := r.db.Where("user_id = ?", userId).First(&artist).Error
	return artist, err
}

func (r *repository) GetArtistAlbums(artistId string) ([]*metadatadb.Album, error) {
	var albums []*metadatadb.Album
	err := r.db.Where("artist_id = ?", artistId).Find(&albums).Error
	return albums, err
}
