package third_party_chembl

import (
	"context"
	"database/sql"
	"errors"
	"hydragen-v2/server/internal/db"
	"hydragen-v2/server/internal/disk"
	"hydragen-v2/server/internal/domain"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func GetImageUrl(inchiKey string) string {
	return "https://www.ebi.ac.uk/chembl/api/data/image/" + inchiKey + "?format=svg"
}

func normalizeMimeType(contentType string) string {
	if contentType == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(contentType, ";")[0])
}

func fetchChemblImage(ctx context.Context, dbConn *sql.DB, compound domain.CompoundMetadata) ([]byte, string, error) {
	cooldownKey := compound.InchiKey
	onCooldown, nextAllowedAt, err := db.IsThirdPartyRequestOnCooldown(ctx, dbConn, "chembl-image", cooldownKey)
	if err != nil {
		return nil, "", err
	}
	if onCooldown {
		slog.Info("[fetchChemblImage]: skipped due to cooldown", "inchiKey", compound.InchiKey, "nextAllowedAt", nextAllowedAt)
		return nil, "", io.EOF
	}

	resp, err := http.Get(GetImageUrl(compound.InchiKey))
	if err != nil {
		_ = db.RegisterThirdPartyRequestFailure(ctx, dbConn, "chembl-image", cooldownKey)
		slog.Error("[fetchChemblImage]: http.Get failed", "inchiKey", compound.InchiKey, "error", err)
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_ = db.RegisterThirdPartyRequestFailure(ctx, dbConn, "chembl-image", cooldownKey)
		slog.Error("[fetchChemblImage]: Non-OK HTTP status code", "status", resp.StatusCode, "inchiKey", compound.InchiKey)
		return nil, "", io.EOF
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = db.RegisterThirdPartyRequestFailure(ctx, dbConn, "chembl-image", cooldownKey)
		slog.Error("[fetchChemblImage]: io.ReadAll failed", "inchiKey", compound.InchiKey, "error", err)
		return nil, "", err
	}
	if len(data) == 0 {
		_ = db.RegisterThirdPartyRequestFailure(ctx, dbConn, "chembl-image", cooldownKey)
		return nil, "", io.EOF
	}

	contentType := normalizeMimeType(resp.Header.Get("Content-Type"))
	if contentType == "" {
		slog.Warn("[fetchChemblImage]: Empty Content-Type, defaulting", "inchiKey", compound.InchiKey)
		contentType = "image/svg+xml"
	}

	if err := db.ClearThirdPartyRequestCooldown(ctx, dbConn, "chembl-image", cooldownKey); err != nil {
		slog.Warn("[fetchChemblImage]: failed to clear cooldown after success", "inchiKey", compound.InchiKey, "error", err)
	}

	return data, contentType, nil
}

func GetChemblImage(ctx context.Context, dbConn *sql.DB, compound domain.CompoundMetadata) ([]byte, string, error) {
	data, mimeType, err := disk.ReadChemblImageFromDisk(compound)
	if err == nil && data != nil && mimeType != "" {
		return data, mimeType, nil
	}

	data, mimeType, err = fetchChemblImage(ctx, dbConn, compound)
	if err != nil {
		return nil, "", err
	}
	if len(data) == 0 {
		return nil, "", io.EOF
	}

	if err := disk.SaveChemblImageToDisk(data, mimeType, compound); err != nil && !errors.Is(err, io.EOF) {
		slog.Warn("[GetChemblImage]: failed to save image to disk", "inchiKey", compound.InchiKey, "error", err)
	}

	return data, mimeType, nil
}
