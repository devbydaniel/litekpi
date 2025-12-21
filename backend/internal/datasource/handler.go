package datasource

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// Handler handles HTTP requests for data sources.
type Handler struct {
	service *Service
}

// NewHandler creates a new data source handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListDataSources handles listing data sources for the authenticated user's organization.
//
//	@Summary		List data sources
//	@Description	Get all data sources for the authenticated user's organization
//	@Tags			data-sources
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListDataSourcesResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/data-sources [get]
func (h *Handler) ListDataSources(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dataSources, err := h.service.ListDataSources(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("list data sources error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list data sources")
		return
	}

	respondJSON(w, http.StatusOK, ListDataSourcesResponse{DataSources: dataSources})
}

// GetDataSource handles getting a single data source by ID.
//
//	@Summary		Get data source
//	@Description	Get a single data source by ID for the authenticated user's organization
//	@Tags			data-sources
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Data Source ID"
//	@Success		200	{object}	DataSource
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/data-sources/{id} [get]
func (h *Handler) GetDataSource(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dataSourceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid data source ID")
		return
	}

	ds, err := h.service.GetDataSource(r.Context(), user.OrganizationID, dataSourceID)
	if err != nil {
		if errors.Is(err, ErrDataSourceNotFound) {
			respondError(w, http.StatusNotFound, "data source not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("get data source error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get data source")
		return
	}

	respondJSON(w, http.StatusOK, ds)
}

// CreateDataSource handles creating a new data source.
//
//	@Summary		Create data source
//	@Description	Create a new data source and return its API key (shown only once). Requires admin role.
//	@Tags			data-sources
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateDataSourceRequest	true	"Data source data"
//	@Success		201		{object}	CreateDataSourceResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		403		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/data-sources [post]
func (h *Handler) CreateDataSource(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateDataSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.service.CreateDataSource(r.Context(), user.OrganizationID, req)
	if err != nil {
		if errors.Is(err, ErrDataSourceNameEmpty) {
			respondError(w, http.StatusBadRequest, "data source name is required")
			return
		}
		log.Printf("create data source error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create data source")
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// DeleteDataSource handles deleting a data source.
//
//	@Summary		Delete data source
//	@Description	Delete a data source by ID. Requires admin role.
//	@Tags			data-sources
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Data Source ID"
//	@Success		200	{object}	MessageResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/data-sources/{id} [delete]
func (h *Handler) DeleteDataSource(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dataSourceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid data source ID")
		return
	}

	err = h.service.DeleteDataSource(r.Context(), user.OrganizationID, dataSourceID)
	if err != nil {
		if errors.Is(err, ErrDataSourceNotFound) {
			respondError(w, http.StatusNotFound, "data source not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("delete data source error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete data source")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "data source deleted"})
}

// RegenerateAPIKey handles regenerating the API key for a data source.
//
//	@Summary		Regenerate API key
//	@Description	Regenerate the API key for a data source (new key shown only once). Requires admin role.
//	@Tags			data-sources
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Data Source ID"
//	@Success		200	{object}	RegenerateKeyResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/data-sources/{id}/regenerate-key [post]
func (h *Handler) RegenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	dataSourceID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid data source ID")
		return
	}

	response, err := h.service.RegenerateAPIKey(r.Context(), user.OrganizationID, dataSourceID)
	if err != nil {
		if errors.Is(err, ErrDataSourceNotFound) {
			respondError(w, http.StatusNotFound, "data source not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("regenerate API key error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to regenerate API key")
		return
	}

	respondJSON(w, http.StatusOK, response)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}
