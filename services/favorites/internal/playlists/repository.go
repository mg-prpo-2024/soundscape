package playlists

import (
	"errors"
	"soundscape/services/favorites/internal/apierrors"
	"soundscape/services/favorites/internal/dtos"
	"soundscape/services/favorites/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreatePlaylist(userId string, playlist dtos.CreatePlaylist) (*models.Playlist, error)
	GetPlaylists(userId string) ([]*models.Playlist, error)
	GetPlaylist(userId, playlistId string) (*models.Playlist, error)
	CreatePlaylistSong(userId, playlistId, songId string) error
	DeletePlaylistSong(userId, playlistId, songId string) error
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreatePlaylist(userId string, playlistDto dtos.CreatePlaylist) (*models.Playlist, error) {
	playlist := &models.Playlist{
		UserId: userId,
		Name:   playlistDto.Name,
	}
	err := r.db.Create(playlist).Error
	return playlist, err
}

func (r *repository) GetPlaylists(userId string) ([]*models.Playlist, error) {
	playlists := []*models.Playlist{}
	err := r.db.Model(&models.Playlist{}).Where("user_id = ?", userId).Preload("Songs").Find(&playlists).Error
	return playlists, err
}

func (r *repository) GetPlaylist(userId, playlistId string) (*models.Playlist, error) {
	playlist := &models.Playlist{Base: models.Base{ID: uuid.MustParse(playlistId)}}
	err := r.db.Model(&models.Playlist{}).Preload("Songs").First(playlist).Error
	return playlist, err
}

func (r *repository) CreatePlaylistSong(userId, playlistId, songId string) error {
	playlist := &models.Playlist{Base: models.Base{ID: uuid.MustParse(playlistId)}}
	err := r.db.First(playlist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierrors.ErrNotFound
		}
		return err
	}
	if playlist.UserId != userId {
		return apierrors.ErrForbidden
	}
	return r.db.Create(&models.Song{
		PlaylistID: uuid.MustParse(playlistId),
		Base: models.Base{
			ID: uuid.MustParse(songId),
		},
	}).Error
}

func (r *repository) DeletePlaylistSong(userId, playlistId, songId string) error {
	playlist := &models.Playlist{Base: models.Base{ID: uuid.MustParse(playlistId)}}
	err := r.db.First(playlist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierrors.ErrNotFound
		}
		return err
	}
	if playlist.UserId != userId {
		return apierrors.ErrForbidden
	}
	return r.db.Delete(&models.Song{}, "playlist_id = ? AND id = ?", playlistId, songId).Error
}
