package postgres

import (
	"context"
	"database/sql"
	"hydragen-v2/server/internal/domain"
	"log/slog"
)

var fallbackSpectra = map[string][]domain.MassSpectraRecord{
	"XLYOFNOQVPJJNP-UHFFFAOYSA-N": {
		{ID: 1, InchiKey: "XLYOFNOQVPJJNP-UHFFFAOYSA-N", MolecularWeight: 32.0419, DBNumber: "MB0001", Source: "fallback", MZ: []float32{15, 29, 31, 32}, Peaks: []int{25, 40, 100, 12}},
	},
	"VNWKTOKETHGBQD-UHFFFAOYSA-N": {
		{ID: 2, InchiKey: "VNWKTOKETHGBQD-UHFFFAOYSA-N", MolecularWeight: 16.0425, DBNumber: "MB0002", Source: "fallback", MZ: []float32{12, 13, 14, 15, 16}, Peaks: []int{15, 35, 80, 60, 100}},
	},
}

const MZ_SCALE = 10000

type PostgresMassSpecStore struct {
	db          *sql.DB
	useFallback bool
}

func NewPostgresMassSpecStore(db *sql.DB, useFallback bool) *PostgresMassSpecStore {
	return &PostgresMassSpecStore{db: db, useFallback: useFallback}
}

func (s *PostgresMassSpecStore) GetSpectra(ctx context.Context, inchiKey string) ([]domain.MassSpectraRecord, error) {
	return GetMassSpectra(ctx, s.db, inchiKey, s.useFallback)
}

func scaleDownMzFromDb(rawValue []int) []float32 {
	var result []float32
	result = make([]float32, len(rawValue))
	for i, val := range rawValue {
		result[i] = float32(val) / float32(MZ_SCALE)
	}
	return result
}

func GetMassSpectra(ctx context.Context, db *sql.DB, inchiKey string, useFallback bool) ([]domain.MassSpectraRecord, error) {
	if useFallback {
		records, ok := fallbackSpectra[inchiKey]
		if !ok {
			return nil, sql.ErrNoRows
		}
		return records, nil

	}
	const spectraSQL = `
		SELECT
			id,
			inchikey,
			molecular_weight,
			exact_mass,
			precursor_mz,
			precursor_type,
			ion_mode,
			collision_energy,
			spectrum_type,
			instrument,
			instrument_type,
			splash,
			db_number,
			source,
			comments,
			m_z,
			peaks
		FROM mass_spectra
		WHERE inchikey = $1
		ORDER BY id ASC
	`

	rows, err := db.QueryContext(ctx, spectraSQL, inchiKey)
	if err != nil {
		slog.Error("GetMassSpectra: database query error", "error", err, "inchiKey", inchiKey)
		return nil, err
	}
	defer rows.Close()

	var spectra []domain.MassSpectraRecord
	for rows.Next() {
		var rec domain.MassSpectraRecord
		var MZ_raw PgInt4Array
		var Peaks_raw PgInt4Array
		err := rows.Scan(
			&rec.ID,
			&rec.InchiKey,
			&rec.MolecularWeight,
			&rec.ExactMass,
			&rec.PrecursorMz,
			&rec.PrecursorType,
			&rec.IonMode,
			&rec.CollisionEnergy,
			&rec.SpectrumType,
			&rec.Instrument,
			&rec.InstrumentType,
			&rec.Splash,
			&rec.DBNumber,
			&rec.Source,
			&rec.Comments,
			&MZ_raw,
			&Peaks_raw,
		)
		rec.MZ = scaleDownMzFromDb(MZ_raw)
		rec.Peaks = Peaks_raw
		if err != nil {
			slog.Error("GetMassSpectra: failed to scan row", "error", err, "inchiKey", inchiKey)
			return nil, err
		}
		spectra = append(spectra, rec)
	}
	if err := rows.Err(); err != nil {
		slog.Error("GetMassSpectra: rows iteration error", "error", err, "inchiKey", inchiKey)
		return nil, err
	}

	return spectra, nil
}
