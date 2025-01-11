package songs

import (
	"soundscape/services/favorites/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetLikedSongs(userId string) ([]*models.LikedSong, error)
	LikeSong(userId, songId string) error
	UnlikeSong(userId, songId string) error
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetLikedSongs(userId string) ([]*models.LikedSong, error) {
	songs := []*models.LikedSong{}
	err := r.db.Model(&models.LikedSong{}).Where("user_id = ?", userId).Find(&songs).Error
	return songs, err
}

func (r *repository) LikeSong(userId, songId string) error {
	return r.db.Create(&models.LikedSong{ID: uuid.MustParse(songId), UserId: userId}).Error
}

func (r *repository) UnlikeSong(userId, songId string) error {
	return r.db.Delete(&models.LikedSong{ID: uuid.MustParse(songId), UserId: userId}).Error
}
