package mcp

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

// MCPKeyContextKey is the context key for the authenticated MCP API key.
const MCPKeyContextKey contextKey = "mcpKey"

// APIKeyMiddleware creates a middleware that validates MCP API keys.
func APIKeyMiddleware(service *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				respondError(w, http.StatusUnauthorized, "unauthorized", "missing X-API-Key header")
				return
			}

			// Validate the API key
			key, err := service.ValidateKey(r.Context(), apiKey)
			if err != nil {
				if err == ErrKeyNotFound {
					respondError(w, http.StatusUnauthorized, "unauthorized", "invalid API key")
					return
				}
				respondError(w, http.StatusInternalServerError, "internal_error", "failed to validate API key")
				return
			}

			// Add MCP key to context
			ctx := context.WithValue(r.Context(), MCPKeyContextKey, key)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// MCPKeyFromContext retrieves the MCP API key from the request context.
func MCPKeyFromContext(ctx context.Context) *MCPAPIKey {
	key, _ := ctx.Value(MCPKeyContextKey).(*MCPAPIKey)
	return key
}

// OrgIDFromContext retrieves the organization ID from the MCP API key in context.
func OrgIDFromContext(ctx context.Context) uuid.UUID {
	key := MCPKeyFromContext(ctx)
	if key == nil {
		return uuid.Nil
	}
	return key.OrganizationID
}

func respondError(w http.ResponseWriter, status int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + errorType + `","message":"` + message + `"}`))
}
