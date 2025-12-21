package report

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

// Handler handles HTTP requests for reports.
type Handler struct {
	service *Service
}

// NewHandler creates a new report handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListReports handles listing all reports for the organization.
//
//	@Summary		List reports
//	@Description	Get all reports for the authenticated user's organization
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListReportsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/reports [get]
func (h *Handler) ListReports(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reports, err := h.service.ListReports(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("list reports error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list reports")
		return
	}

	respondJSON(w, http.StatusOK, ListReportsResponse{Reports: reports})
}

// GetReport handles getting a report with its KPIs.
//
//	@Summary		Get report
//	@Description	Get a report with its KPIs
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Report ID"
//	@Success		200	{object}	ReportWithKPIs
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/reports/{id} [get]
func (h *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	result, err := h.service.GetReport(r.Context(), user.OrganizationID, reportID)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("get report error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get report")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// CreateReport handles creating a new report.
//
//	@Summary		Create report
//	@Description	Create a new report. Requires editor or admin role.
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateReportRequest	true	"Report data"
//	@Success		201		{object}	Report
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/reports [post]
func (h *Handler) CreateReport(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	report, err := h.service.CreateReport(r.Context(), user.OrganizationID, req)
	if err != nil {
		if errors.Is(err, ErrReportNameEmpty) {
			respondError(w, http.StatusBadRequest, "report name is required")
			return
		}
		if errors.Is(err, ErrReportNameTooLong) {
			respondError(w, http.StatusBadRequest, "report name exceeds maximum length")
			return
		}
		log.Printf("create report error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create report")
		return
	}

	respondJSON(w, http.StatusCreated, report)
}

// UpdateReport handles updating a report.
//
//	@Summary		Update report
//	@Description	Update a report's name. Requires editor or admin role.
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Report ID"
//	@Param			request	body		UpdateReportRequest	true	"Report data"
//	@Success		200		{object}	Report
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/reports/{id} [put]
func (h *Handler) UpdateReport(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	var req UpdateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	report, err := h.service.UpdateReport(r.Context(), user.OrganizationID, reportID, req)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if errors.Is(err, ErrReportNameEmpty) {
			respondError(w, http.StatusBadRequest, "report name is required")
			return
		}
		if errors.Is(err, ErrReportNameTooLong) {
			respondError(w, http.StatusBadRequest, "report name exceeds maximum length")
			return
		}
		log.Printf("update report error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update report")
		return
	}

	respondJSON(w, http.StatusOK, report)
}

// DeleteReport handles deleting a report.
//
//	@Summary		Delete report
//	@Description	Delete a report. Requires editor or admin role.
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Report ID"
//	@Success		200	{object}	MessageResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/reports/{id} [delete]
func (h *Handler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	err = h.service.DeleteReport(r.Context(), user.OrganizationID, reportID)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("delete report error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete report")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "report deleted"})
}

// ComputeReport handles computing KPI values for a report.
//
//	@Summary		Compute report
//	@Description	Get a report with computed KPI values
//	@Tags			reports
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Report ID"
//	@Success		200	{object}	ComputedReport
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/reports/{id}/compute [get]
func (h *Handler) ComputeReport(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	result, err := h.service.ComputeReport(r.Context(), user.OrganizationID, reportID)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("compute report error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to compute report")
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// KPI handlers

// ListKPIs handles listing all KPIs for a report.
//
//	@Summary		List report KPIs
//	@Description	Get all KPIs for a report
//	@Tags			report-kpis
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Report ID"
//	@Success		200	{object}	kpi.ListKPIsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/reports/{id}/kpis [get]
func (h *Handler) ListKPIs(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	kpis, err := h.service.GetKPIs(r.Context(), user.OrganizationID, reportID)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("list KPIs error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list KPIs")
		return
	}

	respondJSON(w, http.StatusOK, kpi.ListKPIsResponse{KPIs: kpis})
}

// CreateKPI handles creating a new KPI on a report.
//
//	@Summary		Create report KPI
//	@Description	Create a new KPI on a report. Requires editor or admin role.
//	@Tags			report-kpis
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Report ID"
//	@Param			request	body		kpi.CreateKPIRequest	true	"KPI data"
//	@Success		201		{object}	kpi.KPI
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/reports/{id}/kpis [post]
func (h *Handler) CreateKPI(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	var req kpi.CreateKPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	newKPI, err := h.service.CreateKPI(r.Context(), user.OrganizationID, reportID, req)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
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
//	@Summary		Update report KPI
//	@Description	Update a KPI's configuration. Requires editor or admin role.
//	@Tags			report-kpis
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Report ID"
//	@Param			kpiId	path		string				true	"KPI ID"
//	@Param			request	body		kpi.UpdateKPIRequest	true	"KPI data"
//	@Success		200		{object}	kpi.KPI
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/reports/{id}/kpis/{kpiId} [put]
func (h *Handler) UpdateKPI(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	kpiID, err := uuid.Parse(chi.URLParam(r, "kpiId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid KPI ID")
		return
	}

	var req kpi.UpdateKPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedKPI, err := h.service.UpdateKPI(r.Context(), user.OrganizationID, reportID, kpiID, req)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if errors.Is(err, kpi.ErrKPINotFound) {
			respondError(w, http.StatusNotFound, "KPI not found")
			return
		}
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
//	@Summary		Delete report KPI
//	@Description	Delete a KPI from a report. Requires editor or admin role.
//	@Tags			report-kpis
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string	true	"Report ID"
//	@Param			kpiId	path		string	true	"KPI ID"
//	@Success		200		{object}	MessageResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/reports/{id}/kpis/{kpiId} [delete]
func (h *Handler) DeleteKPI(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	kpiID, err := uuid.Parse(chi.URLParam(r, "kpiId"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid KPI ID")
		return
	}

	err = h.service.DeleteKPI(r.Context(), user.OrganizationID, reportID, kpiID)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if errors.Is(err, kpi.ErrKPINotFound) {
			respondError(w, http.StatusNotFound, "KPI not found")
			return
		}
		log.Printf("delete KPI error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete KPI")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "KPI deleted"})
}

// ReorderKPIs handles reordering KPIs on a report.
//
//	@Summary		Reorder report KPIs
//	@Description	Reorder KPIs on a report. Requires editor or admin role.
//	@Tags			report-kpis
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Report ID"
//	@Param			request	body		kpi.ReorderKPIsRequest	true	"KPI order"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/reports/{id}/kpis/reorder [put]
func (h *Handler) ReorderKPIs(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	reportID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid report ID")
		return
	}

	var req kpi.ReorderKPIsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.service.ReorderKPIs(r.Context(), user.OrganizationID, reportID, req.KPIIDs)
	if err != nil {
		if errors.Is(err, ErrReportNotFound) {
			respondError(w, http.StatusNotFound, "report not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
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
