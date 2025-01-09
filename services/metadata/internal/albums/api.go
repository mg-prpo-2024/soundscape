package albums

import (
	"context"
	"net/http"
	"soundscape/services/metadata/internal/dtos"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB) {
	service := NewService(
		NewRepository(db),
	)
	registerGetAlbum(api, service)
	registerGetSongs(api, service)
}

type GetAlbumInput struct {
	Id string `path:"id" doc:"Album ID" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type GetAlbumOutput struct {
	Body dtos.Album
}

func registerGetAlbum(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-album",
		Method:      http.MethodGet,
		Path:        "/albums/{id}",
		Summary:     "Get an album",
		Description: "Get the album metadata.",
		Tags:        []string{"Albums"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetAlbumInput) (*GetAlbumOutput, error) {
		album, err := service.GetAlbum(input.Id)
		if err != nil {
			return nil, err
		}
		return &GetAlbumOutput{
			Body: *album,
		}, nil
	})
}

type GetSongsInput struct {
	AlbumId string `path:"albumId" doc:"Album ID" example:"550e8400-e29b-41d4-a716-446655440000"`
}
type GetSongsOutput struct {
	Body []*dtos.Song
}

func registerGetSongs(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-album-songs",
		Method:      http.MethodGet,
		Path:        "/albums/{albumId}/songs",
		Summary:     "Get album songs",
		Description: "Get album songs.",
		Tags:        []string{"Albums", "Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetSongsInput) (*GetSongsOutput, error) {
		songs, err := service.GetAlbumSongs(input.AlbumId)
		if err != nil {
			return nil, err
		}
		return &GetSongsOutput{
			Body: songs,
		}, nil
	})
}
