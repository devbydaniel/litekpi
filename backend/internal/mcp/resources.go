package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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
func (r *ResourceRegistry) RegisterResources(server *mcp.Server, getOrgID func(ctx context.Context) uuid.UUID) {
	// Resource template: litekpi://measurements/{dataSourceId}/{measurementName}?timeframe={timeframe}
	server.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "litekpi://measurements/{dataSourceId}/{measurementName}",
			Name:        "Measurement Data",
			Description: "Raw measurement data points. Supports timeframe query parameter: last_7_days, last_30_days, this_month, last_month (default: last_30_days)",
			MIMEType:    "application/json",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			orgID := getOrgID(ctx)
			if orgID == uuid.Nil {
				return nil, fmt.Errorf("unauthorized: no organization context")
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

			// Verify data source belongs to organization
			ds, err := r.dsService.GetDataSource(ctx, orgID, dataSourceID)
			if err != nil {
				return nil, mcp.ResourceNotFoundError(req.Params.URI)
			}

			// Parse timeframe from query params
			timeframe := parsedURI.Query().Get("timeframe")
			if timeframe == "" {
				timeframe = "last_30_days"
			}

			start, end := getTimeframeRange(timeframe)

			// Fetch raw measurements (limit to 1000 points)
			measurements, err := r.ingestService.GetRawMeasurements(ctx, ds.ID, measurementName, start, end, nil, 1000)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch measurements: %w", err)
			}

			// Convert to output format
			dataPoints := make([]RawDataPoint, len(measurements))
			for i, m := range measurements {
				dataPoints[i] = RawDataPoint{
					Timestamp: m.Timestamp,
					Value:     m.Value,
					Metadata:  m.Metadata,
				}
			}

			// Serialize to JSON
			content, err := json.Marshal(dataPoints)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize data: %w", err)
			}

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI:      req.Params.URI,
					MIMEType: "application/json",
					Text:     string(content),
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
