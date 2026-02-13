package handler

import (
	"hydragen-v2/server/utils"
	"log/slog"
	"net/http"
	"time"
)

type healthResponse struct {
	Status string `json:"status"`
}

func (_ *App) HealthHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("[HealthHandler]: start", "method", r.Method, "path", r.URL.Path)
	defer slog.Info("[HealthHandler]: end", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))

	slog.Info("Healthcheck request at /health")
	utils.WriteJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}
