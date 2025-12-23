import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from '@/shared/components/ui/form'
import { DataSourceMultiSelect } from './data-source-multi-select'
import type { DataSource, MCPAPIKey } from '@/shared/api/generated/models'

const editMCPApiKeySchema = z.object({
  dataSourceIds: z
    .array(z.string())
    .min(1, 'At least one data source must be selected'),
})

type EditMCPApiKeyFormValues = z.infer<typeof editMCPApiKeySchema>

interface EditMCPApiKeyDialogProps {
  open: boolean
  apiKey: MCPAPIKey | null
  dataSources: DataSource[]
  isLoading: boolean
  onSave: (dataSourceIds: string[]) => Promise<void>
  onClose: () => void
}

export function EditMCPApiKeyDialog({
  open,
  apiKey,
  dataSources,
  isLoading,
  onSave,
  onClose,
}: EditMCPApiKeyDialogProps) {
  const form = useForm<EditMCPApiKeyFormValues>({
    resolver: zodResolver(editMCPApiKeySchema),
    defaultValues: {
      dataSourceIds: [],
    },
  })

  // Reset form when dialog opens with new key
  useEffect(() => {
    if (open && apiKey) {
      form.reset({
        dataSourceIds: apiKey.allowedDataSourceIds ?? [],
      })
    }
  }, [open, apiKey, form])

  const handleSubmit = async (values: EditMCPApiKeyFormValues) => {
    await onSave(values.dataSourceIds)
  }

  const handleClose = () => {
    form.reset()
    onClose()
  }

  return (
    <Dialog open={open} onOpenChange={(open) => !open && handleClose()}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Edit MCP API Key</DialogTitle>
          <DialogDescription>
            Update the data sources for{' '}
            <span className="font-medium">{apiKey?.name}</span>.
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
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
                onClick={handleClose}
                disabled={isLoading}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading ? 'Saving...' : 'Save'}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
