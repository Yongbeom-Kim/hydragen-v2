package db

import (
	"context"
	"database/sql"
	"hydragen-v2/server/utils"
	"os"
	"time"
)

func Open() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		user := utils.EnvOrDefault("POSTGRES_USER", "postgres")
		password := utils.EnvOrDefault("POSTGRES_PASSWORD", "postgres")
		host := utils.EnvOrDefault("POSTGRES_HOST", "postgres")
		port := utils.EnvOrDefault("POSTGRES_PORT", "5432")
		database := utils.EnvOrDefault("POSTGRES_DB", "postgres")
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
