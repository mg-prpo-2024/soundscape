package songs

import "soundscape/services/metadata/internal/dtos"

type Service interface {
	GetSongs(songIds []string) ([]*dtos.SongFull, error)
}

type service struct {
	repo Repository
	// storage Storage
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) GetSongs(songIds []string) ([]*dtos.SongFull, error) {
	songs, err := s.repo.GetSongs(songIds)
	if err != nil {
		return nil, err
	}
	songDtos := []*dtos.SongFull{}
	for _, song := range songs {
		songDtos = append(songDtos, &dtos.SongFull{
			Id:    song.ID.String(),
			Title: song.Title,
		})
	}
	return songDtos, nil
}
