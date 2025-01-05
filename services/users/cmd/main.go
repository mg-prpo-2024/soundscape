package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"users/internal"
	"users/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *internal.Options) {
		// Create a new router & API
		router := chi.NewMux()
		router.Use(chimiddleware.RequestID)
		router.Use(chimiddleware.Recoverer)
		router.Use(chimiddleware.Logger)
		router.Use(cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"https://*", "http://*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

		db := connect(options)
		internal.AutoMigrate(db)

		registerHealthCheck(router, db)
		registerApi(router, db, options)

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: router,
		}

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			logrus.Infof("Starting server on port %d...\n", options.Port)
			server.ListenAndServe()
		})
		hooks.OnStop(func() {
			// Gracefull shutdown
			logrus.Infof("Shutting down server...\n")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	})

	cli.Run()
}

func connect(opts *internal.Options) *gorm.DB {
	internal.PrettyPrint(opts)
	dsn := opts.PostgresUrl
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func registerHealthCheck(router chi.Router, db *gorm.DB) {
	router.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	router.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			http.Error(w, "not ready: "+err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	})
}

func registerApi(router chi.Router, db *gorm.DB, options *internal.Options) {
	config := huma.DefaultConfig("Users API", "1.0.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"auth0": {
			Type: "oauth2",
			Flows: &huma.OAuthFlows{
				AuthorizationCode: &huma.OAuthFlow{
					AuthorizationURL: fmt.Sprintf("https://%s/authorize", options.Auth0Domain),
					TokenURL:         fmt.Sprintf("https://%s/oauth/token", options.Auth0Domain),
					Scopes: map[string]string{
						"openid": "openid scope description...",
					},
				},
			},
		},
		"Auth0WebhookSecret": {
			Type: "apiKey",
			In:   "header",
			Name: "X-Auth0-Webhook-Secret",
		},
	}

	api := humachi.New(router, config)
	api.UseMiddleware(middleware.NewAuthMiddleware(api, options.Auth0Domain, options.Auth0Audience))
	internal.Register(api, db, options)
}
