package internal

import (
	"soundscape/shared/metadatadb"

	"gorm.io/gorm"
)

type Repository interface {
	GetArtist(userId string) (*metadatadb.Artist, error)
	GetArtistAlbums(artistId string) ([]*metadatadb.Album, error)
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

func (r *repository) GetAlbum(id string) (*metadatadb.Album, error) {
	album := &metadatadb.Album{}
	err := r.db.Where("id = ?", id).First(&album).Error
	return album, err
}

func (r *repository) GetAlbumSongs(albumId string) ([]*metadatadb.Song, error) {
	var songs []*metadatadb.Song
	err := r.db.Where("album_id = ?", albumId).Find(&songs).Error
	return songs, err
}
