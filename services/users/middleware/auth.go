package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// NewJWKSet creates an auto-refreshing key set to validate JWT signatures.
func NewJWKSet(jwkUrl string) jwk.Set {
	jwkCache := jwk.NewCache(context.Background())

	// register a minimum refresh interval for this URL.
	// when not specified, defaults to Cache-Control and similar resp headers
	err := jwkCache.Register(jwkUrl, jwk.WithMinRefreshInterval(10*time.Minute))
	if err != nil {
		panic(fmt.Errorf("failed to register jwk location: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch once on application startup
	_, err = jwkCache.Refresh(ctx, jwkUrl)
	if err != nil {
		panic(fmt.Errorf("failed to fetch jwk on startup: %v", err))
	}
	// create the cached key set
	return jwk.NewCachedSet(jwkCache, jwkUrl)
}

// NewAuthMiddleware creates a middleware that will authorize requests based on
// the required scopes for the operation.
func NewAuthMiddleware(api huma.API, auth0Domain string, auth0Audience string) func(ctx huma.Context, next func(huma.Context)) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", auth0Domain)
	keySet := NewJWKSet(jwksURL)
	issuerUrl := fmt.Sprintf("https://%s/", auth0Domain)

	return func(ctx huma.Context, next func(huma.Context)) {
		var anyOfNeededScopes []string
		isAuthorizationRequired := false
		for _, opScheme := range ctx.Operation().Security {
			var ok bool
			if anyOfNeededScopes, ok = opScheme["auth0"]; ok {
				isAuthorizationRequired = true
				break
			}
		}

		if !isAuthorizationRequired {
			next(ctx)
			return
		}

		token := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")
		if len(token) == 0 {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}
		// Parse and validate the JWT.
		parsed, err := jwt.ParseString(token,
			jwt.WithKeySet(keySet),
			jwt.WithValidate(true),
			jwt.WithIssuer(issuerUrl),
			jwt.WithAudience(auth0Audience),
		)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Ensure the claims required for this operation are present.
		scopes, _ := parsed.Get("scope")
		if scopes, ok := scopes.(string); ok {
			for _, scope := range strings.Fields(scopes) {
				if slices.Contains(anyOfNeededScopes, scope) {
					next(ctx)
					return
				}
			}
		}

		huma.WriteErr(api, ctx, http.StatusForbidden, "Forbidden")
	}
}
