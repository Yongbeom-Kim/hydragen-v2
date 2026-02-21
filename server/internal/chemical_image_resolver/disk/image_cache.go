package chemicalimageresolver_disk

import (
	"context"
	"fmt"
	chemicalimageresolver "hydragen-v2/server/internal/chemical_image_resolver/core"
	"hydragen-v2/server/internal/domain"
	"log/slog"
	"os"
	"path/filepath"
)

const assetDir = "assets"

func getChemicalAssetDir(compound domain.CompoundMetadata, providerType chemicalimageresolver.ProviderType) (string, error) {
	inchiKey := compound.InchiKey
	if len(inchiKey) < 4 {
		return "", fmt.Errorf("InchiKey too short: got length %d, want at least 4", len(inchiKey))
	}
	// Directory structure:
	//   {assetDir}/inchikey/{first2}/{next2}/{fullInchiKey}/{providerType}/
	path := filepath.Join(assetDir, "inchikey", inchiKey[0:2], inchiKey[2:4], inchiKey, string(providerType))
	return path, nil
}

type DiskImageCache struct{}

func (d *DiskImageCache) Fetch(ctx context.Context, provider chemicalimageresolver.ProviderType, c domain.CompoundMetadata) (img *chemicalimageresolver.Image, found bool, err error) {
	imageDir, err := getChemicalAssetDir(c, provider)
	if err != nil {
		slog.Error("[getChemicalAssetDir()]: Error: Unable to get image asset directory", "error", err)
		return nil, false, err
	}

	files, err := filepath.Glob(filepath.Join(imageDir, "*"))
	if err != nil {
		slog.Error("[DiskImageCache.Fetch]: Error reading directory", "dir", imageDir, "error", err)
		return nil, false, err
	}

	if len(files) == 0 {
		err := fmt.Errorf("no files found in image directory: %s", imageDir)
		slog.Error("[DiskImageCache.Fetch]: No files found in image directory", "dir", imageDir, "error", err)
		return nil, false, err
	}

	var (
		selectedFile string
		latestMod    int64
		fileFound    bool
		readable     []string
	)

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		f, err := os.Open(file)
		if err != nil {
			continue
		}
		f.Close()
		readable = append(readable, file)
		modTime := info.ModTime().UnixNano()
		if !fileFound || modTime > latestMod {
			selectedFile = file
			latestMod = modTime
			fileFound = true
		}
	}

	if len(readable) > 1 {
		slog.Warn(
			"Multiple files found in image cache directory; using most recently modified",
			"imageDir", imageDir,
			"files", readable,
			"selectedFile", selectedFile,
		)
	}

	if !fileFound {
		err := fmt.Errorf("no readable files found in image cache directory: %s", imageDir)
		slog.Error(
			"No readable files found in image cache directory",
			"imageDir", imageDir,
			"error", err,
			"files", files,
		)
		return nil, false, err
	}

	data, err := os.ReadFile(selectedFile)
	if err != nil {
		slog.Error(
			"Could not read selected cache file",
			"selectedFile", selectedFile,
			"error", err,
		)
		return nil, false, err
	}
	mimeType := ExtensionToMimeType(selectedFile)
	return &chemicalimageresolver.Image{Bytes: data, MimeType: chemicalimageresolver.MimeType(mimeType)}, true, nil
}

func (d *DiskImageCache) Save(ctx context.Context, provider chemicalimageresolver.ProviderType, c domain.CompoundMetadata, img *chemicalimageresolver.Image, imgMimeType string) error {
	dir, err := getChemicalAssetDir(c, provider)
	if err != nil {
		slog.Error("[DiskImageCache.Save]: unable to get chemical asset directory", "error", err, "provider", provider, "inchiKey", c.InchiKey)
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("[DiskImageCache.Save]: unable to create directory", "dir", dir, "error", err)
		return err
	}
	ext := MimeTypeToExtension(imgMimeType)
	filename := "image" + ext
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, img.Bytes, 0644); err != nil {
		slog.Error("[DiskImageCache.Save]: unable to write file", "path", path, "error", err)
		return err
	}
	return nil
}
