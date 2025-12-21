package datasource

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// RegisterRoutes registers all data source routes.
func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	r.Route("/data-sources", func(r chi.Router) {
		r.Use(authMiddleware)

		// Read operations (all authenticated users)
		r.Get("/", h.ListDataSources)
		r.Get("/{id}", h.GetDataSource)

		// Write operations (admin only)
		r.Group(func(r chi.Router) {
			r.Use(auth.AdminMiddleware)
			r.Post("/", h.CreateDataSource)
			r.Delete("/{id}", h.DeleteDataSource)
			r.Post("/{id}/regenerate-key", h.RegenerateAPIKey)
		})
	})
}
