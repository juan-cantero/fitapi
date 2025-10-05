package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(databaseURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	// Configure connection pool for production best practices
	config.MaxConns = 10                          // Don't flood the pooler (align with Supabase limits)
	config.MinConns = 2                           // Keep some warm connections
	config.MaxConnLifetime = 0                    // No max lifetime (let pooler handle it)
	config.MaxConnIdleTime = 30 * time.Minute     // 30 minutes idle timeout
	config.HealthCheckPeriod = 1 * time.Minute    // Health check every minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully")

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}
