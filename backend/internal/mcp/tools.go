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
	DataSourceID string `json:"dataSourceId" jsonschema:"description=The ID of the data source to query"`
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
func (t *ToolRegistry) RegisterTools(server *mcp.Server, getOrgID func(ctx context.Context) uuid.UUID) {
	// Tool: list_data_sources
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_data_sources",
		Description: "List all data sources for the organization",
	}, func(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, ListDataSourcesOutput, error) {
		orgID := getOrgID(ctx)
		if orgID == uuid.Nil {
			return nil, ListDataSourcesOutput{}, fmt.Errorf("unauthorized: no organization context")
		}

		dataSources, err := t.dsService.ListDataSources(ctx, orgID)
		if err != nil {
			return nil, ListDataSourcesOutput{}, fmt.Errorf("failed to list data sources: %w", err)
		}

		output := ListDataSourcesOutput{
			DataSources: make([]DataSourceOutput, len(dataSources)),
		}
		for i, ds := range dataSources {
			output.DataSources[i] = DataSourceOutput{
				ID:   ds.ID.String(),
				Name: ds.Name,
			}
		}

		return nil, output, nil
	})

	// Tool: list_measurements
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_measurements",
		Description: "List available measurement names and their metadata keys for a data source",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListMeasurementsInput) (*mcp.CallToolResult, ListMeasurementsOutput, error) {
		orgID := getOrgID(ctx)
		if orgID == uuid.Nil {
			return nil, ListMeasurementsOutput{}, fmt.Errorf("unauthorized: no organization context")
		}

		dsID, err := uuid.Parse(input.DataSourceID)
		if err != nil {
			return nil, ListMeasurementsOutput{}, fmt.Errorf("invalid dataSourceId: %w", err)
		}

		// Verify data source belongs to organization
		ds, err := t.dsService.GetDataSource(ctx, orgID, dsID)
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
