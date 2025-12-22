package dashboard

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/auth"
	"github.com/devbydaniel/litekpi/internal/scalarmetric"
	"github.com/devbydaniel/litekpi/internal/timeseries"
)

// Handler handles HTTP requests for dashboards.
type Handler struct {
	service             *Service
	timeSeriesService   *timeseries.Service
	scalarMetricService *scalarmetric.Service
}

// NewHandler creates a new dashboard handler.
func NewHandler(service *Service, timeSeriesService *timeseries.Service, scalarMetricService *scalarmetric.Service) *Handler {
	return &Handler{
		service:             service,
		timeSeriesService:   timeSeriesService,
		scalarMetricService: scalarMetricService,
	}
}

// ListDashboards handles listing all dashboards for the organization.
//
//	@Summary		List dashboards
//	@Description	Get all dashboards for the authenticated user's organization
//	@Tags			dashboards
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListDashboardsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards [get]
func (h *Handler) ListDashboards(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dashboards, err := h.service.ListDashboards(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("list dashboards error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list dashboards")
		return
	}

	respondJSON(w, http.StatusOK, ListDashboardsResponse{Dashboards: dashboards})
}

// GetDashboard handles getting a dashboard with its time series and scalar metrics.
//
//	@Summary		Get dashboard
//	@Description	Get a dashboard with its time series and scalar metrics
//	@Tags			dashboards
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	DashboardWithData
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards/{id} [get]
func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.service.GetDashboard(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("get dashboard error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get dashboard")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// GetDefaultDashboard handles getting the default dashboard.
//
//	@Summary		Get default dashboard
//	@Description	Get the default dashboard with its time series and scalar metrics
//	@Tags			dashboards
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	DashboardWithData
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards/default [get]
func (h *Handler) GetDefaultDashboard(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	result, err := h.service.GetDefaultDashboard(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("get default dashboard error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get default dashboard")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// CreateDashboard handles creating a new dashboard.
//
//	@Summary		Create dashboard
//	@Description	Create a new dashboard. Requires editor or admin role.
//	@Tags			dashboards
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateDashboardRequest	true	"Dashboard data"
//	@Success		201		{object}	Dashboard
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards [post]
func (h *Handler) CreateDashboard(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateDashboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dashboard, err := h.service.CreateDashboard(r.Context(), user.OrganizationID, req)
	if err != nil {
		if errors.Is(err, ErrDashboardNameEmpty) {
			respondError(w, http.StatusBadRequest, "dashboard name is required")
			return
		}
		log.Printf("create dashboard error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create dashboard")
		return
	}

	respondJSON(w, http.StatusCreated, dashboard)
}

// UpdateDashboard handles updating a dashboard.
//
//	@Summary		Update dashboard
//	@Description	Update a dashboard's name. Requires editor or admin role.
//	@Tags			dashboards
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Dashboard ID"
//	@Param			request	body		UpdateDashboardRequest	true	"Dashboard data"
//	@Success		200		{object}	Dashboard
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id} [put]
func (h *Handler) UpdateDashboard(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateDashboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dashboard, err := h.service.UpdateDashboard(r.Context(), user.OrganizationID, dashboardID, req)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if errors.Is(err, ErrDashboardNameEmpty) {
			respondError(w, http.StatusBadRequest, "dashboard name is required")
			return
		}
		log.Printf("update dashboard error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update dashboard")
		return
	}

	respondJSON(w, http.StatusOK, dashboard)
}

// DeleteDashboard handles deleting a dashboard.
//
//	@Summary		Delete dashboard
//	@Description	Delete a dashboard. Cannot delete the default dashboard. Requires editor or admin role.
//	@Tags			dashboards
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id} [delete]
func (h *Handler) DeleteDashboard(w http.ResponseWriter, r *http.Request) {
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

	err = h.service.DeleteDashboard(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if errors.Is(err, ErrCannotDeleteDefault) {
			respondError(w, http.StatusBadRequest, "cannot delete default dashboard")
			return
		}
		log.Printf("delete dashboard error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete dashboard")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "dashboard deleted"})
}

