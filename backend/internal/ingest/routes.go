package ingest

import (
	"github.com/go-chi/chi/v5"

	"github.com/devbydaniel/litekpi/internal/product"
)

// RegisterRoutes registers the ingest routes on the given router.
func (h *Handler) RegisterRoutes(r chi.Router, productRepo *product.Repository) {
	r.Route("/ingest", func(r chi.Router) {
		r.Use(APIKeyMiddleware(productRepo))
		r.Post("/", h.IngestSingle)
		r.Post("/batch", h.IngestBatch)
	})
}
