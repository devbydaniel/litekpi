package ingest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/devbydaniel/litekpi/internal/product"
)

type contextKey string

// ProductContextKey is the context key for the authenticated product.
const ProductContextKey contextKey = "product"

// APIKeyMiddleware creates a middleware that validates API keys.
func APIKeyMiddleware(productRepo *product.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				respondError(w, http.StatusUnauthorized, "unauthorized", "missing X-API-Key header")
				return
			}

			// Hash the API key
			hashBytes := sha256.Sum256([]byte(apiKey))
			keyHash := hex.EncodeToString(hashBytes[:])

			// Look up product by hash
			prod, err := productRepo.GetProductByAPIKeyHash(r.Context(), keyHash)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "internal_error", "failed to validate API key")
				return
			}
			if prod == nil {
				respondError(w, http.StatusUnauthorized, "unauthorized", "invalid API key")
				return
			}

			// Add product to context
			ctx := context.WithValue(r.Context(), ProductContextKey, prod)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ProductFromContext retrieves the product from the request context.
func ProductFromContext(ctx context.Context) *product.Product {
	prod, _ := ctx.Value(ProductContextKey).(*product.Product)
	return prod
}

func respondError(w http.ResponseWriter, status int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + errorType + `","message":"` + message + `"}`))
}
