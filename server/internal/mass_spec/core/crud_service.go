package massspec

import (
	"context"
	"hydragen-v2/server/internal/domain"
	"log/slog"
)

type MassSpecStore interface {
	GetSpectra(ctx context.Context, inchiKey string) ([]domain.MassSpectraRecord, error)
}

type MassSpecCrudService interface {
	GetSpectra(ctx context.Context, inchiKey string) ([]domain.MassSpectraRecord, error)
}

type MassSpecCrudHandler struct {
	store MassSpecStore
}

func (h *MassSpecCrudHandler) GetSpectra(ctx context.Context, inchiKey string) ([]domain.MassSpectraRecord, error) {
	spectra, err := h.store.GetSpectra(ctx, inchiKey)
	if err != nil {
		slog.Error("[MassSpecCrudHandler.GetSpectra] failed to retrieve spectra", "inchiKey", inchiKey, "error", err)
		return nil, err
	}

	processedSpectra := make([]domain.MassSpectraRecord, len(spectra))
	for i := 0; i < len(spectra); i++ {
		processedSpectra[i] = MassSpectraFillInMissingMz(spectra[i])
	}
	return processedSpectra, nil
}

func Round(f float32) int {
	if f >= 0 {
		return int(f + 0.5)
	}
	return int(f - 0.5)
}

func MassSpectraFillInMissingMz(record domain.MassSpectraRecord) domain.MassSpectraRecord {
	maxMz := float32(0)
	for _, mz := range record.MZ {
		if mz > maxMz {
			maxMz = mz
		}
	}

	maxMzRounded := Round(maxMz)

	var resultMz []float32
	var resultPeaks []int
	mz_idx := 0

	for mz := 0; mz <= maxMzRounded; mz++ {
		if mz_idx < len(record.MZ) && mz_idx < len(record.Peaks) && Round(record.MZ[mz_idx]) == mz {
			resultMz = append(resultMz, record.MZ[mz_idx])
			resultPeaks = append(resultPeaks, record.Peaks[mz_idx])
			mz_idx++
		} else {
			resultMz = append(resultMz, float32(mz))
			resultPeaks = append(resultPeaks, 0)
		}
	}

	return domain.MassSpectraRecord{
		InchiKey:        record.InchiKey,
		MolecularWeight: record.MolecularWeight,
		ExactMass:       record.ExactMass,
		PrecursorMz:     record.PrecursorMz,
		PrecursorType:   record.PrecursorType,
		IonMode:         record.IonMode,
		CollisionEnergy: record.CollisionEnergy,
		SpectrumType:    record.SpectrumType,
		Instrument:      record.Instrument,
		InstrumentType:  record.InstrumentType,
		Splash:          record.Splash,
		DBNumber:        record.DBNumber,
		Source:          record.Source,
		Comments:        record.Comments,
		MZ:              resultMz,
		Peaks:           resultPeaks,
	}
}
