package postgres

import (
	"context"
	"database/sql"
	"os"
	"time"
)

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func Open() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		user := envOrDefault("POSTGRES_USER", "postgres")
		password := envOrDefault("POSTGRES_PASSWORD", "postgres")
		host := envOrDefault("POSTGRES_HOST", "postgres")
		port := envOrDefault("POSTGRES_PORT", "5432")
		database := envOrDefault("POSTGRES_DB", "postgres")
		dsn = "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + database + "?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
