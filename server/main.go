package main

import (
	chemicalimageresolver "hydragen-v2/server/internal/chemical_image_resolver/core"
	chemicalimageresolver_disk "hydragen-v2/server/internal/chemical_image_resolver/disk"
	chemicalimageresolver_http "hydragen-v2/server/internal/chemical_image_resolver/http"
	chemicalimageresolver_postgres "hydragen-v2/server/internal/chemical_image_resolver/postgres"
	chemicalimageresolver_thirdparty "hydragen-v2/server/internal/chemical_image_resolver/third_party"
	compoundmetadatastore "hydragen-v2/server/internal/compound_metadata_store/core"
	compoundmetadatastore_http "hydragen-v2/server/internal/compound_metadata_store/http"
	"hydragen-v2/server/internal/http_helper"
	massspecservice "hydragen-v2/server/internal/mass_spec_service/core"
	massspecservice_http "hydragen-v2/server/internal/mass_spec_service/http"
	"hydragen-v2/server/internal/postgres"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	db, err := postgres.Open()
	if err != nil {
		slog.Error("Database unavailable. Using Fallback Data", "error", err)
	}

	compoundStore := postgres.NewPostgresCompoundMetadataStore(db)
	compoundService := compoundmetadatastore.NewService(compoundStore)
	compoundHandler := compoundmetadatastore_http.NewHandler(compoundService)

	massSpecStore := postgres.NewPostgresMassSpecStore(db, db == nil)
	massSpecService := massspecservice.NewMassSpectraCrudService(massSpecStore)
	massSpecHandler := massspecservice_http.NewHandler(massSpecService)

	providers := map[chemicalimageresolver.ProviderType]chemicalimageresolver.ThirdPartyProvider{
		chemicalimageresolver.ProviderType("chembl"): &chemicalimageresolver_thirdparty.ChemblThirdPartyProvider{},
		chemicalimageresolver.ProviderType("cactus"): &chemicalimageresolver_thirdparty.CactusThirdPartyProvider{},
	}
	providerOrder := []chemicalimageresolver.ProviderType{
		chemicalimageresolver.ProviderType("chembl"),
		chemicalimageresolver.ProviderType("cactus"),
	}
	imageResolver := chemicalimageresolver.New(
		chemicalimageresolver_postgres.NewPostgresRequestCooldownStore(db),
		compoundStore,
		&chemicalimageresolver_disk.DiskImageCache{},
		providers,
		providerOrder,
	)
	imageHandler := chemicalimageresolver_http.NewHandler(imageResolver)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		http_helper.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /compounds", compoundHandler.GetCompoundListHandler)
	mux.HandleFunc("GET /compounds/{inchiKey}", compoundHandler.GetCompoundDetailHandler)
	mux.HandleFunc("GET /compounds/{inchiKey}/image", imageHandler.GetCompoundImageHandler)
	mux.HandleFunc("GET /mass-spectra/{inchiKey}", massSpecHandler.GetMassSpectraHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: http_helper.WithCORS(mux),
	}

	log.Println("api server listening on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
