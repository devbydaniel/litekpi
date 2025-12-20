package ingest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/auth"
	"github.com/devbydaniel/litekpi/internal/product"
)

// Handler handles HTTP requests for measurement ingestion.
type Handler struct {
	service        *Service
	productService *product.Service
}

// NewHandler creates a new ingest handler.
func NewHandler(service *Service, productService *product.Service) *Handler {
	return &Handler{service: service, productService: productService}
}

// IngestSingle handles single measurement ingestion.
//
//	@Summary		Ingest single measurement
//	@Description	Ingest a single measurement data point
//	@Tags			ingest
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		IngestRequest	true	"Measurement data"
//	@Success		201		{object}	IngestResponse
//	@Failure		400		{object}	ErrorResponse	"Validation error"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		409		{object}	ErrorResponse	"Duplicate measurement"
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

		// Check for duplicate measurement
		if errors.Is(err, ErrDuplicateMeasurement) {
			respondJSON(w, http.StatusConflict, ErrorResponse{
				Error:   "duplicate_measurement",
				Message: "a measurement with this name and timestamp already exists",
			})
			return
		}

		log.Printf("ingest single error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to ingest measurement",
		})
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// IngestBatch handles batch measurement ingestion.
//
//	@Summary		Ingest batch of measurements
//	@Description	Ingest multiple measurement data points atomically (max 100)
//	@Tags			ingest
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			request	body		BatchIngestRequest	true	"Batch of measurements"
//	@Success		201		{object}	BatchIngestResponse
//	@Failure		400		{object}	ErrorResponse	"Validation error"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		409		{object}	ErrorResponse	"Duplicate measurement"
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

		// Check for duplicate measurement
		if errors.Is(err, ErrDuplicateMeasurement) {
			respondJSON(w, http.StatusConflict, ErrorResponse{
				Error:   "duplicate_measurement",
				Message: "a measurement with this name and timestamp already exists",
			})
			return
		}

		log.Printf("ingest batch error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to ingest measurements",
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

// validateProductOwnership validates that the product belongs to the user's organization.
func (h *Handler) validateProductOwnership(r *http.Request) (*product.Product, error) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	productIDStr := chi.URLParam(r, "productId")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		return nil, errors.New("invalid product ID")
	}

	prod, err := h.productService.GetProduct(r.Context(), user.OrganizationID, productID)
	if err != nil {
		if errors.Is(err, product.ErrProductNotFound) {
			return nil, errors.New("product not found")
		}
		if errors.Is(err, product.ErrUnauthorized) {
			return nil, errors.New("unauthorized")
		}
		return nil, err
	}

	return prod, nil
}

// ListMeasurementNames handles listing unique measurement names for a product.
//
//	@Summary		List measurement names
//	@Description	Get all unique measurement names for a product with their metadata keys
//	@Tags			measurements
//	@Produce		json
//	@Security		BearerAuth
//	@Param			productId	path		string	true	"Product ID"
//	@Success		200			{object}	ListMeasurementNamesResponse
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Product not found"
//	@Failure		500			{object}	ErrorResponse	"Internal error"
//	@Router			/products/{productId}/measurements [get]
func (h *Handler) ListMeasurementNames(w http.ResponseWriter, r *http.Request) {
	prod, err := h.validateProductOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "product not found" || err.Error() == "invalid product ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "product not found",
			})
			return
		}
		log.Printf("validate product ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate product",
		})
		return
	}

	measurements, err := h.service.GetMeasurementNames(r.Context(), prod.ID)
	if err != nil {
		log.Printf("get measurement names error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to get measurement names",
		})
		return
	}

	respondJSON(w, http.StatusOK, ListMeasurementNamesResponse{Measurements: measurements})
}

// GetMetadataValues handles getting metadata filter options for a measurement.
//
//	@Summary		Get metadata values
//	@Description	Get all unique metadata key-value options for filtering a measurement
//	@Tags			measurements
//	@Produce		json
//	@Security		BearerAuth
//	@Param			productId	path		string	true	"Product ID"
//	@Param			name		path		string	true	"Measurement name"
//	@Success		200			{object}	GetMetadataValuesResponse
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Product not found"
//	@Failure		500			{object}	ErrorResponse	"Internal error"
//	@Router			/products/{productId}/measurements/{name}/metadata [get]
func (h *Handler) GetMetadataValues(w http.ResponseWriter, r *http.Request) {
	prod, err := h.validateProductOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "product not found" || err.Error() == "invalid product ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "product not found",
			})
			return
		}
		log.Printf("validate product ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate product",
		})
		return
	}

	measurementName := chi.URLParam(r, "name")
	if measurementName == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "measurement name is required",
		})
		return
	}

	metadata, err := h.service.GetMetadataValues(r.Context(), prod.ID, measurementName)
	if err != nil {
		log.Printf("get metadata values error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to get metadata values",
		})
		return
	}

	respondJSON(w, http.StatusOK, GetMetadataValuesResponse{Metadata: metadata})
}

