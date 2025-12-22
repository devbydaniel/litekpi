import { MoreHorizontal, Trash, Pencil } from 'lucide-react'
import { useState } from 'react'
import { Button } from '@/shared/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog'
import { ScalarMetricCard, ScalarMetricCardSkeleton } from '@/shared/components/scalar-metric-card'
import { ScalarMetricForm } from '@/shared/components/scalar-metric-form'
import type { ComputedScalarMetric, ScalarMetric, UpdateScalarMetricRequest } from '@/shared/api/generated/models'

interface ScalarMetricGridProps {
  metrics: ScalarMetric[]
  computedMetrics: ComputedScalarMetric[]
  isLoading: boolean
  canEdit: boolean
  onUpdate: (metricId: string, metric: UpdateScalarMetricRequest) => Promise<void>
  onDelete: (metricId: string) => Promise<void>
  isUpdating: boolean
}

export function ScalarMetricGrid({
  metrics,
  computedMetrics,
  isLoading,
  canEdit,
  onUpdate,
  onDelete,
  isUpdating,
}: ScalarMetricGridProps) {
  const [editingMetric, setEditingMetric] = useState<ScalarMetric | null>(null)

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {[1, 2, 3, 4].map((i) => (
          <ScalarMetricCardSkeleton key={i} />
        ))}
      </div>
    )
  }

  if (computedMetrics.length === 0) {
    return null
  }

  const handleEdit = async (values: UpdateScalarMetricRequest) => {
    if (!editingMetric?.id) return
    await onUpdate(editingMetric.id, values)
    setEditingMetric(null)
  }

  return (
    <>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {computedMetrics.map((metric) => (
          <div key={metric.id} className="group relative">
            <ScalarMetricCard metric={metric} />
            {canEdit && (
              <div className="absolute right-2 top-2 opacity-0 transition-opacity group-hover:opacity-100">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" className="h-8 w-8">
                      <MoreHorizontal className="h-4 w-4" />
                      <span className="sr-only">Metric options</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      onClick={() => {
                        const originalMetric = metrics.find((m) => m.id === metric.id)
                        if (originalMetric) setEditingMetric(originalMetric)
                      }}
                    >
                      <Pencil className="mr-2 h-4 w-4" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      className="text-destructive focus:text-destructive"
                      onClick={() => metric.id && onDelete(metric.id)}
                    >
                      <Trash className="mr-2 h-4 w-4" />
                      Remove
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            )}
          </div>
        ))}
      </div>

      <Dialog open={!!editingMetric} onOpenChange={(open) => !open && setEditingMetric(null)}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit Metric</DialogTitle>
          </DialogHeader>
          {editingMetric && (
            <ScalarMetricForm
              initialValues={editingMetric}
              onSubmit={handleEdit}
              onCancel={() => setEditingMetric(null)}
              isLoading={isUpdating}
              submitLabel="Save Changes"
            />
          )}
        </DialogContent>
      </Dialog>
    </>
  )
}
