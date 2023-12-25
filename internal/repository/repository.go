package repository

import (
	"context"
	"fmt"

	"github.com/bdrbt/todo/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	pgPool *pgxpool.Pool
	logger *zap.Logger
	// queries map[string]string
}

// New Repository instance.
func New(cfg *config.Config, lg *zap.Logger) (*Repository, error) {
	pg, err := initDB(context.Background(), cfg.PgURL())
	if err != nil {
		return nil, fmt.Errorf("cannot init pgx pool:%w", err)
	}

	repo := &Repository{
		pgPool: pg,
		logger: lg,
	}

	return repo, nil
}

// initDB initialize connection according provided pgURL, attempt to validate pgURL,
// connect to DB and Ping DB.
func initDB(ctx context.Context, pgURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(pgURL)
	if err != nil {
		return nil, fmt.Errorf("invalid postgres connection URL:%w", err)
	}

	pg, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to postgres:%w", err)
	}

	// ensure it's really available
	if err = pg.Ping(ctx); err != nil {
		return nil, fmt.Errorf("cannot ping :%w", err)
	}

	return pg, nil
}

func (repo *Repository) Close() {
	repo.pgPool.Close()
}
