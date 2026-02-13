package third_party_chembl

import (
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

func fetchChemblImage(compound domain.CompoundMetadata) ([]byte, string, error) {
	resp, err := http.Get(GetImageUrl(compound.InchiKey))
	if err != nil {
		slog.Error("[fetchChemblImage]: http.Get failed", "inchiKey", compound.InchiKey, "error", err)
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("[fetchChemblImage]: Non-OK HTTP status code", "status", resp.StatusCode, "inchiKey", compound.InchiKey)
		return nil, "", io.EOF
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("[fetchChemblImage]: io.ReadAll failed", "inchiKey", compound.InchiKey, "error", err)
		return nil, "", err
	}

	contentType := normalizeMimeType(resp.Header.Get("Content-Type"))
	if contentType == "" {
		slog.Warn("[fetchChemblImage]: Empty Content-Type, defaulting", "inchiKey", compound.InchiKey)
		contentType = "image/svg+xml"
	}

	return data, contentType, nil
}

func GetChemblImage(compound domain.CompoundMetadata) ([]byte, string, error) {
	data, mimeType, err := disk.ReadChemblImageFromDisk(compound)
	if err == nil && data != nil && mimeType != "" {
		return data, mimeType, nil
	}

	data, mimeType, err = fetchChemblImage(compound)
	if err != nil {
		return nil, "", err
	}

	_ = disk.SaveChemblImageToDisk(data, mimeType, compound)

	return data, mimeType, nil
}
