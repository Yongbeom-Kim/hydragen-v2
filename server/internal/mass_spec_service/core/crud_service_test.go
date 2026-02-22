package massspecservice

import (
	"context"
	"hydragen-v2/server/internal/domain"
	"testing"
)

// mockMassSpecStore simulates Postgres: returns records with MZ and Peaks directly.
type mockMassSpecStore struct {
	spectra map[string][]domain.MassSpectraRecord
}

func (m *mockMassSpecStore) GetSpectra(ctx context.Context, inchiKey string) ([]domain.MassSpectraRecord, error) {
	recs, ok := m.spectra[inchiKey]
	if !ok {
		return nil, nil
	}
	return recs, nil
}

func makeRecord(inchiKey string, mz []float32, peaks []int) domain.MassSpectraRecord {
	return domain.MassSpectraRecord{
		ID:              125454,
		InchiKey:        inchiKey,
		MolecularWeight: 42,
		DBNumber:        "MSBNK-Fac_Eng_Univ_Tokyo-JP001581",
		Source:          "MassBankDataLoader",
		MZ:              mz,
		Peaks:           peaks,
	}
}

// TestGetSpectra_MockedPostgres_PeaksTranslatedCorrectly verifies that when the store returns
// a spectrum with multiple peaks rounding to the same integer (e.g. 19.5 and 20.0 both round to 20),
// the fill-in logic emits each peak separately (no combining) and all peaks are translated correctly.
func TestGetSpectra_MockedPostgres_PeaksTranslatedCorrectly(t *testing.T) {
	mz := []float32{12, 13, 14, 15, 19, 19.5, 20, 25, 26, 27, 28, 36, 37, 38, 39, 40, 41, 42, 43}
	peaks := []int{11, 16, 35, 50, 29, 18, 20, 26, 116, 265, 22, 30, 158, 213, 734, 303, 999, 661, 23}

	inchiKey := "QQONPFPTGQHPMA-UHFFFAOYSA-N"
	record := makeRecord(inchiKey, mz, peaks)

	mockStore := &mockMassSpecStore{
		spectra: map[string][]domain.MassSpectraRecord{inchiKey: {record}},
	}
	svc := NewMassSpectraCrudService(mockStore)
	ctx := context.Background()

	got, err := svc.GetSpectra(ctx, inchiKey)
	if err != nil {
		t.Fatalf("GetSpectra: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 spectrum, got %d", len(got))
	}

	// Hardcoded expected result: one slot per integer 0..43, with two slots for 20 (19.5 and 20.0 not combined).
	wantMZ := []float32{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		19.5, 20, // two entries for peaks that round to 20
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43,
	}
	wantPeaks := []int{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 16, 35, 50, 0, 0, 0, 29,
		18, 20, // 19.5 -> 18, 20.0 -> 20 (not combined)
		0, 0, 0, 0, 26, 116, 265, 22, 0, 0, 0, 0, 0, 0, 0, 30, 158, 213, 734, 303, 999, 661, 23,
	}

	item := got[0]
	if len(item.MZ) != len(wantMZ) || len(item.Peaks) != len(wantPeaks) {
		t.Fatalf("length mismatch: got MZ=%d Peaks=%d, want MZ=%d Peaks=%d", len(item.MZ), len(item.Peaks), len(wantMZ), len(wantPeaks))
	}
	for i := range wantMZ {
		if item.MZ[i] != wantMZ[i] {
			t.Errorf("MZ[%d]: got %v want %v", i, item.MZ[i], wantMZ[i])
		}
		if item.Peaks[i] != wantPeaks[i] {
			t.Errorf("Peaks[%d]: got %d want %d", i, item.Peaks[i], wantPeaks[i])
		}
	}
}