// Time Series handlers

// CreateTimeSeries handles creating a new time series.
//
//	@Summary		Create time series
//	@Description	Create a new time series on a dashboard. Requires editor or admin role.
//	@Tags			time-series
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string						true	"Dashboard ID"
//	@Param			request	body		timeseries.CreateTimeSeriesRequest	true	"Time series data"
//	@Success		201		{object}	timeseries.TimeSeries
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/time-series [post]
func (h *Handler) CreateTimeSeries(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req timeseries.CreateTimeSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ts, err := h.timeSeriesService.Create(r.Context(), user.OrganizationID, dashboardID, req)
	if err != nil {
		if errors.Is(err, timeseries.ErrMeasurementNameEmpty) {
			respondError(w, http.StatusBadRequest, "measurement name is required")
			return
		}
		if errors.Is(err, timeseries.ErrInvalidChartType) {
			respondError(w, http.StatusBadRequest, "invalid chart type")
			return
		}
		if errors.Is(err, timeseries.ErrInvalidDateRange) {
			respondError(w, http.StatusBadRequest, "invalid date range")
			return
		}
		if errors.Is(err, timeseries.ErrTitleTooLong) {
			respondError(w, http.StatusBadRequest, "title exceeds maximum length")
			return
		}
		log.Printf("create time series error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create time series")
		return
	}

	respondJSON(w, http.StatusCreated, ts)
}

