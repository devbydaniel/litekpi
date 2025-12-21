import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/shared/components/ui/dialog'
import { KpiForm } from '@/shared/components/kpi-form'
import type { CreateKPIRequest } from '@/shared/api/generated/models'

interface AddKpiDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAdd: (kpi: CreateKPIRequest) => Promise<void>
  isLoading: boolean
}

export function AddKpiDialog({
  open,
  onOpenChange,
  onAdd,
  isLoading,
}: AddKpiDialogProps) {
  const handleSubmit = async (values: CreateKPIRequest) => {
    await onAdd(values)
    onOpenChange(false)
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Add KPI</DialogTitle>
          <DialogDescription>
            Create a KPI card to display a single aggregated metric value.
          </DialogDescription>
        </DialogHeader>

        <KpiForm
          onSubmit={handleSubmit}
          onCancel={() => onOpenChange(false)}
          isLoading={isLoading}
          submitLabel="Add KPI"
        />
      </DialogContent>
    </Dialog>
  )
}
