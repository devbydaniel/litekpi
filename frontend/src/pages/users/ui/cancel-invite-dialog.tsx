import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import type { InviteWithInviter } from '@/shared/api/generated/models'

interface CancelInviteDialogProps {
  open: boolean
  invite: InviteWithInviter | null
  isLoading: boolean
  onConfirm: () => Promise<void>
  onClose: () => void
}

export function CancelInviteDialog({
  open,
  invite,
  isLoading,
  onConfirm,
  onClose,
}: CancelInviteDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Cancel Invite</DialogTitle>
          <DialogDescription>
            Are you sure you want to cancel the invite for{' '}
            <span className="font-medium">{invite?.email}</span>? The invitation
            link will no longer work.
          </DialogDescription>
        </DialogHeader>

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={isLoading}>
            Keep Invite
          </Button>
          <Button variant="destructive" onClick={onConfirm} disabled={isLoading}>
            {isLoading ? 'Cancelling...' : 'Cancel Invite'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
