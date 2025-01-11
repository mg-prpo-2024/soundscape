package songs

import (
	"context"
	"net/http"
	"soundscape/services/metadata/internal/dtos"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB) {
	service := NewService(
		NewRepository(db),
	)
	registerGetSongs(api, service)
}

type GetArtistInput struct {
	SongIds string `query:"ids" doc:"The IDs of the songs to get." example:"123e4567-e89b-12d3-a456-426614174000,987fcdeb-51a2-43fe-ba98-237614172000"`
}

type GetArtistOutput struct {
	Body []*dtos.SongFull
}

func registerGetSongs(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-songs",
		Method:      http.MethodGet,
		Path:        "/songs",
		Summary:     "Get an artist",
		Description: "Get the artist metadata for the user.",
		Tags:        []string{"Artists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetArtistInput) (*GetArtistOutput, error) {
		songIds := strings.Split(input.SongIds, ",")
		songs, err := service.GetSongs(songIds)
		if err != nil {
			return nil, err
		}
		return &GetArtistOutput{
			Body: songs,
		}, nil
	})
}
