package metadatadb

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
	UserId string
	Name   string
	Bio    string
	Images []Image
	Songs  []*Song  `gorm:"constraint:OnDelete:CASCADE;many2many:artist_songs;"`
	Albums []*Album `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
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
	Artists    []*Artist `gorm:"constraint:OnDelete:CASCADE;many2many:artist_songs;"`
}

type Album struct {
	Base
	Title    string
	Songs    []Song
	ArtistId uuid.UUID
	Artist   Artist
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