// UpdateTimeSeries handles updating a time series.
//
//	@Summary		Update time series
//	@Description	Update a time series configuration. Requires editor or admin role.
//	@Tags			time-series
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path		string						true	"Dashboard ID"
//	@Param			timeSeriesId	path		string						true	"Time Series ID"
//	@Param			request			body		timeseries.UpdateTimeSeriesRequest	true	"Time series data"
//	@Success		200				{object}	timeseries.TimeSeries
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/time-series/{timeSeriesId} [put]
func (h *Handler) UpdateTimeSeries(w http.ResponseWriter, r *http.Request) {
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

	timeSeriesID, err := uuid.Parse(chi.URLParam(r, "timeSeriesId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid time series ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req timeseries.UpdateTimeSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ts, err := h.timeSeriesService.Update(r.Context(), dashboardID, timeSeriesID, req)
	if err != nil {
		if errors.Is(err, timeseries.ErrTimeSeriesNotFound) {
			respondError(w, http.StatusNotFound, "time series not found")
			return
		}
		if errors.Is(err, timeseries.ErrInvalidChartType) {
			respondError(w, http.StatusBadRequest, "invalid chart type")
			return
		}
		if errors.Is(err, timeseries.ErrInvalidDateRange) {
			respondError(w, http.StatusBadRequest, "invalid date range")
			return
		}
		if errors.Is(err, timeseries.ErrTitleTooLong) {
			respondError(w, http.StatusBadRequest, "title exceeds maximum length")
			return
		}
		log.Printf("update time series error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update time series")
		return
	}

	respondJSON(w, http.StatusOK, ts)
}

// DeleteTimeSeries handles deleting a time series.
//
//	@Summary		Delete time series
//	@Description	Delete a time series from a dashboard. Requires editor or admin role.
//	@Tags			time-series
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path		string	true	"Dashboard ID"
//	@Param			timeSeriesId	path		string	true	"Time Series ID"
//	@Success		200				{object}	MessageResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/time-series/{timeSeriesId} [delete]
func (h *Handler) DeleteTimeSeries(w http.ResponseWriter, r *http.Request) {
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

	timeSeriesID, err := uuid.Parse(chi.URLParam(r, "timeSeriesId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid time series ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	err = h.timeSeriesService.Delete(r.Context(), dashboardID, timeSeriesID)
	if err != nil {
		if errors.Is(err, timeseries.ErrTimeSeriesNotFound) {
			respondError(w, http.StatusNotFound, "time series not found")
			return
		}
		log.Printf("delete time series error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete time series")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "time series deleted"})
}

// ReorderTimeSeries handles reordering time series on a dashboard.
//
//	@Summary		Reorder time series
//	@Description	Reorder time series on a dashboard. Requires editor or admin role.
//	@Tags			time-series
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Dashboard ID"
//	@Param			request	body		timeseries.ReorderTimeSeriesRequest	true	"Time series order"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/time-series/reorder [put]
func (h *Handler) ReorderTimeSeries(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req timeseries.ReorderTimeSeriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.timeSeriesService.Reorder(r.Context(), dashboardID, req.TimeSeriesIDs)
	if err != nil {
		log.Printf("reorder time series error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to reorder time series")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "time series reordered"})
}

// Scalar Metric handlers

// ListScalarMetrics handles listing all scalar metrics for a dashboard.
//
//	@Summary		List dashboard scalar metrics
//	@Description	Get all scalar metrics for a dashboard
//	@Tags			scalar-metrics
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	scalarmetric.ListScalarMetricsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards/{id}/scalar-metrics [get]
func (h *Handler) ListScalarMetrics(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	metrics, err := h.scalarMetricService.GetByDashboardID(r.Context(), dashboardID)
	if err != nil {
		log.Printf("list scalar metrics error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list scalar metrics")
		return
	}

	respondJSON(w, http.StatusOK, scalarmetric.ListScalarMetricsResponse{ScalarMetrics: metrics})
}

// CreateScalarMetric handles creating a new scalar metric on a dashboard.
//
//	@Summary		Create dashboard scalar metric
//	@Description	Create a new scalar metric on a dashboard. Requires editor or admin role.
//	@Tags			scalar-metrics
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string							true	"Dashboard ID"
//	@Param			request	body		scalarmetric.CreateScalarMetricRequest	true	"Scalar metric data"
//	@Success		201		{object}	scalarmetric.ScalarMetric
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/scalar-metrics [post]
func (h *Handler) CreateScalarMetric(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req scalarmetric.CreateScalarMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	metric, err := h.scalarMetricService.Create(r.Context(), user.OrganizationID, dashboardID, req)
	if err != nil {
		if errors.Is(err, scalarmetric.ErrLabelEmpty) {
			respondError(w, http.StatusBadRequest, "label is required")
			return
		}
		if errors.Is(err, scalarmetric.ErrLabelTooLong) {
			respondError(w, http.StatusBadRequest, "label exceeds maximum length")
			return
		}
		if errors.Is(err, scalarmetric.ErrMeasurementNameEmpty) {
			respondError(w, http.StatusBadRequest, "measurement name is required")
			return
		}
		if errors.Is(err, scalarmetric.ErrInvalidTimeframe) {
			respondError(w, http.StatusBadRequest, "invalid timeframe")
			return
		}
		if errors.Is(err, scalarmetric.ErrInvalidAggregation) {
			respondError(w, http.StatusBadRequest, "invalid aggregation type")
			return
		}
		if errors.Is(err, scalarmetric.ErrInvalidComparisonType) {
			respondError(w, http.StatusBadRequest, "invalid comparison display type")
			return
		}
		log.Printf("create scalar metric error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create scalar metric")
		return
	}

	respondJSON(w, http.StatusCreated, metric)
}

// UpdateScalarMetric handles updating a scalar metric.
//
//	@Summary		Update dashboard scalar metric
//	@Description	Update a scalar metric's configuration. Requires editor or admin role.
//	@Tags			scalar-metrics
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path		string							true	"Dashboard ID"
//	@Param			scalarMetricId	path		string							true	"Scalar Metric ID"
//	@Param			request			body		scalarmetric.UpdateScalarMetricRequest	true	"Scalar metric data"
//	@Success		200				{object}	scalarmetric.ScalarMetric
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/scalar-metrics/{scalarMetricId} [put]
func (h *Handler) UpdateScalarMetric(w http.ResponseWriter, r *http.Request) {
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

	scalarMetricID, err := uuid.Parse(chi.URLParam(r, "scalarMetricId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid scalar metric ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req scalarmetric.UpdateScalarMetricRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	metric, err := h.scalarMetricService.Update(r.Context(), dashboardID, scalarMetricID, req)
	if err != nil {
		if errors.Is(err, scalarmetric.ErrScalarMetricNotFound) {
			respondError(w, http.StatusNotFound, "scalar metric not found")
			return
		}
		if errors.Is(err, scalarmetric.ErrLabelEmpty) {
			respondError(w, http.StatusBadRequest, "label is required")
			return
		}
		if errors.Is(err, scalarmetric.ErrLabelTooLong) {
			respondError(w, http.StatusBadRequest, "label exceeds maximum length")
			return
		}
		if errors.Is(err, scalarmetric.ErrInvalidTimeframe) {
			respondError(w, http.StatusBadRequest, "invalid timeframe")
			return
		}
		if errors.Is(err, scalarmetric.ErrInvalidAggregation) {
			respondError(w, http.StatusBadRequest, "invalid aggregation type")
			return
		}
		if errors.Is(err, scalarmetric.ErrInvalidComparisonType) {
			respondError(w, http.StatusBadRequest, "invalid comparison display type")
			return
		}
		log.Printf("update scalar metric error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update scalar metric")
		return
	}

	respondJSON(w, http.StatusOK, metric)
}

// DeleteScalarMetric handles deleting a scalar metric.
//
//	@Summary		Delete dashboard scalar metric
//	@Description	Delete a scalar metric from a dashboard. Requires editor or admin role.
//	@Tags			scalar-metrics
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id				path		string	true	"Dashboard ID"
//	@Param			scalarMetricId	path		string	true	"Scalar Metric ID"
//	@Success		200				{object}	MessageResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/scalar-metrics/{scalarMetricId} [delete]
func (h *Handler) DeleteScalarMetric(w http.ResponseWriter, r *http.Request) {
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

	scalarMetricID, err := uuid.Parse(chi.URLParam(r, "scalarMetricId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid scalar metric ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	err = h.scalarMetricService.Delete(r.Context(), dashboardID, scalarMetricID)
	if err != nil {
		if errors.Is(err, scalarmetric.ErrScalarMetricNotFound) {
			respondError(w, http.StatusNotFound, "scalar metric not found")
			return
		}
		log.Printf("delete scalar metric error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete scalar metric")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "scalar metric deleted"})
}

// ComputeScalarMetrics handles computing scalar metric values for a dashboard.
//
//	@Summary		Compute dashboard scalar metrics
//	@Description	Get computed values for all scalar metrics on a dashboard
//	@Tags			scalar-metrics
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	scalarmetric.ComputeScalarMetricsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards/{id}/scalar-metrics/compute [get]
func (h *Handler) ComputeScalarMetrics(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	metrics, err := h.scalarMetricService.GetByDashboardID(r.Context(), dashboardID)
	if err != nil {
		log.Printf("list scalar metrics error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list scalar metrics")
		return
	}

	computed, err := h.scalarMetricService.Compute(r.Context(), metrics)
	if err != nil {
		log.Printf("compute scalar metrics error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to compute scalar metrics")
		return
	}

	respondJSON(w, http.StatusOK, scalarmetric.ComputeScalarMetricsResponse{ScalarMetrics: computed})
}

// ReorderScalarMetrics handles reordering scalar metrics on a dashboard.
//
//	@Summary		Reorder dashboard scalar metrics
//	@Description	Reorder scalar metrics on a dashboard. Requires editor or admin role.
//	@Tags			scalar-metrics
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string								true	"Dashboard ID"
//	@Param			request	body		scalarmetric.ReorderScalarMetricsRequest	true	"Scalar metric order"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/scalar-metrics/reorder [put]
func (h *Handler) ReorderScalarMetrics(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return
	}

	var req scalarmetric.ReorderScalarMetricsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.scalarMetricService.Reorder(r.Context(), dashboardID, req.ScalarMetricIDs); err != nil {
		log.Printf("reorder scalar metrics error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to reorder scalar metrics")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "scalar metrics reordered"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
