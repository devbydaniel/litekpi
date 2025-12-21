import { useState, useEffect } from 'react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import { Label } from '@/shared/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import type { User, Role } from '@/shared/api/generated/models'

interface ChangeRoleDialogProps {
  open: boolean
  user: User | null
  isLoading: boolean
  onConfirm: (role: Role) => Promise<void>
  onClose: () => void
}

export function ChangeRoleDialog({
  open,
  user,
  isLoading,
  onConfirm,
  onClose,
}: ChangeRoleDialogProps) {
  const [role, setRole] = useState<Role>('viewer')

  useEffect(() => {
    if (user) {
      setRole(user.role as Role)
    }
  }, [user])

  const handleConfirm = async () => {
    await onConfirm(role)
  }

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Change Role</DialogTitle>
          <DialogDescription>
            Change the role for{' '}
            <span className="font-medium">{user?.name}</span>.
          </DialogDescription>
        </DialogHeader>

        <div className="grid gap-2 py-4">
          <Label htmlFor="role">Role</Label>
          <Select value={role} onValueChange={(v) => setRole(v as Role)}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="admin">Admin - Full access</SelectItem>
              <SelectItem value="editor">Editor - Can edit data</SelectItem>
              <SelectItem value="viewer">Viewer - Read-only access</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose} disabled={isLoading}>
            Cancel
          </Button>
          <Button onClick={handleConfirm} disabled={isLoading || role === user?.role}>
            {isLoading ? 'Updating...' : 'Update Role'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
