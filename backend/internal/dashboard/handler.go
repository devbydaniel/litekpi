package dashboard

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/auth"
	"github.com/devbydaniel/litekpi/internal/kpi"
)

// Handler handles HTTP requests for dashboards.
type Handler struct {
	service    *Service
	kpiService *kpi.Service
}

// NewHandler creates a new dashboard handler.
func NewHandler(service *Service, kpiService *kpi.Service) *Handler {
	return &Handler{
		service:    service,
		kpiService: kpiService,
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

// GetDashboard handles getting a dashboard with its widgets.
//
//	@Summary		Get dashboard
//	@Description	Get a dashboard with its widgets
//	@Tags			dashboards
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	DashboardWithWidgets
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
//	@Description	Get the default dashboard with its widgets
//	@Tags			dashboards
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	DashboardWithWidgets
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

// Widget handlers

// CreateWidget handles creating a new widget.
//
//	@Summary		Create widget
//	@Description	Create a new widget on a dashboard. Requires editor or admin role.
//	@Tags			widgets
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Dashboard ID"
//	@Param			request	body		CreateWidgetRequest	true	"Widget data"
//	@Success		201		{object}	Widget
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/widgets [post]
func (h *Handler) CreateWidget(w http.ResponseWriter, r *http.Request) {
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

	var req CreateWidgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	widget, err := h.service.CreateWidget(r.Context(), user.OrganizationID, dashboardID, req)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if errors.Is(err, ErrInvalidChartType) {
			respondError(w, http.StatusBadRequest, "invalid chart type")
			return
		}
		if errors.Is(err, ErrInvalidDateRange) {
			respondError(w, http.StatusBadRequest, "invalid date range")
			return
		}
		log.Printf("create widget error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create widget")
		return
	}

	respondJSON(w, http.StatusCreated, widget)
}

// UpdateWidget handles updating a widget.
//
//	@Summary		Update widget
//	@Description	Update a widget's configuration. Requires editor or admin role.
//	@Tags			widgets
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		string				true	"Dashboard ID"
//	@Param			widgetId	path		string				true	"Widget ID"
//	@Param			request		body		UpdateWidgetRequest	true	"Widget data"
//	@Success		200			{object}	Widget
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/widgets/{widgetId} [put]
func (h *Handler) UpdateWidget(w http.ResponseWriter, r *http.Request) {
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

	widgetID, err := uuid.Parse(chi.URLParam(r, "widgetId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid widget ID")
		return
	}

	var req UpdateWidgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	widget, err := h.service.UpdateWidget(r.Context(), user.OrganizationID, dashboardID, widgetID, req)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrWidgetNotFound) {
			respondError(w, http.StatusNotFound, "widget not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if errors.Is(err, ErrInvalidChartType) {
			respondError(w, http.StatusBadRequest, "invalid chart type")
			return
		}
		if errors.Is(err, ErrInvalidDateRange) {
			respondError(w, http.StatusBadRequest, "invalid date range")
			return
		}
		log.Printf("update widget error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update widget")
		return
	}

	respondJSON(w, http.StatusOK, widget)
}

// DeleteWidget handles deleting a widget.
//
//	@Summary		Delete widget
//	@Description	Delete a widget from a dashboard. Requires editor or admin role.
//	@Tags			widgets
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id			path		string	true	"Dashboard ID"
//	@Param			widgetId	path		string	true	"Widget ID"
//	@Success		200			{object}	MessageResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/widgets/{widgetId} [delete]
func (h *Handler) DeleteWidget(w http.ResponseWriter, r *http.Request) {
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

	widgetID, err := uuid.Parse(chi.URLParam(r, "widgetId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid widget ID")
		return
	}

	err = h.service.DeleteWidget(r.Context(), user.OrganizationID, dashboardID, widgetID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrWidgetNotFound) {
			respondError(w, http.StatusNotFound, "widget not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("delete widget error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete widget")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "widget deleted"})
}

// ReorderWidgets handles reordering widgets on a dashboard.
//
//	@Summary		Reorder widgets
//	@Description	Reorder widgets on a dashboard. Requires editor or admin role.
//	@Tags			widgets
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string					true	"Dashboard ID"
//	@Param			request	body		ReorderWidgetsRequest	true	"Widget order"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/widgets/reorder [put]
func (h *Handler) ReorderWidgets(w http.ResponseWriter, r *http.Request) {
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

	var req ReorderWidgetsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.service.ReorderWidgets(r.Context(), user.OrganizationID, dashboardID, req.WidgetIDs)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("reorder widgets error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to reorder widgets")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "widgets reordered"})
}

