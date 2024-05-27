package app

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oustrix/ozon_journal/config"
	"github.com/oustrix/ozon_journal/internal"
	"github.com/oustrix/ozon_journal/internal/controller/graphql"
	"github.com/oustrix/ozon_journal/internal/repository/inmemory"
	postgresRepository "github.com/oustrix/ozon_journal/internal/repository/postgres"
	"github.com/oustrix/ozon_journal/internal/service"
	"github.com/oustrix/ozon_journal/pkg/httpserver"
	"github.com/oustrix/ozon_journal/pkg/logger"
	"github.com/oustrix/ozon_journal/pkg/postgres"
)

// Run starts the application.
func Run(cfg *config.Config) {
	// Logger
	log := logger.New(cfg.Log.Level)
	log.Debug("Logger initialized", "level", cfg.Log.Level)
	log.Info("Starting application")

	// Repositories
	log.Info("Creating repositories", "storage", cfg.Storage.Type)

	var postRepo internal.PostRepository
	var commentRepo internal.CommentRepository

	if cfg.Storage.Type == "in-memory" {
		log.Debug("Using in-memory storage")
		postRepo = inmemory.NewPostRepository(log)
		commentRepo = inmemory.NewCommentRepository(log)
	} else if cfg.Storage.Type == "postgres" {
		log.Debug("Using postgres storage", "maxPoolSize", cfg.Postgres.MaxPoolSize,
			"connAttempts", cfg.Postgres.ConnAttempts, "connTimeout", cfg.Postgres.ConnTimeout)

		pg, err := postgres.New(cfg.Postgres.DSN,
			postgres.MaxPoolSize(int(cfg.Postgres.MaxPoolSize)),
			postgres.ConnAttempts(int(cfg.Postgres.ConnAttempts)),
			postgres.ConnTimeout(time.Duration(cfg.Postgres.ConnTimeout)*time.Second))
		if err != nil {
			log.Error("Failed to connect to postgres", "error", err.Error())
			return
		}
		defer pg.Close()

		log.Info("Migrating database")
		err = migrateUp(&cfg.Postgres)
		if err != nil {
			log.Error("Failed to apply migrations", "error", err.Error())
			return
		}
		log.Info("Database migrated")

		postRepo = postgresRepository.NewPostRepository(pg, log)
		commentRepo = postgresRepository.NewCommentRepository(pg, log)
	} else {
		log.Error("Unknown storage type", "type", cfg.Storage.Type)
		return
	}
	log.Info("Repositories created")

	// Services
	log.Info("Creating services")
	postService := service.NewPostService(postRepo, &cfg.Post, log)
	commentService := service.NewCommentService(commentRepo, &cfg.Comment, log)
	log.Info("Services created")

	// Router
	var router http.Handler
	log.Debug("Creating router", "environment", cfg.Environment)
	if cfg.Environment == "development" {
		router = graphql.NewRouter(log, true, commentService, postService)
	} else {
		router = graphql.NewRouter(log, false, commentService, postService)

	}
	log.Debug("Router created")

	// HTTP server
	log.Debug("Creating http server", "port", cfg.HTTP.Port)
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))
	log.Info("HTTP server started", "port", cfg.HTTP.Port)

	// Interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Waiting for signal
	select {
	case s := <-interrupt:
		log.Info("Got interrupt signal", "signal", s.String())
	case err := <-httpServer.Notify():
		log.Error("Got error while serving http", "error", err.Error())
	}

	// Shutdown
	log.Info("Shutting down HTTP server")
	err := httpServer.Shutdown()
	if err != nil {
		log.Error("Got error while shutting down http server", "error", err.Error())
	} else {
		log.Info("HTTP server stopped")
	}
}
