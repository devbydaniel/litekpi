import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import type { User } from '@/shared/api/generated/models'

interface RemoveUserDialogProps {
  open: boolean
  user: User | null
  isLoading: boolean
  onConfirm: () => Promise<void>
  onClose: () => void
}

export function RemoveUserDialog({
  open,
  user,
  isLoading,
  onConfirm,
  onClose,
}: RemoveUserDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Remove User</DialogTitle>
          <DialogDescription>
            Are you sure you want to remove{' '}
            <span className="font-medium">{user?.name}</span> from your
            organization? They will lose access immediately.
          </DialogDescription>
        </DialogHeader>

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Removing...' : 'Remove'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
