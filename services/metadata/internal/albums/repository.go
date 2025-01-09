package albums

import (
	"soundscape/shared/metadatadb"

	"gorm.io/gorm"
)

type Repository interface {
	GetAlbums() ([]*metadatadb.Album, error)
	GetAlbum(id string) (*metadatadb.Album, error)
	GetAlbumSongs(albumId string) ([]*metadatadb.Song, error)
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetAlbums() ([]*metadatadb.Album, error) {
	albums := []*metadatadb.Album{}
	err := r.db.Model(&metadatadb.Album{}).Preload("Artist").Limit(10).Find(&albums).Error
	return albums, err
}

func (r *repository) GetAlbum(id string) (*metadatadb.Album, error) {
	album := &metadatadb.Album{}
	err := r.db.Model(&metadatadb.Album{}).Preload("Artist").Where("id = ?", id).First(&album).Error
	return album, err
}

func (r *repository) GetAlbumSongs(albumId string) ([]*metadatadb.Song, error) {
	var songs []*metadatadb.Song
	err := r.db.Where("album_id = ?", albumId).Find(&songs).Error
	return songs, err
}
