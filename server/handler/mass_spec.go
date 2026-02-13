package handler

import (
	"context"
	"hydragen-v2/server/internal/db"
	"hydragen-v2/server/internal/domain"
	"hydragen-v2/server/utils"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type massSpectrumResponse struct {
	InchiKey string                     `json:"inchiKey"`
	Count    int                        `json:"count"`
	Items    []domain.MassSpectraRecord `json:"items"`
}

func (a *App) GetMassSpectraHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("[GetMassSpectraHandler]: start", "method", r.Method, "path", r.URL.Path)
	defer slog.Info("[GetMassSpectraHandler]: end", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))

	inchiKey := strings.TrimSpace(r.PathValue("inchiKey"))
	if inchiKey == "" {
		slog.Error("[GetMassSpectraHandler]: Empty inchiKey", "inchiKey", inchiKey)
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	spectra, err := db.GetMassSpectra(ctx, a.Db, inchiKey, a.UseFallbackData())
	if err != nil {
		slog.Error("[GetMassSpectraHandler]: db.GetMassSpectra error", "inchiKey", inchiKey, "error", err)
		http.NotFound(w, r)
		return
	}
	if len(spectra) == 0 {
		slog.Error("[GetMassSpectraHandler]: No spectra found", "inchiKey", inchiKey)
		http.NotFound(w, r)
		return
	}
	spectra = utils.Map(spectra, domain.MassSpectraFillInMissingMz)

	slog.Info("[GetMassSpectraHandler]: successful response", "inchiKey", inchiKey, "spectraCount", len(spectra))
	utils.WriteJSON(w, http.StatusOK, massSpectrumResponse{
		InchiKey: spectra[0].InchiKey,
		Count:    len(spectra),
		Items:    spectra,
	})
}
