package albums

import (
	"soundscape/services/metadata/internal/dtos"
)

type Service interface {
	GetAlbums() ([]*dtos.Album, error)
	GetAlbum(id string) (*dtos.Album, error)
	GetAlbumSongs(id string) ([]*dtos.Song, error)
}

type service struct {
	repo Repository
	// storage Storage
}

var _ Service = (*service)(nil)

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) GetAlbums() ([]*dtos.Album, error) {
	albums, err := s.repo.GetAlbums()
	if err != nil {
		return nil, err
	}
	albumDtos := []*dtos.Album{}
	for _, album := range albums {
		artist := dtos.Artist{
			Id:   album.Artist.ID.String(),
			Name: album.Artist.Name,
			Bio:  album.Artist.Bio,
		}
		albumDtos = append(albumDtos, &dtos.Album{
			Id:        album.ID.String(),
			Title:     album.Title,
			CreatedAt: album.CreatedAt,
			Artist:    artist,
		})
	}
	return albumDtos, nil
}

func (s *service) GetAlbum(id string) (*dtos.Album, error) {
	album, err := s.repo.GetAlbum(id)
	if err != nil {
		return nil, err
	}
	return &dtos.Album{
		Id:        album.ID.String(),
		Title:     album.Title,
		CreatedAt: album.CreatedAt,
		Artist: dtos.Artist{
			Id:   album.Artist.ID.String(),
			Name: album.Artist.Name,
			Bio:  album.Artist.Bio,
		},
	}, nil
}

func (s *service) GetAlbumSongs(albumId string) ([]*dtos.Song, error) {
	songs, err := s.repo.GetAlbumSongs(albumId)
	if err != nil {
		return nil, err
	}
	songDtos := []*dtos.Song{}
	for _, song := range songs {
		songDtos = append(songDtos, &dtos.Song{
			Id:    song.ID.String(),
			Title: song.Title,
		})
	}
	return songDtos, nil
}
