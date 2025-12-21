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
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import type { Report } from '@/shared/api/generated/models'

interface EditReportDialogProps {
  report: Report | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onUpdate: (id: string, name: string) => Promise<void>
  isLoading: boolean
}

export function EditReportDialog({
  report,
  open,
  onOpenChange,
  onUpdate,
  isLoading,
}: EditReportDialogProps) {
  const [name, setName] = useState('')

  useEffect(() => {
    if (report) setName(report.name ?? '')
  }, [report])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!report?.id || !name.trim()) return
    await onUpdate(report.id, name.trim())
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Rename Report</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="editName">Name</Label>
            <Input
              id="editName"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Report name"
              maxLength={255}
              autoFocus
            />
          </div>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading || !name.trim()}>
              {isLoading ? 'Saving...' : 'Save'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

interface DeleteReportDialogProps {
  report: Report | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onDelete: (id: string) => Promise<void>
  isLoading: boolean
}

export function DeleteReportDialog({
  report,
  open,
  onOpenChange,
  onDelete,
  isLoading,
}: DeleteReportDialogProps) {
  const handleDelete = async () => {
    if (!report?.id) return
    await onDelete(report.id)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Delete Report</DialogTitle>
          <DialogDescription>
            Are you sure you want to delete &quot;{report?.name}&quot;? This will also delete
            all KPIs in this report. This action cannot be undone.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isLoading}
          >
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete} disabled={isLoading}>
            {isLoading ? 'Deleting...' : 'Delete'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
