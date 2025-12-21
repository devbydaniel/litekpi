import { useState } from 'react'
import { MoreHorizontal, Pencil, Trash, BarChart } from 'lucide-react'
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
  DialogDescription,
} from '@/shared/components/ui/dialog'
import { EmptyState } from '@/shared/components/ui/empty-state'
import { KpiCard, KpiCardSkeleton } from '@/shared/components/kpi-card'
import { KpiForm } from '@/shared/components/kpi-form'
import { useReportDetail } from '../hooks/use-report-detail'
import type { Report, Kpi, UpdateKPIRequest, CreateKPIRequest } from '@/shared/api/generated/models'

interface ReportDetailProps {
  report: Report
  canEdit: boolean
  addKpiOpen?: boolean
  onAddKpiOpenChange?: (open: boolean) => void
}

export function ReportDetail({
  report,
  canEdit,
  addKpiOpen = false,
  onAddKpiOpenChange,
}: ReportDetailProps) {
  const [internalAddKpiOpen, setInternalAddKpiOpen] = useState(false)
  const [editingKpi, setEditingKpi] = useState<Kpi | null>(null)

  // Use external state if provided, otherwise internal
  const isAddKpiOpen = onAddKpiOpenChange ? addKpiOpen : internalAddKpiOpen
  const setAddKpiOpen = onAddKpiOpenChange ?? setInternalAddKpiOpen

  const {
    kpis,
    computedKpis,
    isLoadingReport,
    isComputingKpis,
    addKpi,
    updateKpi,
    deleteKpi,
    isAddingKpi,
    isUpdatingKpi,
  } = useReportDetail(report.id)

  const handleAddKpi = async (values: CreateKPIRequest) => {
    await addKpi(values)
    setAddKpiOpen(false)
  }

  const handleEditKpi = async (values: UpdateKPIRequest) => {
    if (!editingKpi?.id) return
    await updateKpi(editingKpi.id, values)
    setEditingKpi(null)
  }

  const isLoading = isLoadingReport || isComputingKpis

  return (
    <>
      {/* KPI Grid */}
      {isLoading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {[1, 2, 3, 4].map((i) => (
            <KpiCardSkeleton key={i} />
          ))}
        </div>
      ) : computedKpis.length === 0 ? (
        <EmptyState
          icon={BarChart}
          title="No KPIs yet"
          description={
            canEdit
              ? 'Add KPIs to start tracking metrics in this report.'
              : 'No KPIs have been added to this report.'
          }
        />
      ) : (
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
                        onClick={() => kpi.id && deleteKpi(kpi.id)}
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
      )}

      {/* Add KPI Dialog */}
      <Dialog open={isAddKpiOpen} onOpenChange={setAddKpiOpen}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Add KPI</DialogTitle>
            <DialogDescription>
              Add a KPI to track a metric in this report.
            </DialogDescription>
          </DialogHeader>
          <KpiForm
            onSubmit={handleAddKpi}
            onCancel={() => setAddKpiOpen(false)}
            isLoading={isAddingKpi}
            submitLabel="Add KPI"
          />
        </DialogContent>
      </Dialog>

      {/* Edit KPI Dialog */}
      <Dialog open={!!editingKpi} onOpenChange={(open) => !open && setEditingKpi(null)}>
        <DialogContent className="max-w-lg">
          <DialogHeader>
            <DialogTitle>Edit KPI</DialogTitle>
          </DialogHeader>
          {editingKpi && (
            <KpiForm
              initialValues={editingKpi}
              onSubmit={handleEditKpi}
              onCancel={() => setEditingKpi(null)}
              isLoading={isUpdatingKpi}
              submitLabel="Save Changes"
            />
          )}
        </DialogContent>
      </Dialog>
    </>
  )
}
