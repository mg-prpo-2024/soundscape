package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Playlist struct {
	ID     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId string
	Name   string
	Public bool
	Songs  []*Song // TODO: maybe inline?
}

type Song struct {
	ID         uuid.UUID `gorm:"primaryKey;type:uuid"`
	PlaylistID uuid.UUID `gorm:"primaryKey"`
	CreatedAt  time.Time // Automatically managed by GORM for creation time
	UpdatedAt  time.Time // Automatically managed by GORM for update time
}

type LikedSong struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserId    string    `gorm:"primaryKey"`
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time // Automatically managed by GORM for update time
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&Playlist{}, &Song{}, &LikedSong{})
	if err != nil {
		panic(err)
	}
}
