package handler

import (
	"context"
	"database/sql"
	"errors"
	"hydragen-v2/server/internal/db"
	cactus "hydragen-v2/server/third_party/cactus"
	chembl "hydragen-v2/server/third_party/chembl"
	"hydragen-v2/server/utils"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func (a *App) GetCompoundImageHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("[GetCompoundImageHandler]: start", "method", r.Method, "path", r.URL.Path)
	defer slog.Info("[GetCompoundImageHandler]: end", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))

	inchiKey := strings.TrimSpace(r.PathValue("inchiKey"))
	if inchiKey == "" {
		slog.Error("[GetCompoundImageHandler]: Empty inchiKey", "inchiKey", inchiKey)
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	compound, err := db.QueryCompoundDetail(ctx, a.Db, inchiKey, a.UseFallbackData())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		slog.Error("[GetCompoundImageHandler]: QueryCompoundDetail error", "inchiKey", inchiKey, "error", err)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	imageSource := "chembl"
	data, mimeType, err := chembl.GetChemblImage(ctx, a.Db, *compound)
	if err != nil || len(data) == 0 {
		if err != nil && !errors.Is(err, io.EOF) {
			slog.Warn("[GetCompoundImageHandler]: GetChemblImage error", "inchiKey", inchiKey, "error", err)
		}
		imageSource = "cactus"
		data, mimeType, err = cactus.GetCactusImage(ctx, a.Db, *compound)
	}
	if err != nil || len(data) == 0 {
		if err != nil && !errors.Is(err, io.EOF) {
			slog.Warn("[GetCompoundImageHandler]: GetCactusImage error", "inchiKey", inchiKey, "error", err)
		}
		http.NotFound(w, r)
		return
	}

	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}
	slog.Info(
		"[GetCompoundImageHandler]: image resolved",
		"inchiKey",
		inchiKey,
		"source",
		imageSource,
		"mimeType",
		mimeType,
		"sizeBytes",
		len(data),
	)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
