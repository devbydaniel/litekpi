package metric

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/auth"
	"github.com/devbydaniel/litekpi/internal/dashboard"
)

// Handler handles HTTP requests for metrics.
type Handler struct {
	service          *Service
	dashboardService *dashboard.Service
}

// NewHandler creates a new metric handler.
func NewHandler(service *Service, dashboardService *dashboard.Service) *Handler {
	return &Handler{
		service:          service,
		dashboardService: dashboardService,
	}
}

// ListMetrics handles listing all metrics for a dashboard.
//
//	@Summary		List dashboard metrics
//	@Description	Get all metrics for a dashboard
//	@Tags			metrics
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	ListMetricsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards/{id}/metrics [get]
func (h *Handler) ListMetrics(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dashboardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid dashboard ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.dashboardService.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, dashboard.ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, dashboard.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	metrics, err := h.service.GetByDashboardID(r.Context(), dashboardID)
	if err != nil {
		log.Printf("list metrics error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list metrics")
		return
	}

	respondJSON(w, http.StatusOK, ListMetricsResponse{Metrics: metrics})
}

// CreateMetric handles creating a new metric on a dashboard.
//
//	@Summary		Create dashboard metric
//	@Description	Create a new metric on a dashboard. Requires editor or admin role.
//	@Tags			metrics
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Dashboard ID"
//	@Param			request	body		CreateMetricRequest	true	"Metric data"
//	@Success		201		{object}	Metric
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/metrics [post]
func (h *Handler) CreateMetric(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dashboardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid dashboard ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.dashboardService.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, dashboard.ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, dashboard.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req CreateMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	metric, err := h.service.Create(r.Context(), user.OrganizationID, dashboardID, req)
	if err != nil {
		if errors.Is(err, ErrLabelEmpty) {
			respondError(w, http.StatusBadRequest, "label is required")
			return
		}
		if errors.Is(err, ErrLabelTooLong) {
			respondError(w, http.StatusBadRequest, "label exceeds maximum length")
			return
		}
		if errors.Is(err, ErrMeasurementNameEmpty) {
			respondError(w, http.StatusBadRequest, "measurement name is required")
			return
		}
		if errors.Is(err, ErrInvalidTimeframe) {
			respondError(w, http.StatusBadRequest, "invalid timeframe")
			return
		}
		if errors.Is(err, ErrInvalidAggregation) {
			respondError(w, http.StatusBadRequest, "invalid aggregation type")
			return
		}
		if errors.Is(err, ErrAggregationKeyRequired) {
			respondError(w, http.StatusBadRequest, "aggregation_key is required for count_unique aggregation")
			return
		}
		if errors.Is(err, ErrInvalidGranularity) {
			respondError(w, http.StatusBadRequest, "invalid granularity")
			return
		}
		if errors.Is(err, ErrInvalidDisplayMode) {
			respondError(w, http.StatusBadRequest, "invalid display mode")
			return
		}
		if errors.Is(err, ErrChartTypeRequired) {
			respondError(w, http.StatusBadRequest, "chart_type is required for time_series display mode")
			return
		}
		if errors.Is(err, ErrInvalidChartType) {
			respondError(w, http.StatusBadRequest, "invalid chart type")
			return
		}
		if errors.Is(err, ErrInvalidComparisonType) {
			respondError(w, http.StatusBadRequest, "invalid comparison display type")
			return
		}
		log.Printf("create metric error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create metric")
		return
	}

	respondJSON(w, http.StatusCreated, metric)
}

