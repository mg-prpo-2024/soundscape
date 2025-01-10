package playlists

import (
	"context"
	"net/http"
	"soundscape/services/favorites/internal/dtos"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB) {
	service := NewService(
		NewRepository(db),
	)
	registerGetAlbums(api, service)
}

type GetPlaylistsOutput struct {
	Body []*dtos.Playlist
}

func registerGetAlbums(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-playlists",
		Method:      http.MethodGet,
		Path:        "/playlists",
		Summary:     "Get playlists",
		Description: "Get playlists.",
		Tags:        []string{"Playlists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *struct{}) (*GetPlaylistsOutput, error) {
		playlists, err := service.GetPlaylists()
		if err != nil {
			return nil, err
		}
		return &GetPlaylistsOutput{
			Body: playlists,
		}, nil
	})
}
