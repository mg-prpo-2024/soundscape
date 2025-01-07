package internal

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateArtist(artist CreateArtistDto) error
	GetArtist(userId string) (*Artist, error)
	CreateAlbum(album CreateAlbumDto) (*Album, error)
	GetAlbum(id string) (*Album, error)
	CreateSong(song CreateSongDto) (*Song, error)
	GetAlbumSongs(albumId string) ([]*Song, error)
	DeleteSong(id string) error
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateArtist(userDto CreateArtistDto) error {
	user := Artist{
		UserId: userDto.UserId,
		Name:   userDto.Name,
		Bio:    userDto.Bio,
	}
	result := r.db.Create(&user)
	return result.Error
}

func (r *repository) GetArtist(userId string) (*Artist, error) {
	artist := &Artist{}
	err := r.db.Where("user_id = ?", userId).First(&artist).Error
	return artist, err
}

func (r *repository) CreateAlbum(albumDto CreateAlbumDto) (*Album, error) {
	album := Album{
		Title: albumDto.Title,
	}
	result := r.db.Create(&album)
	return &album, result.Error
}

func (r *repository) GetAlbum(id string) (*Album, error) {
	album := &Album{}
	err := r.db.Where("id = ?", id).First(&album).Error
	return album, err
}

func (r *repository) CreateSong(songDto CreateSongDto) (*Song, error) {
	song := Song{
		Title:      songDto.Title,
		TrackOrder: songDto.TrackOrder,
		AlbumId:    uuid.MustParse(songDto.AlbumId),
		BlobUrl:    "",
		Artists: []*Artist{
			{Base: Base{ID: uuid.MustParse(songDto.ArtistId)}},
		},
	}
	result := r.db.Omit("Artists.*").Create(&song)
	return &song, result.Error
}

func (r *repository) GetAlbumSongs(albumId string) ([]*Song, error) {
	var songs []*Song
	err := r.db.Where("album_id = ?", albumId).Find(&songs).Error
	return songs, err
}

func (r *repository) DeleteSong(id string) error {
	result := r.db.Delete(&Song{}, "id = ?", id)

	return result.Error
}
