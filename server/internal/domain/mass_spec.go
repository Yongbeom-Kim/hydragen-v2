package domain

type MassSpectraRecord struct {
	ID              int64     `json:"id"`
	InchiKey        string    `json:"inchiKey"`
	MolecularWeight float64   `json:"molecularWeight"`
	ExactMass       *float64  `json:"exactMass"`
	PrecursorMz     *float64  `json:"precursorMz"`
	PrecursorType   *string   `json:"precursorType"`
	IonMode         *string   `json:"ionMode"`
	CollisionEnergy *string   `json:"collisionEnergy"`
	SpectrumType    *string   `json:"spectrumType"`
	Instrument      *string   `json:"instrument"`
	InstrumentType  *string   `json:"instrumentType"`
	Splash          *string   `json:"splash"`
	DBNumber        string    `json:"dbNumber"`
	Source          string    `json:"source"`
	Comments        *string   `json:"comments"`
	MZ              []float32 `json:"mZ"`
	Peaks           []int     `json:"peaks"`
}

func Round(f float32) int {
	if f >= 0 {
		return int(f + 0.5)
	}
	return int(f - 0.5)
}

func MassSpectraFillInMissingMz(record MassSpectraRecord) MassSpectraRecord {
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

	return MassSpectraRecord{
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
