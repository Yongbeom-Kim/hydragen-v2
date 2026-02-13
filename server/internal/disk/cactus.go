package disk

import (
	"fmt"
	"hydragen-v2/server/internal/domain"
	"log/slog"
	"os"
	"path/filepath"
)

func getCactusAssetDir(compound domain.CompoundMetadata) (string, error) {
	chemicalAssetDir, err := getChemicalAssetDir(compound)
	if err != nil {
		slog.Error("[getChemicalAssetDir()]: Error: Unable to determine base chemical asset directory for cactus asset", "error", err)
		return "", err
	}
	return filepath.Join(chemicalAssetDir, "cactus"), nil
}

func SaveCactusImageToDisk(image []byte, mimeType string, compound domain.CompoundMetadata) error {
	dir, err := getCactusAssetDir(compound)
	if err != nil {
		slog.Error("[SaveCactusImageToDisk]: unable to get cactus asset dir", "error", err)
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("[SaveCactusImageToDisk]: unable to create directory", "dir", dir, "error", err)
		return err
	}
	ext := MimeTypeToExtension(mimeType)
	filename := "image" + ext
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, image, 0644); err != nil {
		slog.Error("[SaveCactusImageToDisk]: unable to write file", "path", path, "error", err)
		return err
	}
	return nil
}

func ReadCactusImageFromDisk(compound domain.CompoundMetadata) ([]byte, string, error) {
	imageDir, err := getCactusAssetDir(compound)
	if err != nil {
		slog.Error("[getCactusAssetDir()]: Error: Unable to get cactus image asset directory", "error", err)
		return nil, "", err
	}

	files, err := filepath.Glob(filepath.Join(imageDir, "*"))
	if err != nil {
		slog.Error("[readCactusImage]: Error reading directory", "dir", imageDir, "error", err)
		return nil, "", err
	}

	if len(files) == 0 {
		err := fmt.Errorf("no files found in cactus image directory: %s", imageDir)
		slog.Error("[readCactusImage]: No files found in cactus image directory", "dir", imageDir, "error", err)
		return nil, "", err
	}

	var (
		selectedFile string
		latestMod    int64
		found        bool
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
		if !found || modTime > latestMod {
			selectedFile = file
			latestMod = modTime
			found = true
		}
	}

	if len(readable) > 1 {
		slog.Warn("[readCactusImage]: Multiple files found in cactus image directory; using most recently edited", "dir", imageDir, "files", readable, "selected", selectedFile)
	}

	if !found {
		err := fmt.Errorf("[readCactusImage]: no readable files found in cactus image directory: %s", imageDir)
		slog.Error("[readCactusImage]: No readable files found in cactus image directory", "dir", imageDir, "error", err)
		return nil, "", err
	}

	data, err := os.ReadFile(selectedFile)
	if err != nil {
		slog.Error("[readCactusImage]: Could not read selected file", "file", selectedFile, "error", err)
		return nil, "", err
	}
	mimeType := ExtensionToMimeType(selectedFile)
	return data, mimeType, nil
}
