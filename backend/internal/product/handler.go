package product

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// Handler handles HTTP requests for products.
type Handler struct {
	service *Service
}

// NewHandler creates a new product handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListProducts handles listing products for the authenticated user's organization.
//
//	@Summary		List products
//	@Description	Get all products for the authenticated user's organization
//	@Tags			products
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	ListProductsResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/products [get]
func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	products, err := h.service.ListProducts(r.Context(), user.OrganizationID)
	if err != nil {
		log.Printf("list products error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to list products")
		return
	}

	respondJSON(w, http.StatusOK, ListProductsResponse{Products: products})
}

// GetProduct handles getting a single product by ID.
//
//	@Summary		Get product
//	@Description	Get a single product by ID for the authenticated user's organization
//	@Tags			products
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{object}	Product
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/products/{id} [get]
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	product, err := h.service.GetProduct(r.Context(), user.OrganizationID, productID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			respondError(w, http.StatusNotFound, "product not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("get product error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get product")
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// CreateProduct handles creating a new product.
//
//	@Summary		Create product
//	@Description	Create a new product and return its API key (shown only once)
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		CreateProductRequest	true	"Product data"
//	@Success		201		{object}	CreateProductResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/products [post]
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	response, err := h.service.CreateProduct(r.Context(), user.OrganizationID, req)
	if err != nil {
		if errors.Is(err, ErrProductNameEmpty) {
			respondError(w, http.StatusBadRequest, "product name is required")
			return
		}
		log.Printf("create product error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to create product")
		return
	}

	respondJSON(w, http.StatusCreated, response)
}

// DeleteProduct handles deleting a product.
//
//	@Summary		Delete product
//	@Description	Delete a product by ID
//	@Tags			products
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{object}	MessageResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/products/{id} [delete]
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	err = h.service.DeleteProduct(r.Context(), user.OrganizationID, productID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			respondError(w, http.StatusNotFound, "product not found")
			return
		}
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		log.Printf("delete product error: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to delete product")
		return
	}

	respondJSON(w, http.StatusOK, MessageResponse{Message: "product deleted"})
}

// RegenerateAPIKey handles regenerating the API key for a product.
//
//	@Summary		Regenerate API key
//	@Description	Regenerate the API key for a product (new key shown only once)
//	@Tags			products
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Product ID"
//	@Success		200	{object}	RegenerateKeyResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		403	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/products/{id}/regenerate-key [post]
func (h *Handler) RegenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	user := auth.UserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	response, err := h.service.RegenerateAPIKey(r.Context(), user.OrganizationID, productID)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			respondError(w, http.StatusNotFound, "product not found")
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
