import { Plus, Sparkles } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import { useAuth } from '@/shared/hooks/use-auth'
import { useDataSources } from './hooks/use-data-sources'
import { DataSourceList } from './ui/data-source-list'
import { CreateDataSourceDialog } from './ui/create-data-source-dialog'
import { RegenerateKeyDialog } from './ui/regenerate-key-dialog'
import { DeleteDataSourceDialog } from './ui/delete-data-source-dialog'

export function DataSourcesPage() {
  const { user } = useAuth()
  const isAdmin = user?.role === 'admin'

  const {
    dataSources,
    isLoading,
    createDialogOpen,
    setCreateDialogOpen,
    regenerateDialogOpen,
    deleteDialogOpen,
    selectedDataSource,
    dataSourceToDelete,
    apiKey,
    isCreating,
    isCreatingDemo,
    isRegenerating,
    isDeleting,
    handleCreateDataSource,
    handleCreateDemo,
    handleDeleteDataSource,
    confirmDeleteDataSource,
    handleRegenerateKey,
    confirmRegenerateKey,
    closeDialogs,
  } = useDataSources()

  return (
    <AuthenticatedLayout
      title="Data Sources"
      actions={
        isAdmin ? (
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={handleCreateDemo}
              disabled={isCreatingDemo}
            >
              <Sparkles className="h-4 w-4" />
              Create Demo
            </Button>
            <Button onClick={() => setCreateDialogOpen(true)}>
              <Plus className="h-4 w-4" />
              New Data Source
            </Button>
          </div>
        ) : undefined
      }
    >
      <DataSourceList
        dataSources={dataSources}
        isLoading={isLoading}
        onDelete={handleDeleteDataSource}
        onRegenerateKey={handleRegenerateKey}
      />

      <CreateDataSourceDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onCreate={handleCreateDataSource}
        apiKey={apiKey}
        isLoading={isCreating}
        onClose={closeDialogs}
      />

      <RegenerateKeyDialog
        open={regenerateDialogOpen}
        dataSource={selectedDataSource}
        apiKey={apiKey}
        isLoading={isRegenerating}
        onConfirm={confirmRegenerateKey}
        onClose={closeDialogs}
      />

      <DeleteDataSourceDialog
        open={deleteDialogOpen}
        dataSource={dataSourceToDelete}
        isLoading={isDeleting}
        onConfirm={confirmDeleteDataSource}
        onClose={closeDialogs}
      />
    </AuthenticatedLayout>
  )
}
