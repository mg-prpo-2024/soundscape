package songs

import (
	"context"
	"net/http"
	"soundscape/services/favorites/internal"
	"soundscape/services/favorites/internal/dtos"
	"soundscape/shared"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB) {
	service := NewService(
		NewRepository(db),
		internal.NewMetadataRepository(),
	)
	registerGetSongs(api, service)
	registerLikeSong(api, service)
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
	SongID string `path:"songId" example:"" doc:"Song ID"`
}
type LikeSongOutput struct {
}

func registerLikeSong(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "like-song",
		Method:      http.MethodPost,
		Path:        "/songs/{songId}",
		Summary:     "Like a song",
		Description: "Like a song.",
		Tags:        []string{"Songs"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *LikeSongInput) (*LikeSongOutput, error) {
		userId := ctx.Value(shared.TokenKey{}).(jwt.Token).Subject()
		err := service.LikeSong(userId, input.SongID)
		if err != nil {
			return nil, err
		}
		return &LikeSongOutput{}, nil
	})
}
