package disk

import (
	"fmt"
	"hydragen-v2/server/internal/domain"
	"log/slog"
	"os"
	"path/filepath"
)

func getCwd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		slog.Error("[os.Getwd() fail]: Fatal error", "error", err)
		return "", err
	}
	return wd, nil
}

const assetDir = "assets"

func getChemicalAssetDir(compound domain.CompoundMetadata) (string, error) {
	inchiKey := compound.InchiKey
	if len(inchiKey) < 4 {
		return "", fmt.Errorf("InchiKey too short: got length %d, want at least 4", len(inchiKey))
	}
	// Directory structure:
	//   {assetDir}/inchikey/{first2}/{next2}/{fullInchiKey}/
	path := filepath.Join(assetDir, "inchikey", inchiKey[0:2], inchiKey[2:4], inchiKey) + string(os.PathSeparator)
	return path, nil
}
