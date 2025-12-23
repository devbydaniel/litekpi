package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/devbydaniel/litekpi/internal/auth"
	"github.com/devbydaniel/litekpi/internal/dashboard"
	"github.com/devbydaniel/litekpi/internal/datasource"
	"github.com/devbydaniel/litekpi/internal/demo"
	"github.com/devbydaniel/litekpi/internal/ingest"
	"github.com/devbydaniel/litekpi/internal/mcp"
	"github.com/devbydaniel/litekpi/internal/metric"
	"github.com/devbydaniel/litekpi/internal/platform/config"
	"github.com/devbydaniel/litekpi/internal/platform/database"
	"github.com/devbydaniel/litekpi/internal/platform/email"

	_ "github.com/devbydaniel/litekpi/docs" // Swagger docs
)

// New creates a new Chi router with middleware and routes configured.
func New(db *database.DB, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.AppURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize auth module
	authRepo := auth.NewRepository(db.Pool)
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	emailService := email.NewService(email.Config{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		User:     cfg.SMTP.User,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	})
	authEmailer := auth.NewAuthEmailer(emailService, cfg.AppURL)
	authService := auth.NewService(authRepo, jwtService, authEmailer, cfg)
	authHandler := auth.NewHandler(authService)

	// Initialize data source module
	dsRepo := datasource.NewRepository(db.Pool)
	dsService := datasource.NewService(dsRepo)
	dsHandler := datasource.NewHandler(dsService)

	// Initialize ingest module
	ingestRepo := ingest.NewRepository(db.Pool)
	ingestService := ingest.NewService(ingestRepo)
	ingestHandler := ingest.NewHandler(ingestService, dsService)

	// Initialize dashboard module
	dashboardRepo := dashboard.NewRepository(db.Pool)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	// Initialize metric module (unified metrics)
	metricRepo := metric.NewRepository(db.Pool)
	metricService := metric.NewService(metricRepo, dsService)
	metricHandler := metric.NewHandler(metricService, dashboardService)

	// Initialize demo module
	demoService := demo.NewService(dsService, ingestService)
	demoHandler := demo.NewHandler(demoService)

	// Initialize MCP module
	mcpRepo := mcp.NewRepository(db.Pool)
	mcpService := mcp.NewService(mcpRepo)
	mcpHandler := mcp.NewHandler(mcpService)
	mcpServerFactory := mcp.NewServerFactory(dsService, ingestService)

	// Health check endpoint
	r.Get("/health", healthHandler(db))

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// API status endpoint
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			respondJSON(w, http.StatusOK, map[string]string{
				"message": "LiteKPI API v1",
				"status":  "ok",
			})
		})

		// Register auth routes
		authHandler.RegisterRoutes(r, authService.Middleware)

		// Register data source routes
		dsHandler.RegisterRoutes(r, authService.Middleware)

		// Register dashboard routes
		dashboardHandler.RegisterRoutes(r, authService.Middleware)

		// Register metric routes (unified metrics)
		registerMetricRoutes(r, authService.Middleware, metricHandler)

		// Register demo routes
		demoHandler.RegisterRoutes(r, authService.Middleware)

		// Register ingest routes (uses API key auth, not JWT)
		ingestHandler.RegisterRoutes(r, dsRepo)

		// Register measurement query routes (uses JWT auth)
		ingestHandler.RegisterMeasurementRoutes(r, authService.Middleware)

		// Register MCP key management routes (uses JWT auth, admin only)
		mcpHandler.RegisterRoutes(r, authService.Middleware)

		// Register MCP protocol routes (uses MCP API key auth)
		mcpHandler.RegisterMCPProtocolRoutes(r, mcpServerFactory.MCPHTTPHandler())
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
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// registerMetricRoutes registers the unified metric routes.
func registerMetricRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler, h *metric.Handler) {
	r.Route("/dashboards/{id}/metrics", func(r chi.Router) {
		r.Use(authMiddleware)

		// Read operations
		r.Get("/", h.ListMetrics)
		r.Get("/compute", h.ComputeMetrics)

		// Write operations (editor and admin only)
		r.Group(func(r chi.Router) {
			r.Use(auth.EditorMiddleware)

			r.Post("/", h.CreateMetric)
			r.Put("/{metricId}", h.UpdateMetric)
			r.Delete("/{metricId}", h.DeleteMetric)
			r.Put("/reorder", h.ReorderMetrics)
		})
	})
}
