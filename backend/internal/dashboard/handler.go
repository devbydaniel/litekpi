package dashboard

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// Handler handles HTTP requests for dashboards.
type Handler struct {
	service *Service
}

// NewHandler creates a new dashboard handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
