package disk

import (
	"fmt"
	"hydragen-v2/server/internal/domain"
	"log/slog"
	"os"
	"path/filepath"
)

func getChemblAssetDir(compound domain.CompoundMetadata) (string, error) {
	chemicalAssetDir, err := getChemicalAssetDir(compound)
	if err != nil {
		slog.Error("[getChemicalAssetDir()]: Error: Unable to determine base chemical asset directory for chembl asset", "error", err)
		return "", err
	}
	return filepath.Join(chemicalAssetDir, "chembl"), nil
}

func SaveChemblImageToDisk(image []byte, mimeType string, compound domain.CompoundMetadata) error {
	dir, err := getChemblAssetDir(compound)
	if err != nil {
		slog.Error("[SaveChemblImageToDisk]: unable to get chembl asset dir", "error", err)
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("[SaveChemblImageToDisk]: unable to create directory", "dir", dir, "error", err)
		return err
	}
	ext := MimeTypeToExtension(mimeType)
	filename := "image" + ext
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, image, 0644); err != nil {
		slog.Error("[SaveChemblImageToDisk]: unable to write file", "path", path, "error", err)
		return err
	}
	return nil
}

func ReadChemblImageFromDisk(compound domain.CompoundMetadata) ([]byte, string, error) {
	imageDir, err := getChemblAssetDir(compound)
	if err != nil {
		slog.Error("[getChemblAssetDir()]: Error: Unable to get chembl image asset directory", "error", err)
		return nil, "", err
	}

	files, err := filepath.Glob(filepath.Join(imageDir, "*"))
	if err != nil {
		slog.Error("[readChemblImage]: Error reading directory", "dir", imageDir, "error", err)
		return nil, "", err
	}

	if len(files) == 0 {
		err := fmt.Errorf("no files found in chembl image directory: %s", imageDir)
		slog.Error("[readChemblImage]: No files found in chembl image directory", "dir", imageDir, "error", err)
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
		slog.Warn("[readChemblImage]: Multiple files found in chembl image directory; using most recently edited", "dir", imageDir, "files", readable, "selected", selectedFile)
	}

	if !found {
		err := fmt.Errorf("[readChemblImage]: no readable files found in chembl image directory: %s", imageDir)
		slog.Error("[readChemblImage]: No readable files found in chembl image directory", "dir", imageDir, "error", err)
		return nil, "", err
	}

	data, err := os.ReadFile(selectedFile)
	if err != nil {
		slog.Error("[readChemblImage]: Could not read selected file", "file", selectedFile, "error", err)
		return nil, "", err
	}
	mimeType := ExtensionToMimeType(selectedFile)
	return data, mimeType, nil
}
