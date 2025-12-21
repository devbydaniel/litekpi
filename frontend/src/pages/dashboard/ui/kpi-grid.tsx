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
import { KpiCard, KpiCardSkeleton } from '@/shared/components/kpi-card'
import { KpiForm } from '@/shared/components/kpi-form'
import type { ComputedKPI, Kpi, UpdateKPIRequest } from '@/shared/api/generated/models'

interface KpiGridProps {
  kpis: Kpi[]
  computedKpis: ComputedKPI[]
  isLoading: boolean
  canEdit: boolean
  onUpdate: (kpiId: string, kpi: UpdateKPIRequest) => Promise<void>
  onDelete: (kpiId: string) => Promise<void>
  isUpdating: boolean
}

export function KpiGrid({
  kpis,
  computedKpis,
  isLoading,
  canEdit,
  onUpdate,
  onDelete,
  isUpdating,
}: KpiGridProps) {
  const [editingKpi, setEditingKpi] = useState<Kpi | null>(null)

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {[1, 2, 3, 4].map((i) => (
          <KpiCardSkeleton key={i} />
        ))}
      </div>
    )
  }

  if (computedKpis.length === 0) {
    return null
  }

  const handleEdit = async (values: UpdateKPIRequest) => {
    if (!editingKpi?.id) return
    await onUpdate(editingKpi.id, values)
    setEditingKpi(null)
  }

  return (
    <>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {computedKpis.map((kpi) => (
          <div key={kpi.id} className="group relative">
            <KpiCard kpi={kpi} />
            {canEdit && (
              <div className="absolute right-2 top-2 opacity-0 transition-opacity group-hover:opacity-100">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" className="h-8 w-8">
                      <MoreHorizontal className="h-4 w-4" />
                      <span className="sr-only">KPI options</span>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem
                      onClick={() => {
                        const originalKpi = kpis.find((k) => k.id === kpi.id)
                        if (originalKpi) setEditingKpi(originalKpi)
                      }}
                    >
                      <Pencil className="mr-2 h-4 w-4" />
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      className="text-destructive focus:text-destructive"
                      onClick={() => kpi.id && onDelete(kpi.id)}
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

      <Dialog open={!!editingKpi} onOpenChange={(open) => !open && setEditingKpi(null)}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit KPI</DialogTitle>
          </DialogHeader>
          {editingKpi && (
            <KpiForm
              initialValues={editingKpi}
              onSubmit={handleEdit}
              onCancel={() => setEditingKpi(null)}
              isLoading={isUpdating}
              submitLabel="Save Changes"
            />
          )}
        </DialogContent>
      </Dialog>
    </>
  )
}
