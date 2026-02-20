package compoundmetadatastore_http

import (
	"context"
	"database/sql"
	"errors"
	compoundmetadatastore "hydragen-v2/server/internal/compound_metadata_store/core"
	"hydragen-v2/server/internal/domain"
	"hydragen-v2/server/internal/http_helper"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	service compoundmetadatastore.Service
}

func NewHandler(service *compoundmetadatastore.Service) *Handler {
	return &Handler{service: *service}
}

type compoundsListResponse struct {
	Items    []domain.CompoundMetadata `json:"items"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"pageSize"`
	Total    int                       `json:"total"`
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func (h *Handler) GetCompoundListHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("[GetCompoundListHandler]: start", "method", r.Method, "path", r.URL.Path, "rawQuery", r.URL.RawQuery)
	defer slog.Info("[GetCompoundListHandler]: end", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))

	query := r.URL.Query()
	page := parsePositiveInt(query.Get("page"), 1)
	pageSize := parsePositiveInt(query.Get("pageSize"), 20)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	compounds, err := h.service.List(ctx, page, pageSize)
	if err != nil {
		slog.Error("[GetCompoundListHandler]: service.List error", "error", err, "page", page, "pageSize", pageSize)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	count, err := h.service.Count(ctx)
	if err != nil {
		slog.Error("[GetCompoundListHandler]: db.ListCompounds error", "error", err, "page", page, "pageSize", pageSize)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	slog.Info("[GetCompoundListHandler]: successful response", "page", page, "pageSize", pageSize)

	http_helper.WriteJSON(w, http.StatusOK, compoundsListResponse{Items: compounds, Page: page, PageSize: pageSize, Total: count})
}

func (h *Handler) GetCompoundDetailHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	slog.Info("[GetCompoundDetailHandler]: start", "method", r.Method, "path", r.URL.Path)
	defer slog.Info("[GetCompoundDetailHandler]: end", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))

	inchiKey := r.PathValue("inchiKey")
	if inchiKey == "" {
		slog.Error("[GetCompoundDetailHandler]: Empty inchiKey", "inchiKey", inchiKey)
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	record, err := h.service.Get(ctx, inchiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("[GetCompoundDetailHandler]: Compound not found", "inchiKey", inchiKey, "error", err)
			http.NotFound(w, r)
			return
		}
		slog.Error("[GetCompoundDetailHandler]: QueryCompoundDetail error", "inchiKey", inchiKey, "error", err)
		http_helper.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	slog.Info("[GetCompoundDetailHandler]: successful response", "inchiKey", inchiKey)
	http_helper.WriteJSON(w, http.StatusOK, record)
}
