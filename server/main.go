package main

import (
	"hydragen-v2/server/handler"
	"hydragen-v2/server/internal/db"
	"hydragen-v2/server/utils"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type compoundDetailResponse struct {
	InchiKey         string `json:"inchiKey"`
	Name             string `json:"name"`
	Inchi            string `json:"inchi"`
	Smiles           string `json:"smiles"`
	Formula          string `json:"formula"`
	HasMassSpectrum  bool   `json:"hasMassSpectrum"`
	MassSpectrumHref string `json:"massSpectrumHref,omitempty"`
}

func main() {
	db, err := db.Open()
	if err != nil {
		slog.Error("Database unavailable. Using Fallback Data", "error", err)
	}

	application := &handler.App{Db: db}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", application.HealthHandler)
	mux.HandleFunc("GET /compounds", application.GetCompoundListHandler)
	mux.HandleFunc("GET /compounds/{inchiKey}", application.GetCompoundDetailHandler)
	mux.HandleFunc("GET /mass-spectra/{inchiKey}", application.GetMassSpectraHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: utils.WithCORS(mux),
	}

	log.Println("api server listening on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
