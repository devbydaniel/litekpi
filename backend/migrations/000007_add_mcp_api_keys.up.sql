-- MCP API Keys for organization-level access to MCP endpoints
CREATE TABLE mcp_api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    api_key_hash VARCHAR(255) NOT NULL UNIQUE,
    created_by UUID NOT NULL REFERENCES users(id),
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mcp_api_keys_organization_id ON mcp_api_keys(organization_id);
CREATE INDEX idx_mcp_api_keys_hash ON mcp_api_keys(api_key_hash);
