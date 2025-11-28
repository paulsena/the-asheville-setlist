package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewDatabasePool creates a new PostgreSQL connection pool
func NewDatabasePool(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	// Configure pool
	config, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Pool settings
	config.MaxConns = 25                // Maximum number of connections in the pool
	config.MinConns = 5                 // Minimum number of connections to keep open
	config.MaxConnLifetime = time.Hour  // Close connections after 1 hour
	config.MaxConnIdleTime = 30 * time.Minute // Close idle connections after 30 minutes
	config.HealthCheckPeriod = time.Minute    // Check connection health every minute

	// Create pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection pool established")
	return pool, nil
}
