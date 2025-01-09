package artists

import "soundscape/services/metadata/internal/dtos"

type Service interface {
	GetArtist(userId string) (*dtos.Artist, error)
	GetArtistAlbums(artistId string) ([]*dtos.Album, error)
}

type service struct {
	repo Repository
	// storage Storage
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) GetArtist(userId string) (*dtos.Artist, error) {
	artist, err := s.repo.GetArtist(userId)
	if err != nil {
		return nil, err
	}
	return &dtos.Artist{
		Id:     artist.ID.String(),
		UserId: artist.UserId,
		Name:   artist.Name,
		Bio:    artist.Bio,
	}, nil
}

func (s *service) GetArtistAlbums(artistId string) ([]*dtos.Album, error) {
	albums, err := s.repo.GetArtistAlbums(artistId)
	if err != nil {
		return nil, err
	}
	albumDtos := []*dtos.Album{}
	for _, album := range albums {
		albumDtos = append(albumDtos, &dtos.Album{
			Id:        album.ID.String(),
			Title:     album.Title,
			CreatedAt: album.CreatedAt,
		})
	}
	return albumDtos, nil
}
