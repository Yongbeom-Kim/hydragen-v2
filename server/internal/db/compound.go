package db

import (
	"context"
	"database/sql"
	"hydragen-v2/server/internal/domain"
	"hydragen-v2/server/utils"
	"log/slog"
	"strings"
)

var fallbackCompounds []domain.CompoundMetadata = []domain.CompoundMetadata{
	{InchiKey: "XLYOFNOQVPJJNP-UHFFFAOYSA-N", Name: "Methanol", Inchi: "InChI=1S/CH4O/c1-2/h2H,1H3", Smiles: "CO", Formula: "CH4O", MolecularWeight: utils.Ptr(32.0419), HasMassSpectrum: true},
	{InchiKey: "VNWKTOKETHGBQD-UHFFFAOYSA-N", Name: "Methane", Inchi: "InChI=1S/CH4/h1H4", Smiles: "C", Formula: "CH4", MolecularWeight: utils.Ptr(16.0425), HasMassSpectrum: true},
}

func ListCompounds(ctx context.Context, db *sql.DB, page int, pageSize int, useFallback bool) ([]domain.CompoundMetadata, error) {
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	count := pageSize

	if useFallback {
		return listFallbackCompounds(offset, count)
	}
	return queryCompoundList(ctx, db, count, offset)
}

func listFallbackCompounds(count int, offset int) ([]domain.CompoundMetadata, error) {
	start := offset
	if start > len(fallbackCompounds) {
		start = len(fallbackCompounds)
	}
	end := start + count
	if end > len(fallbackCompounds) {
		end = len(fallbackCompounds)
	}
	return fallbackCompounds[start:end], nil
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

func CountCompounds(ctx context.Context, db *sql.DB, useFallback bool) (int, error) {
	if useFallback {
		return len(fallbackCompounds), nil
	}

	const countSQL = `SELECT COUNT(*) FROM compounds`
	var total int
	err := db.QueryRowContext(ctx, countSQL).Scan(&total)
	if err != nil {
		slog.Error("[DB CountCompounds]: db error", "error", err)
		return 0, err
	}
	return total, nil
}

func QueryCompoundDetail(ctx context.Context, db *sql.DB, inchiKey string, useFallback bool) (*domain.CompoundMetadata, error) {
	if useFallback {
		for _, item := range fallbackCompounds {
			if item.InchiKey == inchiKey {
				result := &domain.CompoundMetadata{
					InchiKey:        item.InchiKey,
					Name:            item.Name,
					Inchi:           item.Inchi,
					Smiles:          item.Smiles,
					Formula:         item.Formula,
					HasMassSpectrum: item.HasMassSpectrum,
				}
				return result, nil
			}
		}
		slog.Error("[DB QueryComopundDetail]: compound not found in fallback compounds", "inchiKey", inchiKey)
		return nil, sql.ErrNoRows
	}

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
	err := db.QueryRowContext(ctx, detailSQL, inchiKey).Scan(
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
