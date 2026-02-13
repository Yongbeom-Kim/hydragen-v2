package handler

import (
	"hydragen-v2/server/utils"
	"log/slog"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func (_ *App) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	slog.Info("Healthcheck request at /health")
	utils.WriteJSON(w, http.StatusOK, healthResponse{Status: "ok"})
}
