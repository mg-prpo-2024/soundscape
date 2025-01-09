package internal

import (
	"soundscape/services/metadata/internal/dtos"
	"time"
)

type Service interface {
	GetArtist(userId string) (*dtos.Artist, error)
	GetArtistAlbums(artistId string) ([]*dtos.Album, error)
	GetAlbum(id string) (*dtos.Album, error)
	GetAlbumSongs(id string) ([]*dtos.Song, error)
}

type service struct {
	repo    Repository
	storage Storage
}

var _ Service = (*service)(nil)

func NewService(repo Repository, storage Storage) *service {
	return &service{repo: repo, storage: storage}
}

func (s *service) GetArtist(userId string) (*dtos.Artist, error) {
	artist, err := s.repo.GetArtist(userId)
	if err != nil {
		return nil, err
	}
	return &dtos.Artist{
		Id: artist.ID.String(),
		CreateArtist: dtos.CreateArtist{
			UserId: artist.UserId,
			Name:   artist.Name,
			Bio:    artist.Bio,
		},
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
			CreatedAt: album.CreatedAt.Format(time.RFC3339),
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
		Id:    album.ID.String(),
		Title: album.Title,
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
