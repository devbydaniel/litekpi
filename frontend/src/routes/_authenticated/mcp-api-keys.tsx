import { createFileRoute } from '@tanstack/react-router'
import { MCPApiKeysPage } from '@/pages/mcp-api-keys'

export const Route = createFileRoute('/_authenticated/mcp-api-keys')({
  component: MCPApiKeysPage,
})
