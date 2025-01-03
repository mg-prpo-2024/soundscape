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
		router.Use(middleware.EnsureValidToken())

		api := humachi.New(router, huma.DefaultConfig("Users API", "1.0.0"))

		db := connect(options)
		internal.Migrate(db)

		registerHealthCheck(router, db)
		internal.Register(api, db, options)

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: router,
		}

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			server.ListenAndServe()
		})
		hooks.OnStop(func() {
			// Gracefull shutdown
			fmt.Printf("Shutting down server...\n")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}

func connect(opts *internal.Options) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		opts.PostgresHost,
		opts.PostgresUser,
		opts.PostgresPassword,
		opts.PostgresDB,
		opts.PostgresPort,
	)
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
