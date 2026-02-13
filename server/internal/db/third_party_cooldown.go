package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

const (
	initialCooldownHours = 1
	maxCooldownHours     = 24 * 7
)

func IsThirdPartyRequestOnCooldown(ctx context.Context, db *sql.DB, origin string, uniqueKey string) (bool, time.Time, error) {
	if db == nil {
		return false, time.Time{}, nil
	}

	const sqlQuery = `
		SELECT earliest_next_request
		FROM third_party_cooldown
		WHERE origin = $1 AND unique_key = $2
	`

	var nextAllowedAt time.Time
	err := db.QueryRowContext(ctx, sqlQuery, origin, uniqueKey).Scan(&nextAllowedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return false, time.Time{}, nil
	}
	if err != nil {
		slog.Error("[DB IsThirdPartyRequestOnCooldown]: postgres error", "error", err, "origin", origin, "uniqueKey", uniqueKey)
		return false, time.Time{}, err
	}

	if time.Now().Before(nextAllowedAt) {
		return true, nextAllowedAt, nil
	}
	return false, nextAllowedAt, nil
}

func RegisterThirdPartyRequestFailure(ctx context.Context, db *sql.DB, origin string, uniqueKey string) error {
	if db == nil {
		return nil
	}

	const sqlQuery = `
		WITH upsert AS (
			INSERT INTO third_party_cooldown (
				origin, unique_key, last_requested,
				current_cooldown_duration_hours, earliest_next_request
			)
			VALUES (
				$1, $2, NOW(),
				$3, NOW() + make_interval(hours => $3)
			)
			ON CONFLICT (origin, unique_key)
			DO UPDATE SET
				last_requested = NOW(),
				current_cooldown_duration_hours =
					LEAST($4, GREATEST($3, third_party_cooldown.current_cooldown_duration_hours * 2))
			RETURNING origin, unique_key, current_cooldown_duration_hours
		)
		UPDATE third_party_cooldown t
		SET earliest_next_request = NOW() + make_interval(hours => upsert.current_cooldown_duration_hours)
		FROM upsert
		WHERE t.origin = upsert.origin AND t.unique_key = upsert.unique_key;
	`

	_, err := db.ExecContext(ctx, sqlQuery, origin, uniqueKey, initialCooldownHours, maxCooldownHours)
	if err != nil {
		slog.Error("[DB RegisterThirdPartyRequestFailure]: postgres error", "error", err, "origin", origin, "uniqueKey", uniqueKey)
		return err
	}
	return nil
}

func ClearThirdPartyRequestCooldown(ctx context.Context, db *sql.DB, origin string, uniqueKey string) error {
	if db == nil {
		return nil
	}

	const sqlQuery = `
		DELETE FROM third_party_cooldown
		WHERE origin = $1 AND unique_key = $2
	`

	_, err := db.ExecContext(ctx, sqlQuery, origin, uniqueKey)
	if err != nil {
		slog.Error("[DB ClearThirdPartyRequestCooldown]: postgres error", "error", err, "origin", origin, "uniqueKey", uniqueKey)
		return err
	}

	return nil
}
