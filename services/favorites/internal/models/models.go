package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time // Automatically managed by GORM for update time
}

type Playlist struct {
	Base
	UserId string
	Name   string
	Public bool
	Songs  []*Song // TODO: maybe inline?
}

type Song struct {
	Base
	PlaylistID uuid.UUID `gorm:"primaryKey;"`
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
