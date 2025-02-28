package songs

import (
	"soundscape/services/favorites/internal"
	"soundscape/services/favorites/internal/dtos"

	"github.com/sirupsen/logrus"
)

type Service interface {
	GetLikedSongs(token, userId string) ([]*dtos.Song, error)
	LikeSong(userId string, songId string) error
	UnlikeSong(userId string, songId string) error
	CheckSongs(userId string, songIds []string) ([]bool, error)
}

type service struct {
	repo     Repository
	metadata internal.MetadataRepository
}

var _ Service = (*service)(nil)

func NewService(repo Repository, metadataRepo internal.MetadataRepository) *service {
	return &service{repo: repo, metadata: metadataRepo}
}

func (s *service) GetLikedSongs(token, userId string) ([]*dtos.Song, error) {
	songDtos := []*dtos.Song{}
	songs, err := s.repo.GetLikedSongs(userId)
	if err != nil {
		return nil, err
	}
	songIds := []string{}
	for _, song := range songs {
		songIds = append(songIds, song.ID.String())
	}
	if len(songIds) == 0 {
		return songDtos, nil
	}
	songsFull, err := s.metadata.GetSongs(token, songIds)
	if err != nil {
		return nil, err
	}
	for i, song := range songsFull {
		logrus.Println(song.Title, songs[i].CreatedAt)
		songDtos = append(songDtos, &dtos.Song{
			Id:        song.ID,
			Title:     song.Title,
			CreatedAt: songs[i].CreatedAt.String(),
		})
	}
	return songDtos, nil
}

func (s *service) LikeSong(userId string, songId string) error {
	return s.repo.LikeSong(userId, songId)
}

func (s *service) UnlikeSong(userId string, songId string) error {
	return s.repo.UnlikeSong(userId, songId)
}

func (s *service) CheckSongs(userId string, songIds []string) ([]bool, error) {
	return s.repo.CheckSongs(userId, songIds)
}
