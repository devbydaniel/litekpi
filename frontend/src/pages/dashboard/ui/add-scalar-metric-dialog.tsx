import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/shared/components/ui/dialog'
import { ScalarMetricForm } from '@/shared/components/scalar-metric-form'
import type { CreateScalarMetricRequest } from '@/shared/api/generated/models'

interface AddScalarMetricDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAdd: (metric: CreateScalarMetricRequest) => Promise<void>
  isLoading: boolean
}

export function AddScalarMetricDialog({
  open,
  onOpenChange,
  onAdd,
  isLoading,
}: AddScalarMetricDialogProps) {
  const handleSubmit = async (values: CreateScalarMetricRequest) => {
    await onAdd(values)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Add Metric</DialogTitle>
          <DialogDescription>
            Create a metric card to display a single aggregated value.
          </DialogDescription>
        </DialogHeader>

        <ScalarMetricForm
          onSubmit={handleSubmit}
          onCancel={() => onOpenChange(false)}
          isLoading={isLoading}
          submitLabel="Add Metric"
        />
      </DialogContent>
    </Dialog>
  )
}
