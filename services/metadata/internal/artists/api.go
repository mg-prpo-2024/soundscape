package artists

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
	registerGetArtist(api, service)
	registerGetArtistAlbums(api, service)
}

type GetArtistInput struct {
	UserId string `path:"userId" doc:"User ID" example:"google-oauth2|106527649641850458478"`
}

type GetArtistOutput struct {
	Body dtos.Artist
}

func registerGetArtist(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-artist",
		Method:      http.MethodGet,
		Path:        "/artists/{userId}",
		Summary:     "Get an artist",
		Description: "Get the artist metadata for the user.",
		Tags:        []string{"Artists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetArtistInput) (*GetArtistOutput, error) {
		artist, err := service.GetArtist(input.UserId)
		if err != nil {
			return nil, err
		}
		return &GetArtistOutput{
			Body: *artist,
		}, nil
	})
}

type GetArtistAlbumsInput struct {
	ArtistId string `path:"artistId" doc:"Artist ID" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type GetArtistAlbumsOutput struct {
	Body []*dtos.Album
}

func registerGetArtistAlbums(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-artist-albums",
		Method:      http.MethodGet,
		Path:        "/artists/{artistId}/albums",
		Summary:     "Get artist albums",
		Description: "Get the albums for the artist.",
		Tags:        []string{"Artists", "Albums"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetArtistAlbumsInput) (*GetArtistAlbumsOutput, error) {
		albums, err := service.GetArtistAlbums(input.ArtistId)
		if err != nil {
			return nil, err
		}
		return &GetArtistAlbumsOutput{
			Body: albums,
		}, nil
	})
}
