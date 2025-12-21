package demo

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devbydaniel/litekpi/internal/auth"
	"github.com/devbydaniel/litekpi/internal/datasource"
)

// Handler handles HTTP requests for demo operations.
type Handler struct {
	service *Service
}

// NewHandler creates a new demo handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateDemoDataSource handles creating a demo data source with sample data.
//
//	@Summary		Create demo data source
//	@Description	Create a demo data source with sample measurements for the last month. Requires admin role.
//	@Tags			data-sources
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	datasource.CreateDataSourceResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/data-sources/demo [post]
func (h *Handler) CreateDemoDataSource(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	response, err := h.service.CreateDemoDataSource(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("create demo data source error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create demo data source")
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error string `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{Error: message})
}

// Ensure datasource.CreateDataSourceResponse is used in swagger docs
var _ = datasource.CreateDataSourceResponse{}
