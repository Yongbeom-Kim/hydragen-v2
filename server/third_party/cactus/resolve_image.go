package third_party_cactus

import (
	"context"
	"database/sql"
	"hydragen-v2/server/internal/db"
	"hydragen-v2/server/internal/disk"
	"hydragen-v2/server/internal/domain"
	"hydragen-v2/server/utils"
	"io"
	"log/slog"
	"net/http"
)

func GetImageUrl(identifier string) string {
	return "https://cactus.nci.nih.gov/chemical/structure/" + identifier + "/image"
}

func GetImageUrlFallbackList(compound domain.CompoundMetadata) []string {
	var identifiers = []string{
		compound.InchiKey,
		compound.Smiles,
		compound.Name,
	}

	identifiers = utils.Filter(identifiers, func(id string) bool { return id != "" })
	return utils.Map(identifiers, GetImageUrl)
}

func fetchCactusImage(ctx context.Context, dbConn *sql.DB, compound domain.CompoundMetadata) ([]byte, string, error) {
	fallbackUrls := GetImageUrlFallbackList(compound)
	for _, url := range fallbackUrls {
		onCooldown, nextAllowedAt, err := db.IsThirdPartyRequestOnCooldown(ctx, dbConn, "cactus-image", url)
		if err != nil {
			return nil, "", err
		}
		if onCooldown {
			slog.Info("[fetchCactusImage]: skipped due to cooldown", "url", url, "nextAllowedAt", nextAllowedAt)
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			_ = db.RegisterThirdPartyRequestFailure(ctx, dbConn, "cactus-image", url)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			_ = db.RegisterThirdPartyRequestFailure(ctx, dbConn, "cactus-image", url)
			continue
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			_ = db.RegisterThirdPartyRequestFailure(ctx, dbConn, "cactus-image", url)
			continue
		}
		if len(data) == 0 {
			_ = db.RegisterThirdPartyRequestFailure(ctx, dbConn, "cactus-image", url)
			continue
		}
		contentType := resp.Header.Get("Content-Type")
		if err := db.ClearThirdPartyRequestCooldown(ctx, dbConn, "cactus-image", url); err != nil {
			slog.Warn("[fetchCactusImage]: failed to clear cooldown after success", "url", url, "error", err)
		}
		return data, contentType, nil
	}
	return nil, "", io.EOF // io.EOF: no image found with any identifier
}

func GetCactusImage(ctx context.Context, dbConn *sql.DB, compound domain.CompoundMetadata) ([]byte, string, error) {
	data, mimeType, err := disk.ReadCactusImageFromDisk(compound)
	if err == nil && data != nil && mimeType != "" {
		return data, mimeType, nil
	}

	data, mimeType, err = fetchCactusImage(ctx, dbConn, compound)
	if err != nil {
		return nil, "", err
	}

	_ = disk.SaveCactusImageToDisk(data, mimeType, compound)

	return data, mimeType, nil
}
