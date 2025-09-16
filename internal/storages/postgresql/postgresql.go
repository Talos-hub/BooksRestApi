package postgresql

import (
	"context"
	"fmt"

	"github.com/Talos-hub/BooksRestApi/internal/storages/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool   *pgxpool.Pool
	config *config.DatabaseConfig
}

// TODO
func NewPostgresStorage(config *config.DatabaseConfig) (*PostgresStorage, error) {
	pgxconf, err := pgxpool.ParseConfig(config.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("faild to parse connection string: %w", err)
	}
	//set connection pool settings
	pgxconf.MaxConns = int32(config.MaxConns)
	pgxconf.MinConns = int32(config.MinConns)
	pgxconf.MaxConnLifetime = config.ConnMaxLifeTime
	pgxconf.MaxConnIdleTime = config.ConnMaxIdleTime
	pgxconf.HealthCheckPeriod = config.HealthCheckPeriod

	// create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), pgxconf)
	if err != nil {
		return nil, fmt.Errorf("faild to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	// ping to database
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("faild to ping database: %w", err)
	}

	//TODO create table book if it doesn't exist

	return &PostgresStorage{
		pool:   pool,
		config: config,
	}, nil
}
