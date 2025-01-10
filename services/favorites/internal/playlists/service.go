package playlists

import (
	"soundscape/services/favorites/internal/dtos"
)

type Service interface {
	GetPlaylists() ([]*dtos.Playlist, error)
}

type service struct {
	repo Repository
	// storage Storage
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) GetPlaylists() ([]*dtos.Playlist, error) {
	albums, err := s.repo.GetPlaylists()
	if err != nil {
		return nil, err
	}
	albumDtos := []*dtos.Playlist{}
	for _, album := range albums {
		albumDtos = append(albumDtos, &dtos.Playlist{
			Id: album.ID.String(),
			// Name: album.Name,
		})
	}
	return albumDtos, nil
}