// UpdateMetric handles updating a metric.
//
//	@Summary		Update dashboard metric
//	@Description	Update a metric's configuration. Requires editor or admin role.
//	@Tags			metrics
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		string				true	"Dashboard ID"
//	@Param			metricId	path		string				true	"Metric ID"
//	@Param			request		body		UpdateMetricRequest	true	"Metric data"
//	@Success		200			{object}	Metric
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		403			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/dashboards/{id}/metrics/{metricId} [put]
func (h *Handler) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dashboardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid dashboard ID")
		return
	}

	metricID, err := uuid.Parse(chi.URLParam(r, "metricId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid metric ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.dashboardService.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, dashboard.ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, dashboard.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req UpdateMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	metric, err := h.service.Update(r.Context(), dashboardID, metricID, req)
	if err != nil {
		if errors.Is(err, ErrMetricNotFound) {
			respondError(w, http.StatusNotFound, "metric not found")
			return
		}
		if errors.Is(err, ErrLabelEmpty) {
			respondError(w, http.StatusBadRequest, "label is required")
			return
		}
		if errors.Is(err, ErrLabelTooLong) {
			respondError(w, http.StatusBadRequest, "label exceeds maximum length")
			return
		}
		if errors.Is(err, ErrInvalidTimeframe) {
			respondError(w, http.StatusBadRequest, "invalid timeframe")
			return
		}
		if errors.Is(err, ErrInvalidAggregation) {
			respondError(w, http.StatusBadRequest, "invalid aggregation type")
			return
		}
		if errors.Is(err, ErrAggregationKeyRequired) {
			respondError(w, http.StatusBadRequest, "aggregation_key is required for count_unique aggregation")
			return
		}
		if errors.Is(err, ErrInvalidGranularity) {
			respondError(w, http.StatusBadRequest, "invalid granularity")
			return
		}
		if errors.Is(err, ErrInvalidDisplayMode) {
			respondError(w, http.StatusBadRequest, "invalid display mode")
			return
		}
		if errors.Is(err, ErrChartTypeRequired) {
			respondError(w, http.StatusBadRequest, "chart_type is required for time_series display mode")
			return
		}
		if errors.Is(err, ErrInvalidChartType) {
			respondError(w, http.StatusBadRequest, "invalid chart type")
			return
		}
		if errors.Is(err, ErrInvalidComparisonType) {
			respondError(w, http.StatusBadRequest, "invalid comparison display type")
			return
		}
		log.Printf("update metric error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update metric")
		return
	}

	respondJSON(w, http.StatusOK, metric)
}

// DeleteMetric handles deleting a metric.
//
//	@Summary		Delete dashboard metric
//	@Description	Delete a metric from a dashboard. Requires editor or admin role.
//	@Tags			metrics
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		string	true	"Dashboard ID"
//	@Param			metricId	path		string	true	"Metric ID"
//	@Success		200			{object}	MessageResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		403			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/dashboards/{id}/metrics/{metricId} [delete]
func (h *Handler) DeleteMetric(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dashboardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid dashboard ID")
		return
	}

	metricID, err := uuid.Parse(chi.URLParam(r, "metricId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid metric ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.dashboardService.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, dashboard.ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, dashboard.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	err = h.service.Delete(r.Context(), dashboardID, metricID)
	if err != nil {
		if errors.Is(err, ErrMetricNotFound) {
			respondError(w, http.StatusNotFound, "metric not found")
			return
		}
		log.Printf("delete metric error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete metric")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "metric deleted"})
}

// ComputeMetrics handles computing metric values for a dashboard.
//
//	@Summary		Compute dashboard metrics
//	@Description	Get computed values for all metrics on a dashboard
//	@Tags			metrics
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	ComputeMetricsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards/{id}/metrics/compute [get]
func (h *Handler) ComputeMetrics(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dashboardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid dashboard ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.dashboardService.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, dashboard.ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, dashboard.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	metrics, err := h.service.GetByDashboardID(r.Context(), dashboardID)
	if err != nil {
		log.Printf("list metrics error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list metrics")
		return
	}

	computed, err := h.service.Compute(r.Context(), metrics)
	if err != nil {
		log.Printf("compute metrics error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to compute metrics")
		return
	}

	respondJSON(w, http.StatusOK, ComputeMetricsResponse{Metrics: computed})
}

// ReorderMetrics handles reordering metrics on a dashboard.
//
//	@Summary		Reorder dashboard metrics
//	@Description	Reorder metrics on a dashboard. Requires editor or admin role.
//	@Tags			metrics
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Dashboard ID"
//	@Param			request	body		ReorderMetricsRequest	true	"Metric order"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/metrics/reorder [put]
func (h *Handler) ReorderMetrics(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dashboardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid dashboard ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.dashboardService.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, dashboard.ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, dashboard.ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req ReorderMetricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Reorder(r.Context(), dashboardID, req.MetricIDs); err != nil {
		log.Printf("reorder metrics error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to reorder metrics")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "metrics reordered"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
