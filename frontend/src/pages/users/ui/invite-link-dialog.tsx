import { useState } from 'react'
import { Copy, Check } from 'lucide-react'
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

interface InviteLinkDialogProps {
  open: boolean
  inviteLink: string | null
  onClose: () => void
}

export function InviteLinkDialog({
  open,
  inviteLink,
  onClose,
}: InviteLinkDialogProps) {
  const [copied, setCopied] = useState(false)

  const handleCopy = async () => {
    if (inviteLink) {
      await navigator.clipboard.writeText(inviteLink)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  const handleClose = () => {
    setCopied(false)
    onClose()
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Invite Link Created</DialogTitle>
          <DialogDescription>
            Email is not configured. Share this link with the user to invite them
            to your organization. The link expires in 7 days.
          </DialogDescription>
        </DialogHeader>

        <div className="flex gap-2">
          <Input
            value={inviteLink ?? ''}
            readOnly
            className="font-mono text-sm"
          />
          <Button variant="outline" size="icon" onClick={handleCopy}>
            {copied ? (
              <Check className="h-4 w-4 text-green-500" />
            ) : (
              <Copy className="h-4 w-4" />
            )}
          </Button>
        </div>

        <DialogFooter>
          <Button onClick={handleClose}>Done</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
