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
	"github.com/devbydaniel/litekpi/internal/datasource"
)

// Handler handles HTTP requests for measurement ingestion.
type Handler struct {
	service          *Service
	dataSourceService *datasource.Service
}

// NewHandler creates a new ingest handler.
func NewHandler(service *Service, dataSourceService *datasource.Service) *Handler {
	return &Handler{service: service, dataSourceService: dataSourceService}
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
	ds := DataSourceFromContext(r.Context())
	if ds == nil {
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

	response, err := h.service.IngestSingle(r.Context(), ds.ID, req)
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
	ds := DataSourceFromContext(r.Context())
	if ds == nil {
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

	response, err := h.service.IngestBatch(r.Context(), ds.ID, req)
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

// validateDataSourceOwnership validates that the data source belongs to the user's organization.
func (h *Handler) validateDataSourceOwnership(r *http.Request) (*datasource.DataSource, error) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	dataSourceIDStr := chi.URLParam(r, "dataSourceId")
	dataSourceID, err := uuid.Parse(dataSourceIDStr)
	if err != nil {
		return nil, errors.New("invalid data source ID")
	}

	ds, err := h.dataSourceService.GetDataSource(r.Context(), user.OrganizationID, dataSourceID)
	if err != nil {
		if errors.Is(err, datasource.ErrDataSourceNotFound) {
			return nil, errors.New("data source not found")
		}
		if errors.Is(err, datasource.ErrUnauthorized) {
			return nil, errors.New("unauthorized")
		}
		return nil, err
	}

	return ds, nil
}

// ListMeasurementNames handles listing unique measurement names for a data source.
//
//	@Summary		List measurement names
//	@Description	Get all unique measurement names for a data source with their metadata keys
//	@Tags			measurements
//	@Produce		json
//	@Security		BearerAuth
//	@Param			dataSourceId	path		string	true	"Data Source ID"
//	@Success		200				{object}	ListMeasurementNamesResponse
//	@Failure		401				{object}	ErrorResponse	"Unauthorized"
//	@Failure		404				{object}	ErrorResponse	"Data source not found"
//	@Failure		500				{object}	ErrorResponse	"Internal error"
//	@Router			/data-sources/{dataSourceId}/measurements [get]
func (h *Handler) ListMeasurementNames(w http.ResponseWriter, r *http.Request) {
	ds, err := h.validateDataSourceOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "data source not found" || err.Error() == "invalid data source ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "data source not found",
			})
			return
		}
		log.Printf("validate data source ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate data source",
		})
		return
	}

	measurements, err := h.service.GetMeasurementNames(r.Context(), ds.ID)
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
//	@Param			dataSourceId	path		string	true	"Data Source ID"
//	@Param			name			path		string	true	"Measurement name"
//	@Success		200				{object}	GetMetadataValuesResponse
//	@Failure		401				{object}	ErrorResponse	"Unauthorized"
//	@Failure		404				{object}	ErrorResponse	"Data source not found"
//	@Failure		500				{object}	ErrorResponse	"Internal error"
//	@Router			/data-sources/{dataSourceId}/measurements/{name}/metadata [get]
func (h *Handler) GetMetadataValues(w http.ResponseWriter, r *http.Request) {
	ds, err := h.validateDataSourceOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "data source not found" || err.Error() == "invalid data source ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "data source not found",
			})
			return
		}
		log.Printf("validate data source ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate data source",
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

	metadata, err := h.service.GetMetadataValues(r.Context(), ds.ID, measurementName)
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
//	@Description	Get daily aggregated data points for a measurement with optional metadata filtering
//	@Tags			measurements
//	@Produce		json
//	@Security		BearerAuth
//	@Param			dataSourceId	path		string	true	"Data Source ID"
//	@Param			name			path		string	true	"Measurement name"
//	@Param			start			query		string	true	"Start date (ISO 8601)"
//	@Param			end				query		string	true	"End date (ISO 8601)"
//	@Success		200				{object}	GetMeasurementDataResponse
//	@Failure		400				{object}	ErrorResponse	"Validation error"
//	@Failure		401				{object}	ErrorResponse	"Unauthorized"
//	@Failure		404				{object}	ErrorResponse	"Data source not found"
//	@Failure		500				{object}	ErrorResponse	"Internal error"
//	@Router			/data-sources/{dataSourceId}/measurements/{name}/data [get]
func (h *Handler) GetMeasurementData(w http.ResponseWriter, r *http.Request) {
	ds, err := h.validateDataSourceOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "data source not found" || err.Error() == "invalid data source ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "data source not found",
			})
			return
		}
		log.Printf("validate data source ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate data source",
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

	startDate, endDate, err := h.parseDateRange(r)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: err.Error(),
		})
		return
	}

	metadataFilters := h.parseMetadataFilters(r)

	dataPoints, err := h.service.GetAggregatedMeasurements(r.Context(), ds.ID, measurementName, startDate, endDate, metadataFilters)
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

