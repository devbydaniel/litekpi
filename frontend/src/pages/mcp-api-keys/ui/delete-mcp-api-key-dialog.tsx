import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import type { MCPAPIKey } from '@/shared/api/generated/models'

interface DeleteMCPApiKeyDialogProps {
  open: boolean
  apiKey: MCPAPIKey | null
  isLoading: boolean
  onConfirm: () => Promise<void>
  onClose: () => void
}

export function DeleteMCPApiKeyDialog({
  open,
  apiKey,
  isLoading,
  onConfirm,
  onClose,
}: DeleteMCPApiKeyDialogProps) {
  return (
    <Dialog open={open} onOpenChange={(open) => !open && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete MCP API Key</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete{' '}
            <span className="font-medium">{apiKey?.name}</span>? Any
            applications using this key will lose access immediately. This
            action cannot be undone.
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