// KPI handlers

// ListKPIs handles listing all KPIs for a dashboard.
//
//	@Summary		List dashboard KPIs
//	@Description	Get all KPIs for a dashboard
//	@Tags			kpis
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	kpi.ListKPIsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards/{id}/kpis [get]
func (h *Handler) ListKPIs(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.GetDashboard(r.Context(), user.OrganizationID, dashboardID)
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

	kpis, err := h.kpiService.GetKPIsByDashboardID(r.Context(), dashboardID)
	if err != nil {
		log.Printf("list KPIs error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list KPIs")
		return
	}

	respondJSON(w, http.StatusOK, kpi.ListKPIsResponse{KPIs: kpis})
}

// CreateKPI handles creating a new KPI on a dashboard.
//
//	@Summary		Create dashboard KPI
//	@Description	Create a new KPI on a dashboard. Requires editor or admin role.
//	@Tags			kpis
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Dashboard ID"
//	@Param			request	body		kpi.CreateKPIRequest	true	"KPI data"
//	@Success		201		{object}	kpi.KPI
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/kpis [post]
func (h *Handler) CreateKPI(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.GetDashboard(r.Context(), user.OrganizationID, dashboardID)
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

	var req kpi.CreateKPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	newKPI, err := h.kpiService.CreateKPIForDashboard(r.Context(), user.OrganizationID, dashboardID, req)
	if err != nil {
		if errors.Is(err, kpi.ErrLabelEmpty) {
			respondError(w, http.StatusBadRequest, "label is required")
			return
		}
		if errors.Is(err, kpi.ErrLabelTooLong) {
			respondError(w, http.StatusBadRequest, "label exceeds maximum length")
			return
		}
		if errors.Is(err, kpi.ErrMeasurementNameEmpty) {
			respondError(w, http.StatusBadRequest, "measurement name is required")
			return
		}
		if errors.Is(err, kpi.ErrInvalidTimeframe) {
			respondError(w, http.StatusBadRequest, "invalid timeframe")
			return
		}
		if errors.Is(err, kpi.ErrInvalidAggregation) {
			respondError(w, http.StatusBadRequest, "invalid aggregation type")
			return
		}
		if errors.Is(err, kpi.ErrInvalidComparisonType) {
			respondError(w, http.StatusBadRequest, "invalid comparison display type")
			return
		}
		log.Printf("create KPI error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create KPI")
		return
	}

	respondJSON(w, http.StatusCreated, newKPI)
}

