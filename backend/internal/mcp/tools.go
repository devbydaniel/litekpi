package mcp

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/devbydaniel/litekpi/internal/datasource"
	"github.com/devbydaniel/litekpi/internal/ingest"
)

// DataSourceOutput represents a data source in tool output.
type DataSourceOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListDataSourcesOutput is the output for list_data_sources tool.
type ListDataSourcesOutput struct {
	DataSources []DataSourceOutput `json:"dataSources"`
}

// ListMeasurementsInput is the input for list_measurements tool.
type ListMeasurementsInput struct {
	DataSourceID string `json:"dataSourceId"`
}

// MeasurementOutput represents a measurement in tool output.
type MeasurementOutput struct {
	Name         string   `json:"name"`
	MetadataKeys []string `json:"metadataKeys"`
}

// ListMeasurementsOutput is the output for list_measurements tool.
type ListMeasurementsOutput struct {
	Measurements []MeasurementOutput `json:"measurements"`
}

// ToolRegistry holds references to services needed by MCP tools.
type ToolRegistry struct {
	dsService     *datasource.Service
	ingestService *ingest.Service
}

// NewToolRegistry creates a new tool registry.
func NewToolRegistry(dsService *datasource.Service, ingestService *ingest.Service) *ToolRegistry {
	return &ToolRegistry{
		dsService:     dsService,
		ingestService: ingestService,
	}
}

// RegisterTools registers all MCP tools with the server.
func (t *ToolRegistry) RegisterTools(server *mcp.Server, getMCPKey func(ctx context.Context) *MCPAPIKey) {
	// Tool: list_data_sources
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_data_sources",
		Description: "List all data sources accessible by this API key",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, ListDataSourcesOutput, error) {
		mcpKey := getMCPKey(ctx)
		if mcpKey == nil {
			return nil, ListDataSourcesOutput{}, fmt.Errorf("unauthorized: no API key context")
		}

		// Return only allowed data sources for this key
		output := ListDataSourcesOutput{
			DataSources: make([]DataSourceOutput, 0, len(mcpKey.AllowedDataSourceIDs)),
		}

		for _, dsID := range mcpKey.AllowedDataSourceIDs {
			ds, err := t.dsService.GetDataSource(ctx, mcpKey.OrganizationID, dsID)
			if err != nil {
				continue // Skip if data source no longer exists
			}
			output.DataSources = append(output.DataSources, DataSourceOutput{
				ID:   ds.ID.String(),
				Name: ds.Name,
			})
		}

		return nil, output, nil
	})

	// Tool: list_measurements
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_measurements",
		Description: "List available measurement names and their metadata keys for a data source",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListMeasurementsInput) (*mcp.CallToolResult, ListMeasurementsOutput, error) {
		mcpKey := getMCPKey(ctx)
		if mcpKey == nil {
			return nil, ListMeasurementsOutput{}, fmt.Errorf("unauthorized: no API key context")
		}

		dsID, err := uuid.Parse(input.DataSourceID)
		if err != nil {
			return nil, ListMeasurementsOutput{}, fmt.Errorf("invalid dataSourceId: %w", err)
		}

		// Check if this key has access to this data source
		if !hasAccess(mcpKey, dsID) {
			return nil, ListMeasurementsOutput{}, fmt.Errorf("unauthorized: API key does not have access to this data source")
		}

		// Verify data source belongs to organization
		ds, err := t.dsService.GetDataSource(ctx, mcpKey.OrganizationID, dsID)
		if err != nil {
			return nil, ListMeasurementsOutput{}, fmt.Errorf("data source not found or unauthorized: %w", err)
		}

		measurements, err := t.ingestService.GetMeasurementNames(ctx, ds.ID)
		if err != nil {
			return nil, ListMeasurementsOutput{}, fmt.Errorf("failed to list measurements: %w", err)
		}

		output := ListMeasurementsOutput{
			Measurements: make([]MeasurementOutput, len(measurements)),
		}
		for i, m := range measurements {
			output.Measurements[i] = MeasurementOutput{
				Name:         m.Name,
				MetadataKeys: m.MetadataKeys,
			}
		}

		return nil, output, nil
	})
}

// hasAccess checks if the MCP API key has access to the given data source.
func hasAccess(key *MCPAPIKey, dsID uuid.UUID) bool {
	for _, allowedID := range key.AllowedDataSourceIDs {
		if allowedID == dsID {
			return true
		}
	}
	return false
}
