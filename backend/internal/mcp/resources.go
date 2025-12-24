package mcp

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/devbydaniel/litekpi/internal/datasource"
	"github.com/devbydaniel/litekpi/internal/ingest"
)

// RawDataPoint represents a measurement data point for MCP resource output.
type RawDataPoint struct {
	Timestamp time.Time         `json:"timestamp"`
	Value     float64           `json:"value"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// ResourceRegistry holds references to services needed by MCP resources.
type ResourceRegistry struct {
	dsService     *datasource.Service
	ingestService *ingest.Service
}

// NewResourceRegistry creates a new resource registry.
func NewResourceRegistry(dsService *datasource.Service, ingestService *ingest.Service) *ResourceRegistry {
	return &ResourceRegistry{
		dsService:     dsService,
		ingestService: ingestService,
	}
}

// RegisterResources registers all MCP resources with the server.
func (r *ResourceRegistry) RegisterResources(server *mcp.Server, getMCPKey func(ctx context.Context) *MCPAPIKey) {
	// Resource template: litekpi://measurements/{dataSourceId}/{measurementName}{?timeframe,metadata}
	server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "litekpi://measurements/{dataSourceId}/{measurementName}{?timeframe,metadata}",
			Name:        "measurement_data",
			Description: `Raw measurement data points as CSV. Query parameters:
- timeframe: last_7_days, last_30_days, this_month, last_month (default: last_30_days)
- metadata: URL-encoded JSON object to filter by metadata, e.g. %7B%22region%22%3A%22eu%22%7D for {"region":"eu"}`,
			MIMEType: "text/csv",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			mcpKey := getMCPKey(ctx)
			if mcpKey == nil {
				return nil, fmt.Errorf("unauthorized: no API key context")
			}

			// Parse URI: litekpi://measurements/{dataSourceId}/{measurementName}?timeframe=...
			parsedURI, err := url.Parse(req.Params.URI)
			if err != nil {
				return nil, fmt.Errorf("invalid URI: %w", err)
			}

			// Extract path parts: /measurements/{dataSourceId}/{measurementName}
			pathParts := strings.Split(strings.Trim(parsedURI.Path, "/"), "/")
			if len(pathParts) < 2 {
				return nil, fmt.Errorf("invalid URI path: expected /measurements/{dataSourceId}/{measurementName}")
			}

			dataSourceIDStr := pathParts[0]
			measurementName := pathParts[1]

			dataSourceID, err := uuid.Parse(dataSourceIDStr)
			if err != nil {
				return nil, fmt.Errorf("invalid dataSourceId: %w", err)
			}

			// Check if this key has access to this data source
			if !hasAccess(mcpKey, dataSourceID) {
				return nil, mcp.ResourceNotFoundError(req.Params.URI)
			}

			// Verify data source belongs to organization
			ds, err := r.dsService.GetDataSource(ctx, mcpKey.OrganizationID, dataSourceID)
			if err != nil {
				return nil, mcp.ResourceNotFoundError(req.Params.URI)
			}

			// Parse timeframe from query params
			query := parsedURI.Query()
			timeframe := query.Get("timeframe")
			if timeframe == "" {
				timeframe = "last_30_days"
			}

			// Parse metadata filter from query params (JSON-encoded)
			var metadataFilter map[string]string
			if metadataJSON := query.Get("metadata"); metadataJSON != "" {
				if err := json.Unmarshal([]byte(metadataJSON), &metadataFilter); err != nil {
					return nil, fmt.Errorf("invalid metadata JSON: %w", err)
				}
			}

			start, end := getTimeframeRange(timeframe)

			// Fetch raw measurements (limit to 1000 points)
			measurements, err := r.ingestService.GetRawMeasurements(ctx, ds.ID, measurementName, start, end, metadataFilter, 1000)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch measurements: %w", err)
			}

			// Collect all unique metadata keys
			metaKeys := make(map[string]struct{})
			for _, m := range measurements {
				for k := range m.Metadata {
					metaKeys[k] = struct{}{}
				}
			}

			// Sort metadata keys for consistent column order
			sortedMetaKeys := make([]string, 0, len(metaKeys))
			for k := range metaKeys {
				sortedMetaKeys = append(sortedMetaKeys, k)
			}
			sort.Strings(sortedMetaKeys)

			// Build CSV string
			var buf bytes.Buffer
			w := csv.NewWriter(&buf)

			// Write header: timestamp, value, [metadata columns...]
			header := append([]string{"timestamp", "value"}, sortedMetaKeys...)
			w.Write(header)

			// Write data rows
			for _, m := range measurements {
				row := []string{
					m.Timestamp.Format(time.RFC3339),
					fmt.Sprintf("%g", m.Value),
				}
				for _, k := range sortedMetaKeys {
					row = append(row, m.Metadata[k])
				}
				w.Write(row)
			}
			w.Flush()

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI:      req.Params.URI,
					MIMEType: "text/csv",
					Text:     buf.String(),
				}},
			}, nil
		},
	)
}

// getTimeframeRange calculates the start and end times for a given timeframe.
func getTimeframeRange(timeframe string) (time.Time, time.Time) {
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	switch timeframe {
	case "last_7_days":
		return today.AddDate(0, 0, -7), today.AddDate(0, 0, 1)
	case "last_30_days":
		return today.AddDate(0, 0, -30), today.AddDate(0, 0, 1)
	case "this_month":
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		return start, today.AddDate(0, 0, 1)
	case "last_month":
		firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		return firstOfThisMonth.AddDate(0, -1, 0), firstOfThisMonth
	default:
		// Default to last 30 days
		return today.AddDate(0, 0, -30), today.AddDate(0, 0, 1)
	}
}
