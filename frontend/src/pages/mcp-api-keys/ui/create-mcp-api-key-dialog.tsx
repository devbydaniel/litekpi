import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from '@/shared/components/ui/form'
import { useCreateMCPApiKeyForm } from '../hooks/use-create-mcp-api-key-form'
import { DataSourceMultiSelect } from './data-source-multi-select'
import { ApiKeyDisplay } from '@/pages/data-sources/ui/api-key-display'
import type { DataSource } from '@/shared/api/generated/models'

interface CreateMCPApiKeyDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onCreate: (name: string, dataSourceIds: string[]) => Promise<void>
  dataSources: DataSource[]
  apiKey: string | null
  isLoading: boolean
  onClose: () => void
}

export function CreateMCPApiKeyDialog({
  open,
  onOpenChange,
  onCreate,
  dataSources,
  apiKey,
  isLoading,
  onClose,
}: CreateMCPApiKeyDialogProps) {
  const { form, reset } = useCreateMCPApiKeyForm()

  const handleSubmit = async (values: {
    name: string
    dataSourceIds: string[]
  }) => {
    await onCreate(values.name, values.dataSourceIds)
  }

  const handleClose = () => {
    reset()
    onClose()
  }

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      handleClose()
    } else {
      onOpenChange(open)
    }
  }

  // Show API key success state
  if (apiKey) {
    return (
      <Dialog open={open} onOpenChange={handleOpenChange}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>MCP API Key Created</DialogTitle>
            <DialogDescription>
              Copy your API key now. You will not be able to see it again.
            </DialogDescription>
          </DialogHeader>

          <ApiKeyDisplay apiKey={apiKey} />

          <DialogFooter>
            <Button onClick={handleClose}>Done</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    )
  }

  // Show create form
  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Create MCP API Key</DialogTitle>
          <DialogDescription>
            Create an API key for MCP access. Select which data sources this key
            can access.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Key Name</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="Production MCP Key"
                      disabled={isLoading}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="dataSourceIds"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Data Sources</FormLabel>
                  <FormControl>
                    <DataSourceMultiSelect
                      dataSources={dataSources}
                      selectedIds={field.value}
                      onChange={field.onChange}
                      disabled={isLoading}
                    />
                  </FormControl>
                  <FormDescription>
                    This API key will only have access to the selected data
                    sources.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => handleOpenChange(false)}
                disabled={isLoading}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Creating...' : 'Create'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
