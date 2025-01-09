package internal

import (
	"soundscape/shared/metadatadb"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateArtist(artist CreateArtistDto) error
	GetArtist(userId string) (*metadatadb.Artist, error)
	GetArtistAlbums(artistId string) ([]*metadatadb.Album, error)
	CreateAlbum(album CreateAlbumDto) (*metadatadb.Album, error)
	GetAlbum(id string) (*metadatadb.Album, error)
	CreateSong(song CreateSongDto) (*metadatadb.Song, error)
	GetAlbumSongs(albumId string) ([]*metadatadb.Song, error)
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
	user := metadatadb.Artist{
		UserId: userDto.UserId,
		Name:   userDto.Name,
		Bio:    userDto.Bio,
	}
	result := r.db.Create(&user)
	return result.Error
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

func (r *repository) CreateAlbum(albumDto CreateAlbumDto) (*metadatadb.Album, error) {
	album := metadatadb.Album{
		Title:    albumDto.Title,
		ArtistId: uuid.MustParse(albumDto.ArtistId),
	}
	result := r.db.Create(&album)
	return &album, result.Error
}

func (r *repository) GetAlbum(id string) (*metadatadb.Album, error) {
	album := &metadatadb.Album{}
	err := r.db.Where("id = ?", id).First(&album).Error
	return album, err
}

func (r *repository) CreateSong(songDto CreateSongDto) (*metadatadb.Song, error) {
	song := metadatadb.Song{
		Title:      songDto.Title,
		TrackOrder: songDto.TrackOrder,
		AlbumId:    uuid.MustParse(songDto.AlbumId),
		BlobUrl:    "",
		Artists: []*metadatadb.Artist{
			{Base: metadatadb.Base{ID: uuid.MustParse(songDto.ArtistId)}},
		},
	}
	result := r.db.Omit("Artists.*").Create(&song)
	return &song, result.Error
}

func (r *repository) GetAlbumSongs(albumId string) ([]*metadatadb.Song, error) {
	var songs []*metadatadb.Song
	err := r.db.Where("album_id = ?", albumId).Find(&songs).Error
	return songs, err
}

func (r *repository) DeleteSong(id string) error {
	result := r.db.Delete(&metadatadb.Song{}, "id = ?", id)

	return result.Error
}
