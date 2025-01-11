package songs

import (
	"soundscape/shared/metadatadb"

	"gorm.io/gorm"
)

type Repository interface {
	GetSongs(songIds []string) ([]*metadatadb.Song, error)
}

type repository struct {
	db *gorm.DB
}

var _ Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetSongs(songIds []string) ([]*metadatadb.Song, error) {
	var songs []*metadatadb.Song
	err := r.db.Where("id in (?)", songIds).Find(&songs).Error
	// Sort songs to match input order
	songMap := make(map[string]*metadatadb.Song)
	for _, song := range songs {
		songMap[song.ID.String()] = song
	}
	sortedSongs := make([]*metadatadb.Song, 0, len(songIds))
	for _, id := range songIds {
		if song, ok := songMap[id]; ok {
			sortedSongs = append(sortedSongs, song)
		}
	}
	return sortedSongs, err
}
