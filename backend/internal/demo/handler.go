package demo

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/devbydaniel/litekpi/internal/auth"
	"github.com/devbydaniel/litekpi/internal/product"
)

// Handler handles HTTP requests for demo operations.
type Handler struct {
	service *Service
}

// NewHandler creates a new demo handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateDemoProduct handles creating a demo product with sample data.
//
//	@Summary		Create demo product
//	@Description	Create a demo product with sample measurements for the last month
//	@Tags			products
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	product.CreateProductResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/products/demo [post]
func (h *Handler) CreateDemoProduct(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	response, err := h.service.CreateDemoProduct(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("create demo product error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create demo product")
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

// Ensure product.CreateProductResponse is used in swagger docs
var _ = product.CreateProductResponse{}