// GetMeasurementData handles getting aggregated chart data for a measurement.
//
//	@Summary		Get measurement data
//	@Description	Get daily aggregated data points for a measurement with optional metadata filtering and split-by
//	@Tags			measurements
//	@Produce		json
//	@Security		BearerAuth
//	@Param			productId	path		string	true	"Product ID"
//	@Param			name		path		string	true	"Measurement name"
//	@Param			start		query		string	true	"Start date (ISO 8601)"
//	@Param			end			query		string	true	"End date (ISO 8601)"
//	@Param			splitBy		query		string	false	"Metadata key to split data by"
//	@Success		200			{object}	GetMeasurementDataResponse
//	@Failure		400			{object}	ErrorResponse	"Validation error"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Product not found"
//	@Failure		500			{object}	ErrorResponse	"Internal error"
//	@Router			/products/{productId}/measurements/{name}/data [get]
func (h *Handler) GetMeasurementData(w http.ResponseWriter, r *http.Request) {
	prod, err := h.validateProductOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "product not found" || err.Error() == "invalid product ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "product not found",
			})
			return
		}
		log.Printf("validate product ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate product",
		})
		return
	}

	measurementName := chi.URLParam(r, "name")
	if measurementName == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "measurement name is required",
		})
		return
	}

	// Parse start date
	startStr := r.URL.Query().Get("start")
	if startStr == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "start date is required",
		})
		return
	}
	startDate, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		// Try parsing as date-only
		startDate, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, ErrorResponse{
				Error:   "validation_failed",
				Message: "invalid start date format (use ISO 8601)",
			})
			return
		}
	}

	// Parse end date
	endStr := r.URL.Query().Get("end")
	if endStr == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "end date is required",
		})
		return
	}
	endDate, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		// Try parsing as date-only
		endDate, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, ErrorResponse{
				Error:   "validation_failed",
				Message: "invalid end date format (use ISO 8601)",
			})
			return
		}
		// For date-only, add one day to make it inclusive
		endDate = endDate.AddDate(0, 0, 1)
	}

	// Parse metadata filters (metadata.key=value format)
	metadataFilters := make(map[string]string)
	for key, values := range r.URL.Query() {
		if strings.HasPrefix(key, "metadata.") && len(values) > 0 {
			metadataKey := strings.TrimPrefix(key, "metadata.")
			metadataFilters[metadataKey] = values[0]
		}
	}

	// Check for splitBy parameter
	splitByKey := r.URL.Query().Get("splitBy")
	if splitByKey != "" {
		series, err := h.service.GetAggregatedMeasurementsSplitBy(r.Context(), prod.ID, measurementName, startDate, endDate, metadataFilters, splitByKey)
		if err != nil {
			log.Printf("get measurement data split by error: %v", err)
			respondJSON(w, http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: "failed to get measurement data",
			})
			return
		}

		respondJSON(w, http.StatusOK, GetMeasurementDataSplitResponse{
			Name:    measurementName,
			SplitBy: splitByKey,
			Series:  series,
		})
		return
	}

	dataPoints, err := h.service.GetAggregatedMeasurements(r.Context(), prod.ID, measurementName, startDate, endDate, metadataFilters)
	if err != nil {
		log.Printf("get measurement data error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to get measurement data",
		})
		return
	}

	respondJSON(w, http.StatusOK, GetMeasurementDataResponse{
		Name:       measurementName,
		DataPoints: dataPoints,
	})
}

// GetPreferences handles getting saved chart preferences for a measurement.
//
//	@Summary		Get measurement preferences
//	@Description	Get saved chart preferences for a measurement
//	@Tags			measurements
//	@Produce		json
//	@Security		BearerAuth
//	@Param			productId	path		string	true	"Product ID"
//	@Param			name		path		string	true	"Measurement name"
//	@Success		200			{object}	GetPreferencesResponse
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Product not found"
//	@Failure		500			{object}	ErrorResponse	"Internal error"
//	@Router			/products/{productId}/measurements/{name}/preferences [get]
func (h *Handler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	prod, err := h.validateProductOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "product not found" || err.Error() == "invalid product ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "product not found",
			})
			return
		}
		log.Printf("validate product ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate product",
		})
		return
	}

	measurementName := chi.URLParam(r, "name")
	if measurementName == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "measurement name is required",
		})
		return
	}

	prefs, err := h.service.GetPreferences(r.Context(), prod.ID, measurementName)
	if err != nil {
		log.Printf("get preferences error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to get preferences",
		})
		return
	}

	respondJSON(w, http.StatusOK, GetPreferencesResponse{Preferences: prefs})
}

// SavePreferences handles saving chart preferences for a measurement.
//
//	@Summary		Save measurement preferences
//	@Description	Save chart preferences for a measurement (creates or updates)
//	@Tags			measurements
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			productId	path		string					true	"Product ID"
//	@Param			name		path		string					true	"Measurement name"
//	@Param			request		body		SavePreferencesRequest	true	"Preferences data"
//	@Success		200			{object}	map[string]string		"success message"
//	@Failure		400			{object}	ErrorResponse			"Validation error"
//	@Failure		401			{object}	ErrorResponse			"Unauthorized"
//	@Failure		404			{object}	ErrorResponse			"Product not found"
//	@Failure		500			{object}	ErrorResponse			"Internal error"
//	@Router			/products/{productId}/measurements/{name}/preferences [post]
func (h *Handler) SavePreferences(w http.ResponseWriter, r *http.Request) {
	prod, err := h.validateProductOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "product not found" || err.Error() == "invalid product ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "product not found",
			})
			return
		}
		log.Printf("validate product ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate product",
		})
		return
	}

	measurementName := chi.URLParam(r, "name")
	if measurementName == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "measurement name is required",
		})
		return
	}

	var req SavePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "invalid request body",
		})
		return
	}

	err = h.service.SavePreferences(r.Context(), prod.ID, measurementName, &req.Preferences)
	if err != nil {
		if ve, ok := IsValidationError(err); ok {
			respondJSON(w, http.StatusBadRequest, ErrorResponse{
				Error:   ve.errorType,
				Message: ve.message,
			})
			return
		}
		log.Printf("save preferences error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to save preferences",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "preferences saved"})
}