// GetMeasurementDataSplit handles getting aggregated chart data split by a metadata key.
//
//	@Summary		Get measurement data split by metadata
//	@Description	Get daily aggregated data points for a measurement split by a metadata key
//	@Tags			measurements
//	@Produce		json
//	@Security		BearerAuth
//	@Param			dataSourceId	path		string	true	"Data Source ID"
//	@Param			name			path		string	true	"Measurement name"
//	@Param			start			query		string	true	"Start date (ISO 8601)"
//	@Param			end				query		string	true	"End date (ISO 8601)"
//	@Param			splitBy			query		string	true	"Metadata key to split data by"
//	@Success		200				{object}	GetMeasurementDataSplitResponse
//	@Failure		400				{object}	ErrorResponse	"Validation error"
//	@Failure		401				{object}	ErrorResponse	"Unauthorized"
//	@Failure		404				{object}	ErrorResponse	"Data source not found"
//	@Failure		500				{object}	ErrorResponse	"Internal error"
//	@Router			/data-sources/{dataSourceId}/measurements/{name}/data/split [get]
func (h *Handler) GetMeasurementDataSplit(w http.ResponseWriter, r *http.Request) {
	ds, err := h.validateDataSourceOwnership(r)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondJSON(w, http.StatusUnauthorized, ErrorResponse{
				Error:   "unauthorized",
				Message: "unauthorized",
			})
			return
		}
		if err.Error() == "data source not found" || err.Error() == "invalid data source ID" {
			respondJSON(w, http.StatusNotFound, ErrorResponse{
				Error:   "not_found",
				Message: "data source not found",
			})
			return
		}
		log.Printf("validate data source ownership error: %v", err)
		respondJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "failed to validate data source",
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

	splitByKey := r.URL.Query().Get("splitBy")
	if splitByKey == "" {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: "splitBy parameter is required",
		})
		return
	}

	startDate, endDate, err := h.parseDateRange(r)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, ErrorResponse{
			Error:   "validation_failed",
			Message: err.Error(),
		})
		return
	}

	metadataFilters := h.parseMetadataFilters(r)

	series, err := h.service.GetAggregatedMeasurementsSplitBy(r.Context(), ds.ID, measurementName, startDate, endDate, metadataFilters, splitByKey)
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
}

// parseDateRange extracts and validates start/end dates from request query params.
func (h *Handler) parseDateRange(r *http.Request) (time.Time, time.Time, error) {
	startStr := r.URL.Query().Get("start")
	if startStr == "" {
		return time.Time{}, time.Time{}, errors.New("start date is required")
	}
	startDate, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		startDate, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("invalid start date format (use ISO 8601)")
		}
	}

	endStr := r.URL.Query().Get("end")
	if endStr == "" {
		return time.Time{}, time.Time{}, errors.New("end date is required")
	}
	endDate, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		endDate, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return time.Time{}, time.Time{}, errors.New("invalid end date format (use ISO 8601)")
		}
		// For date-only, add one day to make it inclusive
		endDate = endDate.AddDate(0, 0, 1)
	}

	return startDate, endDate, nil
}

// parseMetadataFilters extracts metadata.key=value filters from request query params.
func (h *Handler) parseMetadataFilters(r *http.Request) map[string]string {
	metadataFilters := make(map[string]string)
	for key, values := range r.URL.Query() {
		if strings.HasPrefix(key, "metadata.") && len(values) > 0 {
			metadataKey := strings.TrimPrefix(key, "metadata.")
			metadataFilters[metadataKey] = values[0]
		}
	}
	return metadataFilters
}
