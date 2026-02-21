package chemicalimagersolver_postgres

import (
	"context"
	"database/sql"
	"errors"
	chemicalimageresolver "hydragen-v2/server/internal/chemical_image_resolver/core"
	"hydragen-v2/server/internal/domain"
	"log/slog"
	"time"
)

func NewPostgresRequestCooldownStore(db *sql.DB) *PostgresRequestCooldownStore {
	return &PostgresRequestCooldownStore{
		db: db,
	}
}

type PostgresRequestCooldownStore struct {
	db *sql.DB
}

func (p *PostgresRequestCooldownStore) OnCooldown(ctx context.Context, provider chemicalimageresolver.ProviderType, c domain.CompoundMetadata) (bool, error) {
	if p.db == nil {
		return false, nil
	}

	const sqlQuery = `
		SELECT earliest_next_request
		FROM third_party_cooldown
		WHERE origin = $1 AND unique_key = $2
	`

	var nextAllowedAt time.Time
	err := p.db.QueryRowContext(ctx, sqlQuery, provider, c.InchiKey).Scan(&nextAllowedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		slog.Error("[DB IsThirdPartyRequestOnCooldown]: postgres error", "error", err, "origin", provider, "uniqueKey", c.InchiKey)
		return false, err
	}

	if time.Now().Before(nextAllowedAt) {
		return true, nil
	}
	return false, nil
}

const (
	initialCooldownHours = 1
	maxCooldownHours     = 24 * 7
)

func (p *PostgresRequestCooldownStore) Add(ctx context.Context, provider chemicalimageresolver.ProviderType, c domain.CompoundMetadata) error {
	if p.db == nil {
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

	_, err := p.db.ExecContext(ctx, sqlQuery, provider, c.InchiKey, initialCooldownHours, maxCooldownHours)
	if err != nil {
		slog.Error("[DB RegisterThirdPartyRequestFailure]: postgres error", "error", err, "origin", provider, "uniqueKey", c.InchiKey)
		return err
	}
	return nil
}

func (p *PostgresRequestCooldownStore) Remove(ctx context.Context, provider chemicalimageresolver.ProviderType, c domain.CompoundMetadata) error {
	if p.db == nil {
		return nil
	}

	const sqlQuery = `
		DELETE FROM third_party_cooldown
		WHERE origin = $1 AND unique_key = $2
	`

	_, err := p.db.ExecContext(ctx, sqlQuery, provider, c.InchiKey)
	if err != nil {
		slog.Error("[DB ClearThirdPartyRequestCooldown]: postgres error", "error", err, "origin", provider, "uniqueKey", c.InchiKey)
		return err
	}

	return nil
}
