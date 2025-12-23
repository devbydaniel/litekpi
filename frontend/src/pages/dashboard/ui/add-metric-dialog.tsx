import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/shared/components/ui/dialog'
import { MetricForm } from './metric-form'
import type { CreateMetricRequest } from '@/shared/api/generated/models'

interface AddMetricDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAdd: (metric: CreateMetricRequest) => Promise<void>
  isLoading: boolean
}

export function AddMetricDialog({
  open,
  onOpenChange,
  onAdd,
  isLoading,
}: AddMetricDialogProps) {
  const handleSubmit = async (values: CreateMetricRequest) => {
    await onAdd(values)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Add Metric</DialogTitle>
          <DialogDescription>
            Create a metric to display data as a single value or a chart over time.
          </DialogDescription>
        </DialogHeader>

        <MetricForm
          onSubmit={handleSubmit}
          onCancel={() => onOpenChange(false)}
          isLoading={isLoading}
          submitLabel="Add Metric"
        />
      </DialogContent>
    </Dialog>
  )
}
