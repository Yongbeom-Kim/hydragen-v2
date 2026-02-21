package db

import (
	"context"
	"database/sql"
	"hydragen-v2/server/internal/domain"
	"log/slog"
	"strings"
)

type PostgresCompoundMetadataStore struct {
	db *sql.DB
}

func (store *PostgresCompoundMetadataStore) Get(ctx context.Context, inchiKey string) (*domain.CompoundMetadata, error) {
	const detailSQL = `
		SELECT
			c.inchikey,
			c.name,
			c.inchi,
			c.smiles,
			c.formula,
			COUNT(ms.id) > 0 AS has_mass_spectrum
		FROM compounds c
		LEFT JOIN mass_spectra ms ON ms.inchikey = c.inchikey
		WHERE c.inchikey = $1
		GROUP BY c.inchikey, c.name, c.inchi, c.smiles, c.formula
	`

	var result domain.CompoundMetadata
	err := store.db.QueryRowContext(ctx, detailSQL, inchiKey).Scan(
		&result.InchiKey,
		&result.Name,
		&result.Inchi,
		&result.Smiles,
		&result.Formula,
		&result.HasMassSpectrum,
	)
	if err != nil {
		slog.Error("[DB QueryComopundDetail]: postgres error", "error", err, "inchiKey", inchiKey)
		return nil, err
	}
	result.InchiKey = strings.TrimSpace(result.InchiKey)
	return &result, nil
}

func ListCompounds(ctx context.Context, db *sql.DB, page int, pageSize int) ([]domain.CompoundMetadata, error) {
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	count := pageSize

	return queryCompoundList(ctx, db, count, offset)
}

func queryCompoundList(ctx context.Context, db *sql.DB, count int, offset int) ([]domain.CompoundMetadata, error) {
	const listCompoundsSQL = `
		SELECT
			c.inchikey,
			c.name,
			c.inchi,
			c.smiles,
			c.formula,
			MIN(ms.molecular_weight) AS molecular_weight,
			COUNT(ms.id) > 0 AS has_mass_spectrum
		FROM compounds c
		LEFT JOIN mass_spectra ms ON ms.inchikey = c.inchikey
		GROUP BY c.inchikey, c.name, c.inchi, c.smiles, c.formula
		ORDER BY MIN(ms.molecular_weight) ASC NULLS LAST, c.name ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.QueryContext(ctx, listCompoundsSQL, count, offset)
	if err != nil {
		slog.Error("[DB QueryCompoundList]: postgres error", "error", err, "count", count, "offset", offset)
		return nil, err
	}
	defer rows.Close()

	var results []domain.CompoundMetadata

	for rows.Next() {
		var item domain.CompoundMetadata
		err := rows.Scan(
			&item.InchiKey,
			&item.Name,
			&item.Inchi,
			&item.Smiles,
			&item.Formula,
			&item.MolecularWeight,
			&item.HasMassSpectrum,
		)
		if err != nil {
			slog.Error("[DB QueryCompoundList]: failed to scan compound list row", "error", err, "count", count, "offset", offset)
			return nil, err
		}
		results = append(results, item)
	}
	if err := rows.Err(); err != nil {
		slog.Error("[DB QueryCompoundList]: rows iteration error in compound list", "error", err, "count", count, "offset", offset)
		return nil, err
	}
	return results, nil

}

func CountCompounds(ctx context.Context, db *sql.DB) (int, error) {
	const countSQL = `SELECT COUNT(*) FROM compounds`
	var total int
	err := db.QueryRowContext(ctx, countSQL).Scan(&total)
	if err != nil {
		slog.Error("[DB CountCompounds]: db error", "error", err)
		return 0, err
	}
	return total, nil
}
