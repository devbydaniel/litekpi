package ingest

import (
	"net/http"

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

// RegisterMeasurementRoutes registers measurement query routes on the given router.
func (h *Handler) RegisterMeasurementRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	r.Route("/products/{productId}/measurements", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/", h.ListMeasurementNames)
		r.Get("/{name}/metadata", h.GetMetadataValues)
		r.Get("/{name}/data", h.GetMeasurementData)
		r.Get("/{name}/data/split", h.GetMeasurementDataSplit)
		r.Get("/{name}/preferences", h.GetPreferences)
		r.Post("/{name}/preferences", h.SavePreferences)
	})
}
