-- Join table for MCP API keys to data sources (many-to-many)
CREATE TABLE mcp_api_key_data_sources (
    mcp_api_key_id UUID NOT NULL REFERENCES mcp_api_keys(id) ON DELETE CASCADE,
    data_source_id UUID NOT NULL REFERENCES data_sources(id) ON DELETE CASCADE,
    PRIMARY KEY (mcp_api_key_id, data_source_id)
);

-- Index for querying allowed data sources by key
CREATE INDEX idx_mcp_api_key_data_sources_key_id ON mcp_api_key_data_sources(mcp_api_key_id);

-- Index for checking if a data source is in use by any MCP key
CREATE INDEX idx_mcp_api_key_data_sources_ds_id ON mcp_api_key_data_sources(data_source_id);
