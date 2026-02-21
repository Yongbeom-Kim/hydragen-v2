package disk

import (
	"log/slog"
	"os"
)

func getCwd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		slog.Error("[os.Getwd() fail]: Fatal error", "error", err)
		return "", err
	}
	return wd, nil
}
