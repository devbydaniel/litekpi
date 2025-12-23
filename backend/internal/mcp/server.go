package mcp

import (
	"context"
	"net/http"

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
// The getMCPKey function is called to get the MCP API key from context.
func (f *ServerFactory) CreateServer(getMCPKey func(ctx context.Context) *MCPAPIKey) *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "litekpi",
			Version: "1.0.0",
		},
		nil,
	)

	// Register tools
	f.toolRegistry.RegisterTools(server, getMCPKey)

	// Register resources
	f.resourceRegistry.RegisterResources(server, getMCPKey)

	return server
}

// MCPHTTPHandler creates an HTTP handler for MCP protocol requests.
// It uses the StreamableHTTPHandler from the MCP SDK.
func (f *ServerFactory) MCPHTTPHandler() http.Handler {
	return mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			// Create a server for this request with MCP key context from middleware
			return f.CreateServer(func(ctx context.Context) *MCPAPIKey {
				return MCPKeyFromContext(r.Context())
			})
		},
		nil,
	)
}
