import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import type { DataSource } from '@/shared/api/generated/models'

interface DeleteDataSourceDialogProps {
  open: boolean
  dataSource: DataSource | null
  isLoading: boolean
  onConfirm: () => Promise<void>
  onClose: () => void
}

export function DeleteDataSourceDialog({
  open,
  dataSource,
  isLoading,
  onConfirm,
  onClose,
}: DeleteDataSourceDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Data Source</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete{' '}
            <span className="font-medium">{dataSource?.name}</span>? This will
            permanently delete all measurements associated with this data source.
            This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
