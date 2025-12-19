package ingest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Handler handles HTTP requests for metric ingestion.
type Handler struct {
	service *Service
}

// NewHandler creates a new ingest handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// IngestSingle handles single metric ingestion.
//
//	@Summary		Ingest single metric
//	@Description	Ingest a single metric data point
//	@Tags			ingest
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		IngestRequest	true	"Metric data"
//	@Success		201		{object}	IngestResponse
//	@Failure		400		{object}	ErrorResponse	"Validation error"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		409		{object}	ErrorResponse	"Duplicate metric"
//	@Failure		500		{object}	ErrorResponse	"Internal error"
//	@Router			/ingest [post]
func (h *Handler) IngestSingle(w http.ResponseWriter, r *http.Request) {
	prod := ProductFromContext(r.Context())
	if prod == nil {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "missing or invalid API key",
		})
		return
	}

	var req IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "invalid request body",
		})
		return
	}

	response, err := h.service.IngestSingle(r.Context(), prod.ID, req)
	if err != nil {
		// Check for validation errors
		if ve, ok := IsValidationError(err); ok {
			respondJSON(w, http.StatusBadRequest, ErrorResponse{
				Error:   ve.errorType,
				Message: ve.message,
			})
			return
		}

		// Check for duplicate metric
		if errors.Is(err, ErrDuplicateMetric) {
			respondJSON(w, http.StatusConflict, ErrorResponse{
				Error:   "duplicate_metric",
				Message: "a metric with this name and timestamp already exists",
			})
			return
		}

		log.Printf("ingest single error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to ingest metric",
		})
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// IngestBatch handles batch metric ingestion.
//
//	@Summary		Ingest batch of metrics
//	@Description	Ingest multiple metric data points atomically (max 100)
//	@Tags			ingest
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		BatchIngestRequest	true	"Batch of metrics"
//	@Success		201		{object}	BatchIngestResponse
//	@Failure		400		{object}	ErrorResponse	"Validation error"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		409		{object}	ErrorResponse	"Duplicate metric"
//	@Failure		500		{object}	ErrorResponse	"Internal error"
//	@Router			/ingest/batch [post]
func (h *Handler) IngestBatch(w http.ResponseWriter, r *http.Request) {
	prod := ProductFromContext(r.Context())
	if prod == nil {
		respondJSON(w, http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "missing or invalid API key",
		})
		return
	}

	var req BatchIngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "invalid request body",
		})
		return
	}

	response, err := h.service.IngestBatch(r.Context(), prod.ID, req)
	if err != nil {
		// Check for validation errors
		if ve, ok := IsValidationError(err); ok {
			respondJSON(w, http.StatusBadRequest, ErrorResponse{
				Error:   ve.errorType,
				Message: ve.message,
			})
			return
		}

		// Check for duplicate metric
		if errors.Is(err, ErrDuplicateMetric) {
			respondJSON(w, http.StatusConflict, ErrorResponse{
				Error:   "duplicate_metric",
				Message: "a metric with this name and timestamp already exists",
			})
			return
		}

		log.Printf("ingest batch error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to ingest metrics",
		})
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
