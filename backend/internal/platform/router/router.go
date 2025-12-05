package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/trackable/trackable/internal/platform/database"
)

// New creates a new Chi router with middleware and routes configured.
func New(db *database.DB, appURL string) *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{appURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", healthHandler(db))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Placeholder for future routes
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusOK, map[string]string{
				"message": "Trackable API v1",
				"status":  "ok",
			})
		})
	})

	return r
}

// healthHandler returns a health check handler.
func healthHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Check database connection
		if err := db.Health(ctx); err != nil {
			respondJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status":   "unhealthy",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"status":   "healthy",
			"database": "connected",
		})
	}
}

// respondJSON writes a JSON response.
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
