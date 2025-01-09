package internal

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetArtist(userId string) (*Artist, error)
	GetArtistAlbums(artistId string) ([]*Album, error)
	GetAlbum(id string) (*Album, error)
	GetAlbumSongs(albumId string) ([]*Song, error)
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetArtist(userId string) (*Artist, error) {
	artist := &Artist{}
	err := r.db.Where("user_id = ?", userId).First(&artist).Error
	return artist, err
}

func (r *repository) GetArtistAlbums(artistId string) ([]*Album, error) {
	var albums []*Album
	err := r.db.Where("artist_id = ?", artistId).Find(&albums).Error
	return albums, err
}

func (r *repository) GetAlbum(id string) (*Album, error) {
	album := &Album{}
	err := r.db.Where("id = ?", id).First(&album).Error
	return album, err
}

func (r *repository) GetAlbumSongs(albumId string) ([]*Song, error) {
	var songs []*Song
	err := r.db.Where("album_id = ?", albumId).Find(&songs).Error
	return songs, err
}
