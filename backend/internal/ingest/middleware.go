package ingest

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/devbydaniel/litekpi/internal/datasource"
)

type contextKey string

// DataSourceContextKey is the context key for the authenticated data source.
const DataSourceContextKey contextKey = "dataSource"

// APIKeyMiddleware creates a middleware that validates API keys.
func APIKeyMiddleware(dsRepo *datasource.Repository) func(http.Handler) http.Handler {
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

			// Look up data source by hash
			ds, err := dsRepo.GetDataSourceByAPIKeyHash(r.Context(), keyHash)
			if err != nil {
				respondError(w, http.StatusInternalServerError, "internal_error", "failed to validate API key")
				return
			}
			if ds == nil {
				respondError(w, http.StatusUnauthorized, "unauthorized", "invalid API key")
				return
			}

			// Add data source to context
			ctx := context.WithValue(r.Context(), DataSourceContextKey, ds)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// DataSourceFromContext retrieves the data source from the request context.
func DataSourceFromContext(ctx context.Context) *datasource.DataSource {
	ds, _ := ctx.Value(DataSourceContextKey).(*datasource.DataSource)
	return ds
}

func respondError(w http.ResponseWriter, status int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + errorType + `","message":"` + message + `"}`))
}
