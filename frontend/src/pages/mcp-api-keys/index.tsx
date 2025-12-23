import { Plus } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import { useAuthStore } from '@/shared/stores/auth-store'
import { useMCPApiKeys } from './hooks/use-mcp-api-keys'
import { MCPApiKeyList } from './ui/mcp-api-key-list'
import { CreateMCPApiKeyDialog } from './ui/create-mcp-api-key-dialog'
import { EditMCPApiKeyDialog } from './ui/edit-mcp-api-key-dialog'
import { DeleteMCPApiKeyDialog } from './ui/delete-mcp-api-key-dialog'

export function MCPApiKeysPage() {
  const { user } = useAuthStore()
  const isAdmin = user?.role === 'admin'

  const {
    keys,
    dataSources,
    isLoading,
    createDialogOpen,
    setCreateDialogOpen,
    editDialogOpen,
    keyToEdit,
    deleteDialogOpen,
    keyToDelete,
    apiKey,
    isCreating,
    isUpdating,
    isDeleting,
    handleCreateKey,
    handleEditKey,
    handleUpdateKey,
    handleDeleteKey,
    confirmDeleteKey,
    closeDialogs,
  } = useMCPApiKeys()

  // Non-admins see restricted access message
  if (!isAdmin) {
    return (
      <AuthenticatedLayout title="MCP API Keys">
        <div className="text-center text-muted-foreground">
          Admin access required to manage MCP API keys.
        </div>
      </AuthenticatedLayout>
    )
  }

  return (
    <AuthenticatedLayout
      title="MCP API Keys"
      actions={
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Plus className="h-4 w-4" />
          New API Key
        </Button>
      }
    >
      <MCPApiKeyList
        keys={keys}
        isLoading={isLoading}
        onEdit={handleEditKey}
        onDelete={handleDeleteKey}
      />

      <CreateMCPApiKeyDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onCreate={handleCreateKey}
        dataSources={dataSources}
        apiKey={apiKey}
        isLoading={isCreating}
        onClose={closeDialogs}
      />

      <EditMCPApiKeyDialog
        open={editDialogOpen}
        apiKey={keyToEdit}
        dataSources={dataSources}
        isLoading={isUpdating}
        onSave={handleUpdateKey}
        onClose={closeDialogs}
      />

      <DeleteMCPApiKeyDialog
        open={deleteDialogOpen}
        apiKey={keyToDelete}
        isLoading={isDeleting}
        onConfirm={confirmDeleteKey}
        onClose={closeDialogs}
      />
    </AuthenticatedLayout>
  )
}
