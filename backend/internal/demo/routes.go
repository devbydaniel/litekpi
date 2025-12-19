package demo

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all demo routes.
func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	// Register directly under /products/demo without creating a new /products group
	// (the /products group is already registered by the product module)
	r.With(authMiddleware).Post("/products/demo", h.CreateDemoProduct)
}
