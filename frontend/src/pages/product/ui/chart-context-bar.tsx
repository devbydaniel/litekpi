import { Loader2, Save, X } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { ActiveFilterChips } from './active-filter-chips'

interface ChartContextBarProps {
  filters: Record<string, string>
  isDirty: boolean
  isSaving: boolean
  onRemoveFilter: (key: string) => void
  onClearAll: () => void
  onSave: () => void
}

export function ChartContextBar({
  filters,
  isDirty,
  isSaving,
  onRemoveFilter,
  onClearAll,
  onSave,
}: ChartContextBarProps) {
  const activeFilters = Object.entries(filters)
    .filter(([, value]) => value !== undefined && value !== '')
    .map(([key, value]) => ({ key, value }))

  const hasActiveFilters = activeFilters.length > 0

  // Only render if there are active filters or unsaved changes
  if (!hasActiveFilters && !isDirty) {
    return null
  }

  return (
    <div className="flex flex-col gap-2 border-b px-6 py-2 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex flex-wrap items-center gap-2">
        <ActiveFilterChips filters={activeFilters} onRemove={onRemoveFilter} />
        {hasActiveFilters && (
          <Button
            variant="ghost"
            size="sm"
            className="h-6 gap-1 px-2 text-xs text-muted-foreground"
            onClick={onClearAll}
          >
            <X className="h-3 w-3" />
            Clear all
          </Button>
        )}
      </div>
      {isDirty && (
        <Button
          variant="outline"
          size="sm"
          className="gap-2"
          onClick={onSave}
          disabled={isSaving}
        >
          {isSaving ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Save className="h-4 w-4" />
          )}
          Save as default
        </Button>
      )}
    </div>
  )
}
