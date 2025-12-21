package dashboard

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// RegisterRoutes registers all dashboard routes.
func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	r.Route("/dashboards", func(r chi.Router) {
		r.Use(authMiddleware)

		// Read operations (all authenticated users)
		r.Get("/", h.ListDashboards)
		r.Get("/default", h.GetDefaultDashboard)
		r.Get("/{id}", h.GetDashboard)

		// Write operations (editor and admin only)
		r.Group(func(r chi.Router) {
			r.Use(auth.EditorMiddleware)

			r.Post("/", h.CreateDashboard)
			r.Put("/{id}", h.UpdateDashboard)
			r.Delete("/{id}", h.DeleteDashboard)

			// Widget routes
			r.Post("/{id}/widgets", h.CreateWidget)
			r.Put("/{id}/widgets/{widgetId}", h.UpdateWidget)
			r.Delete("/{id}/widgets/{widgetId}", h.DeleteWidget)
			r.Put("/{id}/widgets/reorder", h.ReorderWidgets)
		})
	})
}
