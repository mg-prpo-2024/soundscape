package internal

import (
	"time"
)

type Service interface {
	CreateArtist(artist CreateArtistDto) error
	GetArtist(userId string) (*ArtistDto, error)
	CreateAlbum(album CreateAlbumDto) (*AlbumDto, error)
	GetAlbum(id string) (*AlbumDto, error)
	CreateSong(song CreateSongDto) (*CreateSongResponseDto, error)
	GetAlbumSongs(id string) ([]*SongDto, error)
	DeleteSong(id string) error
}

type service struct {
	repo    Repository
	storage Storage
}

var _ Service = (*service)(nil)

func NewService(repo Repository, storage Storage) *service {
	return &service{repo: repo, storage: storage}
}

func (s *service) CreateArtist(artist CreateArtistDto) error {
	return s.repo.CreateArtist(artist)
}

func (s *service) GetArtist(userId string) (*ArtistDto, error) {
	artist, err := s.repo.GetArtist(userId)
	if err != nil {
		return nil, err
	}
	return &ArtistDto{
		Id: artist.ID.String(),
		CreateArtistDto: CreateArtistDto{
			UserId: artist.UserId,
			Name:   artist.Name,
			Bio:    artist.Bio,
		},
	}, nil
}

func (s *service) CreateAlbum(albumDto CreateAlbumDto) (*AlbumDto, error) {
	album, err := s.repo.CreateAlbum(albumDto)
	if err != nil {
		return nil, err
	}
	return &AlbumDto{
		Id: album.ID.String(),
		CreateAlbumDto: CreateAlbumDto{
			Title: album.Title,
		},
	}, err
}

func (s *service) GetAlbum(id string) (*AlbumDto, error) {
	album, err := s.repo.GetAlbum(id)
	if err != nil {
		return nil, err
	}
	return &AlbumDto{
		Id: album.ID.String(),
		CreateAlbumDto: CreateAlbumDto{
			Title: album.Title,
		},
	}, nil
}

func (s *service) CreateSong(songDto CreateSongDto) (*CreateSongResponseDto, error) {
	song, err := s.repo.CreateSong(songDto)
	if err != nil {
		return nil, err
	}
	id := song.ID.String()
	presignedUrl, err := s.storage.GeneratePresignedURL(id, 5*time.Hour)
	if err != nil {
		return nil, err
	}
	return &CreateSongResponseDto{
		Id:       id,
		UpoadUrl: presignedUrl,
	}, nil
}

func (s *service) GetAlbumSongs(albumId string) ([]*SongDto, error) {
	songs, err := s.repo.GetAlbumSongs(albumId)
	if err != nil {
		return nil, err
	}
	songDtos := []*SongDto{}
	for _, song := range songs {
		songDtos = append(songDtos, &SongDto{
			Id:    song.ID.String(),
			Title: song.Title,
		})
	}
	return songDtos, nil
}

func (s *service) DeleteSong(id string) error {
	// TODO: this should be transactional
	err := s.storage.DeleteFile(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteSong(id)
}
