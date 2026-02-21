package chemicalimageresolver_http

import (
	"context"
	chemicalimageresolver "hydragen-v2/server/internal/chemical_image_resolver/core"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	resolver *chemicalimageresolver.Resolver
}

func NewHandler(resolver *chemicalimageresolver.Resolver) *Handler {
	return &Handler{resolver: resolver}
}

func (a *Handler) GetCompoundImageHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("[GetCompoundImageHandler]: start", "method", r.Method, "path", r.URL.Path)
	defer slog.Info("[GetCompoundImageHandler]: end", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	inchiKey := strings.TrimSpace(r.PathValue("inchiKey"))
	if inchiKey == "" {
		http.NotFound(w, r)
		return
	}

	image, ok := a.resolver.Image(ctx, inchiKey)
	if !ok || image == nil {
		http.NotFound(w, r)
		return
	}
	data := image.Bytes

	mimeType := image.MimeType

	if mimeType != "" {
		w.Header().Set("Content-Type", string(mimeType))
	}
	slog.Info(
		"[GetCompoundImageHandler]: image resolved",
		"inchiKey",
		inchiKey,
		"mimeType",
		mimeType,
		"sizeBytes",
		len(data),
	)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
