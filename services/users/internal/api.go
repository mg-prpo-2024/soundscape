package internal

import (
	"context"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/danielgtaylor/huma/v2"
)

func Register(api huma.API) {
	RegisterSignIn(api)
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

func RegisterSignIn(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "sign-in",
		Method:      http.MethodPost,
		Path:        "/sign-in",
		Summary:     "Sign in",
		Description: "Sign in to the system.",
		Tags:        []string{"Sign In"},
	}, func(ctx context.Context, input *SignInInput) (*SignInOutput, error) {
		resp := &SignInOutput{}
		token := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		fmt.Printf("token %+v", token)
		// err := Login("", token.RegisteredClaims.Subject)
		// if err != nil {
		// 	return nil, err
		// }
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Body.Name)
		return resp, nil
	})
}