// UpdateKPI handles updating a KPI.
//
//	@Summary		Update dashboard KPI
//	@Description	Update a KPI's configuration. Requires editor or admin role.
//	@Tags			kpis
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Dashboard ID"
//	@Param			kpiId	path		string				true	"KPI ID"
//	@Param			request	body		kpi.UpdateKPIRequest	true	"KPI data"
//	@Success		200		{object}	kpi.KPI
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/kpis/{kpiId} [put]
func (h *Handler) UpdateKPI(w http.ResponseWriter, r *http.Request) {
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

	kpiID, err := uuid.Parse(chi.URLParam(r, "kpiId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid KPI ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.service.GetDashboard(r.Context(), user.OrganizationID, dashboardID)
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

	// Verify KPI belongs to this dashboard
	existingKPI, err := h.kpiService.GetKPIByID(r.Context(), kpiID)
	if err != nil {
		if errors.Is(err, kpi.ErrKPINotFound) {
			respondError(w, http.StatusNotFound, "KPI not found")
			return
		}
		log.Printf("get KPI error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get KPI")
		return
	}
	if existingKPI.DashboardID == nil || *existingKPI.DashboardID != dashboardID {
		respondError(w, http.StatusNotFound, "KPI not found")
		return
	}

	var req kpi.UpdateKPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedKPI, err := h.kpiService.UpdateKPI(r.Context(), kpiID, req)
	if err != nil {
		if errors.Is(err, kpi.ErrLabelEmpty) {
			respondError(w, http.StatusBadRequest, "label is required")
			return
		}
		if errors.Is(err, kpi.ErrLabelTooLong) {
			respondError(w, http.StatusBadRequest, "label exceeds maximum length")
			return
		}
		if errors.Is(err, kpi.ErrInvalidTimeframe) {
			respondError(w, http.StatusBadRequest, "invalid timeframe")
			return
		}
		if errors.Is(err, kpi.ErrInvalidAggregation) {
			respondError(w, http.StatusBadRequest, "invalid aggregation type")
			return
		}
		if errors.Is(err, kpi.ErrInvalidComparisonType) {
			respondError(w, http.StatusBadRequest, "invalid comparison display type")
			return
		}
		log.Printf("update KPI error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update KPI")
		return
	}

	respondJSON(w, http.StatusOK, updatedKPI)
}

// DeleteKPI handles deleting a KPI.
//
//	@Summary		Delete dashboard KPI
//	@Description	Delete a KPI from a dashboard. Requires editor or admin role.
//	@Tags			kpis
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string	true	"Dashboard ID"
//	@Param			kpiId	path		string	true	"KPI ID"
//	@Success		200		{object}	MessageResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/kpis/{kpiId} [delete]
func (h *Handler) DeleteKPI(w http.ResponseWriter, r *http.Request) {
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

	kpiID, err := uuid.Parse(chi.URLParam(r, "kpiId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid KPI ID")
		return
	}

	// Verify dashboard ownership
	_, err = h.service.GetDashboard(r.Context(), user.OrganizationID, dashboardID)
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

	// Verify KPI belongs to this dashboard
	existingKPI, err := h.kpiService.GetKPIByID(r.Context(), kpiID)
	if err != nil {
		if errors.Is(err, kpi.ErrKPINotFound) {
			respondError(w, http.StatusNotFound, "KPI not found")
			return
		}
		log.Printf("get KPI error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get KPI")
		return
	}
	if existingKPI.DashboardID == nil || *existingKPI.DashboardID != dashboardID {
		respondError(w, http.StatusNotFound, "KPI not found")
		return
	}

	if err := h.kpiService.DeleteKPI(r.Context(), kpiID); err != nil {
		log.Printf("delete KPI error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete KPI")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "KPI deleted"})
}

// ComputeKPIs handles computing KPI values for a dashboard.
//
//	@Summary		Compute dashboard KPIs
//	@Description	Get computed values for all KPIs on a dashboard
//	@Tags			kpis
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Dashboard ID"
//	@Success		200	{object}	kpi.ComputeKPIsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/dashboards/{id}/kpis/compute [get]
func (h *Handler) ComputeKPIs(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.GetDashboard(r.Context(), user.OrganizationID, dashboardID)
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

	kpis, err := h.kpiService.GetKPIsByDashboardID(r.Context(), dashboardID)
	if err != nil {
		log.Printf("list KPIs error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list KPIs")
		return
	}

	computedKPIs, err := h.kpiService.ComputeKPIs(r.Context(), kpis)
	if err != nil {
		log.Printf("compute KPIs error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to compute KPIs")
		return
	}

	respondJSON(w, http.StatusOK, kpi.ComputeKPIsResponse{KPIs: computedKPIs})
}

// ReorderKPIs handles reordering KPIs on a dashboard.
//
//	@Summary		Reorder dashboard KPIs
//	@Description	Reorder KPIs on a dashboard. Requires editor or admin role.
//	@Tags			kpis
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Dashboard ID"
//	@Param			request	body		kpi.ReorderKPIsRequest	true	"KPI order"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/dashboards/{id}/kpis/reorder [put]
func (h *Handler) ReorderKPIs(w http.ResponseWriter, r *http.Request) {
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
	_, err = h.service.GetDashboard(r.Context(), user.OrganizationID, dashboardID)
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

	var req kpi.ReorderKPIsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.kpiService.ReorderKPIsForDashboard(r.Context(), dashboardID, req.KPIIDs); err != nil {
		log.Printf("reorder KPIs error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to reorder KPIs")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "KPIs reordered"})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
