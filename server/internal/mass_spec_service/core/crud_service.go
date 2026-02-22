package massspecservice

import (
	"context"
	"hydragen-v2/server/internal/domain"
	"log/slog"
)

type MassSpecStore interface {
	GetSpectra(ctx context.Context, inchiKey string) ([]domain.MassSpectraRecord, error)
}

type Service interface {
	GetSpectra(ctx context.Context, inchiKey string) ([]domain.MassSpectraRecord, error)
}

type MassSpectraCrudService struct {
	store MassSpecStore
}

func NewMassSpectraCrudService(store MassSpecStore) *MassSpectraCrudService {
	return &MassSpectraCrudService{store: store}
}

func (h *MassSpectraCrudService) GetSpectra(ctx context.Context, inchiKey string) ([]domain.MassSpectraRecord, error) {
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
		// Consume all peaks that round to this integer (e.g. 19.5 and 20.0 both round to 20).
		// Emit one slot per peak so we don't combine; otherwise we'd never advance mz_idx and drop later peaks.
		emitted := false
		for mz_idx < len(record.MZ) && Round(record.MZ[mz_idx]) == mz {
			resultMz = append(resultMz, record.MZ[mz_idx])
			if mz_idx < len(record.Peaks) {
				resultPeaks = append(resultPeaks, record.Peaks[mz_idx])
			} else {
				resultPeaks = append(resultPeaks, 0)
			}
			mz_idx++
			emitted = true
		}
		if !emitted {
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
