import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import type { Product } from '@/shared/api/generated/models'
import { ApiKeyDisplay } from './api-key-display'

interface RegenerateKeyDialogProps {
  open: boolean
  product: Product | null
  apiKey: string | null
  isLoading: boolean
  onConfirm: () => Promise<void>
  onClose: () => void
}

export function RegenerateKeyDialog({
  open,
  product,
  apiKey,
  isLoading,
  onConfirm,
  onClose,
}: RegenerateKeyDialogProps) {
  const handleClose = () => {
    onClose()
  }

  // Show new API key
  if (apiKey) {
    return (
      <Dialog open={open} onOpenChange={handleClose}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>API Key Regenerated</DialogTitle>
            <DialogDescription>
              Copy your new API key now. You won't be able to see it again.
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

  // Show confirmation
  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Regenerate API Key</DialogTitle>
          <DialogDescription>
            Are you sure you want to regenerate the API key for{' '}
            <span className="font-medium">{product?.name}</span>? The current
            key will be invalidated immediately.
          </DialogDescription>
        </DialogHeader>

        <DialogFooter>
          <Button variant="outline" onClick={handleClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Regenerating...' : 'Regenerate'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
