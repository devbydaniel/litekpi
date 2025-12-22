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
	return &Handler{
		service: service,
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

// GetDashboard handles getting a dashboard.
//
//	@Summary		Get dashboard
//	@Description	Get a dashboard by ID. Metrics are fetched separately via /metrics endpoints.
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
//	@Description	Get the default dashboard. Metrics are fetched separately via /metrics endpoints.
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

// VerifyDashboardOwnership verifies that a dashboard belongs to an organization.
// This is used by other handlers (e.g., metric handler) to check ownership.
func (h *Handler) VerifyDashboardOwnership(w http.ResponseWriter, r *http.Request) (*Dashboard, bool) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}

	dashboardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid dashboard ID")
		return nil, false
	}

	dashboard, err := h.service.VerifyDashboardOwnership(r.Context(), user.OrganizationID, dashboardID)
	if err != nil {
		if errors.Is(err, ErrDashboardNotFound) {
			respondError(w, http.StatusNotFound, "dashboard not found")
			return nil, false
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return nil, false
		}
		log.Printf("verify dashboard ownership error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to verify dashboard")
		return nil, false
	}

	return dashboard, true
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
