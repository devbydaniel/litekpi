package report

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// RegisterRoutes registers all report routes.
func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	r.Route("/reports", func(r chi.Router) {
		r.Use(authMiddleware)

		// Read operations (all authenticated users)
		r.Get("/", h.ListReports)
		r.Get("/{id}", h.GetReport)
		r.Get("/{id}/compute", h.ComputeReport)
		r.Get("/{id}/kpis", h.ListKPIs)

		// Write operations (editor and admin only)
		r.Group(func(r chi.Router) {
			r.Use(auth.EditorMiddleware)

			r.Post("/", h.CreateReport)
			r.Put("/{id}", h.UpdateReport)
			r.Delete("/{id}", h.DeleteReport)

			// KPI routes
			r.Post("/{id}/kpis", h.CreateKPI)
			r.Put("/{id}/kpis/{kpiId}", h.UpdateKPI)
			r.Delete("/{id}/kpis/{kpiId}", h.DeleteKPI)
			r.Put("/{id}/kpis/reorder", h.ReorderKPIs)
		})
	})
}
