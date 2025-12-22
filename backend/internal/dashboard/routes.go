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
		r.Get("/{id}/scalar-metrics", h.ListScalarMetrics)
		r.Get("/{id}/scalar-metrics/compute", h.ComputeScalarMetrics)

		// Write operations (editor and admin only)
		r.Group(func(r chi.Router) {
			r.Use(auth.EditorMiddleware)

			r.Post("/", h.CreateDashboard)
			r.Put("/{id}", h.UpdateDashboard)
			r.Delete("/{id}", h.DeleteDashboard)

			// Time Series routes
			r.Post("/{id}/time-series", h.CreateTimeSeries)
			r.Put("/{id}/time-series/{timeSeriesId}", h.UpdateTimeSeries)
			r.Delete("/{id}/time-series/{timeSeriesId}", h.DeleteTimeSeries)
			r.Put("/{id}/time-series/reorder", h.ReorderTimeSeries)

			// Scalar Metric routes
			r.Post("/{id}/scalar-metrics", h.CreateScalarMetric)
			r.Put("/{id}/scalar-metrics/{scalarMetricId}", h.UpdateScalarMetric)
			r.Delete("/{id}/scalar-metrics/{scalarMetricId}", h.DeleteScalarMetric)
			r.Put("/{id}/scalar-metrics/reorder", h.ReorderScalarMetrics)
		})
	})
}
