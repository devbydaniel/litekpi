package demo

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// RegisterRoutes registers all demo routes.
func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	// Register demo data source creation (admin only)
	r.With(authMiddleware, auth.AdminMiddleware).Post("/data-sources/demo", h.CreateDemoDataSource)
}
