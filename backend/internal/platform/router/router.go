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
	"github.com/devbydaniel/litekpi/internal/demo"
	"github.com/devbydaniel/litekpi/internal/ingest"
	"github.com/devbydaniel/litekpi/internal/platform/config"
	"github.com/devbydaniel/litekpi/internal/platform/database"
	"github.com/devbydaniel/litekpi/internal/product"

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
	emailService := auth.NewEmailService(auth.EmailConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		User:     cfg.SMTP.User,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
		AppURL:   cfg.AppURL,
	})
	authService := auth.NewService(authRepo, jwtService, emailService, cfg)
	authHandler := auth.NewHandler(authService)

	// Initialize product module
	productRepo := product.NewRepository(db.Pool)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	// Initialize ingest module
	ingestRepo := ingest.NewRepository(db.Pool)
	ingestService := ingest.NewService(ingestRepo)
	ingestHandler := ingest.NewHandler(ingestService, productService)

	// Initialize demo module (uses both product and ingest services)
	demoService := demo.NewService(productService, ingestService)
	demoHandler := demo.NewHandler(demoService)

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

		// Register product routes
		productHandler.RegisterRoutes(r, authService.Middleware)

		// Register demo routes (must be before ingest to handle /products/demo)
		demoHandler.RegisterRoutes(r, authService.Middleware)

		// Register ingest routes (uses API key auth, not JWT)
		ingestHandler.RegisterRoutes(r, productRepo)

		// Register measurement query routes (uses JWT auth)
		ingestHandler.RegisterMeasurementRoutes(r, authService.Middleware)
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
