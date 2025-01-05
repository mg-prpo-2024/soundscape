package internal

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB, config *Config) {
	service := NewService(NewRepository(db))
	registerCreateArtist(api, service)
	// registerGetCustomer(api, service)
}

type CreateArtistInput struct {
	Secret string `header:"X-Auth0-Webhook-Secret" doc:"Auth0 Webhook Secret"`
	Body   ArtistDto
}

type CreateArtistOutput struct{}

func registerCreateArtist(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "create-artist",
		Method:      http.MethodPost,
		Path:        "/artists",
		Summary:     "Create an artist",
		Description: "Create an artist.",
		Tags:        []string{"Artists"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *CreateArtistInput) (*CreateArtistOutput, error) {
		err := service.CreateArtist(input.Body)
		if err != nil {
			return nil, err
		}
		return &CreateArtistOutput{}, nil
	})
}
