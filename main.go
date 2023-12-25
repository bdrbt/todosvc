package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdrbt/todo/database"
	"github.com/bdrbt/todo/internal/config"
	"github.com/bdrbt/todo/internal/handlers"
	"github.com/bdrbt/todo/internal/logger"
	"github.com/bdrbt/todo/internal/repository"
	"github.com/bdrbt/todo/internal/usecases"
	"github.com/bdrbt/todo/internal/vcs"
	"go.uber.org/zap"
)

var version = vcs.Version()

// @title Insof Service

func main() {
	// initial config truct with Environment to determine our loglevel.
	cfg := config.New()

	// logger shared across handlers,repositry,usecases, whatever.
	logger := logger.New(cfg.Environment)

	// load rest of config from environment variables.
	// TODO add coffer.
	err := cfg.Load()
	if err != nil {
		logger.Panic("cannot load config", zap.Error(err))
	}

	// try to connect and migrate DB using golang-migrate.
	pgURL := cfg.PgURL()
	if err = database.MigrateDB(pgURL); err != nil {
		logger.Error("cannot migrate db", zap.Error(err))
		os.Exit(-1)
	}

	// setup repository
	repo, err := repository.New(cfg, logger)
	if err != nil {
		logger.Error("cannot create repository", zap.Error(err))
		os.Exit(-1)
	}

	// setup usecases
	ucs := usecases.New(repo, logger)

	// setup router
	router := handlers.MountHandlers(cfg, ucs, logger)

	// setup http server
	if err := serveHTTP(logger, cfg, router); err != nil {
		logger.Error("cannot start http server", zap.Error(err))
	}

	// final cleanup
	repo.Close()
	//nolint: errcheck
	logger.Sync()
}

func serveHTTP(logger *zap.Logger, cfg *config.Config, router http.Handler) error {
	restAPI := http.Server{
		Addr:         cfg.Rest.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Rest.ReadTimeout,
		WriteTimeout: cfg.Rest.WriteTimeout,
		IdleTimeout:  cfg.Rest.IdleTimeout,
	}

	logger.Info("startin service", zap.String("version", version))

	serverErrors := make(chan error, 1)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Sturt
	go func() {
		logger.Info("http server started", zap.String("addr", cfg.Rest.Addr))

		serverErrors <- restAPI.ListenAndServe()
	}()

	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Info("shutdown started", zap.String("signal", sig.String()))
		defer logger.Info("shutdown complete", zap.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Rest.ShutdownTimeout)
		defer cancel()

		if err := restAPI.Shutdown(ctx); err != nil {
			restAPI.Close()

			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
