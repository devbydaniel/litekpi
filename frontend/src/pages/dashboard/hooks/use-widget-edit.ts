import { useState, useCallback, useMemo, useEffect } from 'react'
import type { Widget, Filter, UpdateWidgetRequest } from '@/shared/api/generated/models'

export type ChartType = 'area' | 'bar' | 'line'
export type DateRangeValue = 'last_24_hours' | 'last_7_days' | 'last_30_days'

export interface WidgetEditState {
  title: string | undefined
  chartType: ChartType
  dateRange: DateRangeValue
  splitBy: string | undefined
  filters: Filter[]
}

interface UseWidgetEditOptions {
  widget: Widget
  onSave: (widgetId: string, update: UpdateWidgetRequest) => Promise<void>
}

function getInitialState(widget: Widget): WidgetEditState {
  return {
    title: widget.title ?? undefined,
    chartType: (widget.chartType as ChartType) ?? 'area',
    dateRange: (widget.dateRange as DateRangeValue) ?? 'last_7_days',
    splitBy: widget.splitBy ?? undefined,
    filters: widget.filters ?? [],
  }
}

function statesEqual(a: WidgetEditState, b: WidgetEditState): boolean {
  if (a.title !== b.title) return false
  if (a.chartType !== b.chartType) return false
  if (a.dateRange !== b.dateRange) return false
  if (a.splitBy !== b.splitBy) return false
  if (a.filters.length !== b.filters.length) return false

  const aFilters = [...a.filters].sort((x, y) => (x.key ?? '').localeCompare(y.key ?? ''))
  const bFilters = [...b.filters].sort((x, y) => (x.key ?? '').localeCompare(y.key ?? ''))

  for (let i = 0; i < aFilters.length; i++) {
    if (aFilters[i].key !== bFilters[i].key || aFilters[i].value !== bFilters[i].value) {
      return false
    }
  }

  return true
}

export function useWidgetEdit({ widget, onSave }: UseWidgetEditOptions) {
  const [isEditing, setIsEditing] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [state, setState] = useState<WidgetEditState>(() => getInitialState(widget))

  // Reset state when widget changes (e.g., after save)
  useEffect(() => {
    setState(getInitialState(widget))
  }, [widget])

  const savedState = useMemo(() => getInitialState(widget), [widget])
  const isDirty = useMemo(() => !statesEqual(state, savedState), [state, savedState])

  const setTitle = useCallback((title: string | undefined) => {
    setState((prev) => ({ ...prev, title }))
  }, [])

  const setChartType = useCallback((chartType: ChartType) => {
    setState((prev) => ({ ...prev, chartType }))
  }, [])

  const setDateRange = useCallback((dateRange: DateRangeValue) => {
    setState((prev) => ({ ...prev, dateRange }))
  }, [])

  const setSplitBy = useCallback((splitBy: string | undefined) => {
    setState((prev) => ({ ...prev, splitBy }))
  }, [])

  const setFilter = useCallback((key: string, value: string | undefined) => {
    setState((prev) => {
      const newFilters = prev.filters.filter((f) => f.key !== key)
      if (value !== undefined) {
        newFilters.push({ key, value })
      }
      return { ...prev, filters: newFilters }
    })
  }, [])

  const clearAllFilters = useCallback(() => {
    setState((prev) => ({ ...prev, filters: [] }))
  }, [])

  const reset = useCallback(() => {
    setState(getInitialState(widget))
  }, [widget])

  const save = useCallback(async () => {
    if (!isDirty || !widget.id) return

    setIsSaving(true)
    try {
      await onSave(widget.id, {
        title: state.title,
        chartType: state.chartType,
        dateRange: state.dateRange,
        splitBy: state.splitBy,
        filters: state.filters,
      })
    } finally {
      setIsSaving(false)
    }
  }, [isDirty, widget.id, state, onSave])

  const toggleEditing = useCallback(() => {
    setIsEditing((prev) => !prev)
  }, [])

  const stopEditing = useCallback(() => {
    setIsEditing(false)
  }, [])

  // Build a widget-like object with current edits for preview
  const previewWidget = useMemo(
    (): Widget => ({
      ...widget,
      title: state.title,
      chartType: state.chartType,
      dateRange: state.dateRange,
      splitBy: state.splitBy,
      filters: state.filters,
    }),
    [widget, state]
  )

  return {
    // State
    isEditing,
    isSaving,
    isDirty,
    state,
    previewWidget,

    // Setters
    setTitle,
    setChartType,
    setDateRange,
    setSplitBy,
    setFilter,
    clearAllFilters,

    // Actions
    toggleEditing,
    stopEditing,
    reset,
    save,
  }
}
