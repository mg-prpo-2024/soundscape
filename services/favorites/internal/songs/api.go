package songs

import (
	"context"
	"net/http"
	"soundscape/services/favorites/internal"
	"soundscape/services/favorites/internal/dtos"
	"soundscape/shared"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB, config *internal.Config) {
	service := NewService(
		NewRepository(db),
		internal.NewMetadataRepository(config.MetadataServiceUrl),
	)
	registerCheckSongs(api, service)
	registerGetSongs(api, service)
	registerLikeSong(api, service)
	registerUnlikeSong(api, service)
}

type GetLikedSongsInput struct {
}
type GetLikedSongsOutput struct {
	Body []*dtos.Song
}

func registerGetSongs(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "get-liked-songs",
		Method:      http.MethodGet,
		Path:        "/songs",
		Summary:     "Get a user's liked songs",
		Description: "Get a user's liked songs.",
		Tags:        []string{"Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *GetLikedSongsInput) (*GetLikedSongsOutput, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		token := ctx.Value(shared.RawTokenKey{}).(string)
		songs, err := service.GetLikedSongs(token, userId)
		if err != nil {
			return nil, err
		}
		return &GetLikedSongsOutput{
			Body: songs,
		}, nil
	})
}

type LikeSongInput struct {
	Body struct {
		SongID string `json:"song_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Song ID"`
	}
}
type LikeSongOutput struct {
}

func registerLikeSong(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "like-song",
		Method:      http.MethodPost,
		Path:        "/songs",
		Summary:     "Like a song",
		Description: "Like a song.",
		Tags:        []string{"Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *LikeSongInput) (*LikeSongOutput, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		err := service.LikeSong(userId, input.Body.SongID)
		if err != nil {
			return nil, err
		}
		return &LikeSongOutput{}, nil
	})
}

type UnlikeSongInput struct {
	SongID string `path:"songId" example:"123e4567-e89b-12d3-a456-426614174000" doc:"Song ID"`
}
type UnlikeSongOutput struct {
}

func registerUnlikeSong(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "unlike-song",
		Method:      http.MethodDelete,
		Path:        "/songs/{songId}",
		Summary:     "Unlike a song",
		Description: "Unlike a song.",
		Tags:        []string{"Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *UnlikeSongInput) (*UnlikeSongOutput, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		err := service.UnlikeSong(userId, input.SongID)
		if err != nil {
			return nil, err
		}
		return &UnlikeSongOutput{}, nil
	})
}

type CheckSongsInput struct {
	SongIDs string `query:"ids" example:"123e4567-e89b-12d3-a456-426614174000,123e4567-e89b-12d3-a456-426614174001" doc:"Comma separated song IDs"`
}
type CheckSongsOutput struct {
	Body []bool
}

func registerCheckSongs(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "check-songs",
		Method:      http.MethodGet,
		Path:        "/songs/check",
		Summary:     "Check user's liked songs",
		Description: "Check if songs are liked by a user.",
		Tags:        []string{"Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CheckSongsInput) (*CheckSongsOutput, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		songIds := strings.Split(input.SongIDs, ",")
		checks, err := service.CheckSongs(userId, songIds)
		if err != nil {
			return nil, err
		}
		return &CheckSongsOutput{
			Body: checks,
		}, nil
	})
}
