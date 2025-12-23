package mcp

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/devbydaniel/litekpi/internal/auth"
)

// RegisterRoutes registers all MCP routes.
func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(next http.Handler) http.Handler) {
	// Key management routes (JWT auth, admin only)
	r.Route("/mcp/keys", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(auth.AdminMiddleware)

		r.Post("/", h.CreateKey)
		r.Get("/", h.ListKeys)
		r.Put("/{id}", h.UpdateKey)
		r.Delete("/{id}", h.DeleteKey)
	})
}

// RegisterMCPProtocolRoutes registers the MCP protocol endpoint (API key auth).
func (h *Handler) RegisterMCPProtocolRoutes(r chi.Router, mcpHandler http.Handler) {
	r.Route("/mcp", func(r chi.Router) {
		r.Use(APIKeyMiddleware(h.service))
		r.Handle("/*", mcpHandler)
	})
}
