package playlists

import (
	"soundscape/services/favorites/internal"
	"soundscape/services/favorites/internal/dtos"

	"github.com/sirupsen/logrus"
)

type Service interface {
	CreatePlaylist(userId string, playlist dtos.CreatePlaylist) (*dtos.Playlist, error)
	GetPlaylists(userId string) ([]*dtos.Playlist, error)
	GetPlaylist(token, userId, playlistId string) (*dtos.PlaylistFull, error)
	CreatePlaylistSong(userId, playlistId, songId string) error
	DeletePlaylistSong(userid, playlistId, songId string) error
}

type service struct {
	repo     Repository
	metadata internal.MetadataRepository
}

var _ Service = (*service)(nil)

func NewService(repo Repository, metadataRepo internal.MetadataRepository) *service {
	return &service{repo: repo, metadata: metadataRepo}
}

func (s *service) CreatePlaylist(userId string, playlistDto dtos.CreatePlaylist) (*dtos.Playlist, error) {
	playlist, err := s.repo.CreatePlaylist(userId, playlistDto)
	if err != nil {
		return nil, err
	}
	return &dtos.Playlist{
		Id:     playlist.ID.String(),
		UserId: userId,
		Name:   playlist.Name,
	}, nil
}

func (s *service) GetPlaylists(userId string) ([]*dtos.Playlist, error) {
	playlists, err := s.repo.GetPlaylists(userId)
	if err != nil {
		return nil, err
	}
	playlistDtos := []*dtos.Playlist{}
	for _, playlist := range playlists {
		playlistDtos = append(playlistDtos, &dtos.Playlist{
			Id:         playlist.ID.String(),
			UserId:     userId,
			Name:       playlist.Name,
			TotalSongs: len(playlist.Songs),
		})
	}
	return playlistDtos, nil
}

func (s *service) GetPlaylist(token, userId, playlistId string) (*dtos.PlaylistFull, error) {
	playlist, err := s.repo.GetPlaylist(userId, playlistId)
	if err != nil {
		return nil, err
	}
	logrus.Println("playlist", playlist)
	songIds := []string{}
	for _, song := range playlist.Songs {
		songIds = append(songIds, song.ID.String())
	}
	songs := []*dtos.Song{}
	if len(songIds) > 0 {
		resSongs, err := s.metadata.GetSongs(token, songIds)
		if err != nil {
			logrus.Errorf("error getting songs from metadata api: %v", err)
			return nil, err
		}
		for _, song := range resSongs {
			songs = append(songs, &dtos.Song{
				Id:    song.ID,
				Title: song.Title,
			})
		}
	}
	return &dtos.PlaylistFull{
		Id:     playlist.ID.String(),
		UserId: userId,
		Name:   playlist.Name,
		Songs:  songs,
	}, nil
}

func (s *service) CreatePlaylistSong(userId, playlistId, songId string) error {
	return s.repo.CreatePlaylistSong(userId, playlistId, songId)
}

func (s *service) DeletePlaylistSong(userid, playlistId, songId string) error {
	return s.repo.DeletePlaylistSong(userid, playlistId, songId)
}
