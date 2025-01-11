package playlists

import (
	"context"
	"errors"
	"net/http"
	"soundscape/services/favorites/internal"
	"soundscape/services/favorites/internal/apierrors"
	"soundscape/services/favorites/internal/dtos"
	"soundscape/shared"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB) {
	service := NewService(
		NewRepository(db),
		internal.NewMetadataRepository(),
	)
	registerCreatePlaylist(api, service)
	registerGetPlaylists(api, service)
	registerGetPlaylist(api, service)
	registerCreatePlaylistSong(api, service)
	registerDeletePlaylistSong(api, service)
}

type CreatePlaylistInput struct {
	Body dtos.CreatePlaylist
}
type CreatePlaylistOutput struct {
	Body dtos.Playlist
}

func registerCreatePlaylist(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "create-user-playlist",
		Method:      http.MethodPost,
		Path:        "/playlists",
		Summary:     "Create a user playlist",
		Description: "Create a playlist for a user.",
		Tags:        []string{"Playlists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CreatePlaylistInput) (*CreatePlaylistOutput, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		playlist, err := service.CreatePlaylist(userId, input.Body)
		if err != nil {
			return nil, err
		}
		return &CreatePlaylistOutput{
			Body: *playlist,
		}, nil
	})
}

type GetPlaylistsOutput struct {
	Body []*dtos.Playlist
}

func registerGetPlaylists(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user-playlists",
		Method:      http.MethodGet,
		Path:        "/playlists",
		Summary:     "Get users playlists",
		Description: "Get users playlists.",
		Tags:        []string{"Playlists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *struct{}) (*GetPlaylistsOutput, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		playlists, err := service.GetPlaylists(userId)
		if err != nil {
			return nil, err
		}
		return &GetPlaylistsOutput{
			Body: playlists,
		}, nil
	})
}

type GetPlaylistsInput struct {
	PlaylistId string `path:"playlistId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Playlist ID"`
}
type GetPlaylistOutput struct {
	Body *dtos.PlaylistFull
}

func registerGetPlaylist(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user-playlist",
		Method:      http.MethodGet,
		Path:        "/playlists/{playlistId}",
		Summary:     "Get users playlist",
		Description: "Get users playlist.",
		Tags:        []string{"Playlists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetPlaylistsInput) (*GetPlaylistOutput, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		token := ctx.Value(shared.RawTokenKey{}).(string)
		logrus.Println("getting a playlist", input.PlaylistId)
		playlist, err := service.GetPlaylist(token, userId, input.PlaylistId)
		if err != nil {
			return nil, err
		}
		return &GetPlaylistOutput{
			Body: playlist,
		}, nil
	})
}

type CreatePlaylistSongInput struct {
	PlaylistId string `path:"playlistId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Playlist ID"`
	Body       struct {
		SongId string `json:"song_id" example:"550e8400-e29b-41d4-a716-446655440000" description:"Song ID"`
	}
}

func registerCreatePlaylistSong(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "create-user-playlist-song",
		Method:      http.MethodPost,
		Path:        "/playlists/{playlistId}/songs",
		Summary:     "Add a song to a user playlist",
		Description: "Add a song to a user playlist.",
		Tags:        []string{"Playlists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CreatePlaylistSongInput) (*struct{}, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		if err := service.CreatePlaylistSong(userId, input.PlaylistId, input.Body.SongId); err != nil {
			if errors.Is(err, apierrors.ErrNotFound) {
				return nil, huma.Error404NotFound("Playlist not found")
			}
			if errors.Is(err, apierrors.ErrForbidden) {
				return nil, huma.Error403Forbidden("Unauthorized")
			}
			return nil, err
		}
		return nil, nil
	})
}

type DeletePlaylistSongInput struct {
	PlaylistId string `path:"playlistId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Playlist ID"`
	SongId     string `path:"songId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Song ID"`
}

func registerDeletePlaylistSong(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "delete-user-playlist-song",
		Method:      http.MethodDelete,
		Path:        "/playlists/{playlistId}/songs/{songId}",
		Summary:     "Delete a song from a user playlist",
		Description: "Delete a song from a user playlist.",
		Tags:        []string{"Playlists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *DeletePlaylistSongInput) (*struct{}, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		if err := service.DeletePlaylistSong(userId, input.PlaylistId, input.SongId); err != nil {
			if errors.Is(err, apierrors.ErrNotFound) {
				return nil, huma.Error404NotFound("Playlist not found")
			}
			if errors.Is(err, apierrors.ErrForbidden) {
				return nil, huma.Error403Forbidden("Unauthorized")
			}
		}
		return nil, nil
	})
}
