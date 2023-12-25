package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	EnvDevelopment = 4
	EnvTesting     = 2
	EnvProduction  = 0
)

var ErrConfig = errors.New("configuration error")

type Config struct {
	// Environment service running mode enum
	// one of EnvProduction, EnvTesting or EnvDevelopment.
	Environment int
	Postgres    struct {
		Host         string
		Port         int
		User         string
		Password     string
		Database     string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Rest struct {
		Addr            string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
	}
}

func New() *Config {
	cfg := &Config{}
	en := os.Getenv("ENVIRONMENT")

	switch en {
	case "development":
		cfg.Environment = EnvDevelopment
	case "testing":
		cfg.Environment = EnvTesting
	default:
		cfg.Environment = EnvProduction
	}

	return cfg
}

func (cfg *Config) Load() error {
	var err error

	cfg.Postgres.Host = os.Getenv("PG_HOST")
	cfg.Postgres.User = os.Getenv("PG_USER")
	cfg.Postgres.Password = os.Getenv("PG_PASSWORD")
	cfg.Postgres.Database = os.Getenv("PG_DATABASE")

	if cfg.Postgres.MaxIdleTime = os.Getenv("PG_IDLE_TIME"); cfg.Postgres.MaxIdleTime == "" {
		cfg.Postgres.MaxIdleTime = "10m"
	}

	if cfg.Postgres.Port, err = strconv.Atoi(os.Getenv("PG_PORT")); err != nil {
		cfg.Postgres.Port = 5432
	}

	if cfg.Postgres.MaxOpenConns, err = strconv.Atoi(os.Getenv("PG_OPEN_CONNS")); err != nil {
		cfg.Postgres.MaxOpenConns = 4
	}

	if cfg.Postgres.MaxIdleConns, err = strconv.Atoi(os.Getenv("PG_IDLE_CONNS")); err != nil {
		cfg.Postgres.MaxIdleConns = 4
	}

	if cfg.Rest.Addr = os.Getenv("ADDR"); cfg.Rest.Addr == "" {
		cfg.Rest.Addr = ":8080"
	}

	return cfg.Validate()
}

func (cfg *Config) Validate() error {
	if cfg.Postgres.Host == "" {
		return fmt.Errorf("%w: postgres host is not set", ErrConfig)
	}

	if cfg.Postgres.User == "" || cfg.Postgres.Password == "" {
		return fmt.Errorf("%w: postgrs username/password is not set", ErrConfig)
	}

	if cfg.Postgres.Database == "" {
		return fmt.Errorf("d%w: atabase name is not set", ErrConfig)
	}

	return nil
}

//nolint:nosprintfhostport
func (cfg *Config) PgURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
}
