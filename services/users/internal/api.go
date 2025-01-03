package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Register(api huma.API, db *gorm.DB, options *Options) {
	service := NewService(NewRepository(db))
	RegisterSignIn(api, service)
}

type SignInInput struct {
	Body struct {
		Name string `json:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}
}

type SignInOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func RegisterSignIn(api huma.API, service Service) {
	huma.Register(api, huma.Operation{
		OperationID: "sign-in",
		Method:      http.MethodPost,
		Path:        "/sign-in",
		Summary:     "Sign in",
		Description: "Sign in to the system.",
		Tags:        []string{"Sign In"},
		Security: []map[string][]string{
			{"auth0": {"openid"}},
		},
	}, func(ctx context.Context, input *SignInInput) (*SignInOutput, error) {
		resp := &SignInOutput{}
		// token := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		// fmt.Printf("token %+v", token)
		// err := Login("", token.RegisteredClaims.Subject)
		// if err != nil {
		// 	return nil, err
		// }
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Body.Name)
		// service.Login(token.Raw, token.RegisteredClaims.Subject)
		return resp, nil
	})
}
