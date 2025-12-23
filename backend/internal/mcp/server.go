package mcp

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/devbydaniel/litekpi/internal/datasource"
	"github.com/devbydaniel/litekpi/internal/ingest"
)

// ServerFactory creates MCP servers for authenticated requests.
type ServerFactory struct {
	toolRegistry     *ToolRegistry
	resourceRegistry *ResourceRegistry
}

// NewServerFactory creates a new MCP server factory.
func NewServerFactory(dsService *datasource.Service, ingestService *ingest.Service) *ServerFactory {
	return &ServerFactory{
		toolRegistry:     NewToolRegistry(dsService, ingestService),
		resourceRegistry: NewResourceRegistry(dsService, ingestService),
	}
}

// CreateServer creates a new MCP server configured with tools and resources.
// The getOrgID function is called to get the organization ID from context.
func (f *ServerFactory) CreateServer(getOrgID func(ctx context.Context) uuid.UUID) *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "litekpi",
			Version: "1.0.0",
		},
		nil,
	)

	// Register tools
	f.toolRegistry.RegisterTools(server, getOrgID)

	// Register resources
	f.resourceRegistry.RegisterResources(server, getOrgID)

	return server
}

// MCPHTTPHandler creates an HTTP handler for MCP protocol requests.
// It uses the StreamableHTTPHandler from the MCP SDK.
func (f *ServerFactory) MCPHTTPHandler() http.Handler {
	return mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			// Create a server for this request with org context from middleware
			return f.CreateServer(func(ctx context.Context) uuid.UUID {
				return OrgIDFromContext(r.Context())
			})
		},
		nil,
	)
}
