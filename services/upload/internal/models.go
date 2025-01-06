package internal

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

type Artist struct {
	Base
	Name   string
	Bio    string
	Images []Image
	Songs  []*Song `gorm:"many2many:artist_songs;"`
}

type Image struct {
	Base
	BlobUrl  string
	ArtistId uuid.UUID
}

type Song struct {
	Base
	Title      string
	BlobUrl    string
	TrackOrder uint
	AlbumId    uuid.UUID
	Artists    []*Artist `gorm:"many2many:artist_songs;"`
}

type Album struct {
	Base
	Title string
	Songs []Song
}

// type Genre struct {
// 	base
// 	Name string
// }

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&Artist{}, &Album{}, &Image{}, &Song{})
	if err != nil {
		panic(err)
	}
}
