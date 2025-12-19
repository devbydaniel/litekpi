package product

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all product routes.
func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	r.Route("/products", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", h.ListProducts)
		r.Post("/", h.CreateProduct)
		r.Delete("/{id}", h.DeleteProduct)
		r.Post("/{id}/regenerate-key", h.RegenerateAPIKey)
	})
}
