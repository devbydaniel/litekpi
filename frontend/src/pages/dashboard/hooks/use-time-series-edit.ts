import { useState, useCallback, useMemo, useEffect } from 'react'
import type { TimeSeries, Filter, UpdateTimeSeriesRequest } from '@/shared/api/generated/models'

export type ChartType = 'area' | 'bar' | 'line'
export type DateRangeValue = 'last_24_hours' | 'last_7_days' | 'last_30_days'

export interface TimeSeriesEditState {
  title: string | undefined
  chartType: ChartType
  dateRange: DateRangeValue
  splitBy: string | undefined
  filters: Filter[]
}

interface UseTimeSeriesEditOptions {
  timeSeries: TimeSeries
  onSave: (timeSeriesId: string, update: UpdateTimeSeriesRequest) => Promise<void>
}

function getInitialState(timeSeries: TimeSeries): TimeSeriesEditState {
  return {
    title: timeSeries.title ?? undefined,
    chartType: (timeSeries.chartType as ChartType) ?? 'area',
    dateRange: (timeSeries.dateRange as DateRangeValue) ?? 'last_7_days',
    splitBy: timeSeries.splitBy ?? undefined,
    filters: timeSeries.filters ?? [],
  }
}

function statesEqual(a: TimeSeriesEditState, b: TimeSeriesEditState): boolean {
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

export function useTimeSeriesEdit({ timeSeries, onSave }: UseTimeSeriesEditOptions) {
  const [isEditing, setIsEditing] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [state, setState] = useState<TimeSeriesEditState>(() => getInitialState(timeSeries))

  // Reset state when time series changes (e.g., after save)
  useEffect(() => {
    setState(getInitialState(timeSeries))
  }, [timeSeries])

  const savedState = useMemo(() => getInitialState(timeSeries), [timeSeries])
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
    setState(getInitialState(timeSeries))
  }, [timeSeries])

  const save = useCallback(async () => {
    if (!isDirty || !timeSeries.id) return

    setIsSaving(true)
    try {
      await onSave(timeSeries.id, {
        title: state.title,
        chartType: state.chartType,
        dateRange: state.dateRange,
        splitBy: state.splitBy,
        filters: state.filters,
      })
    } finally {
      setIsSaving(false)
    }
  }, [isDirty, timeSeries.id, state, onSave])

  const toggleEditing = useCallback(() => {
    setIsEditing((prev) => !prev)
  }, [])

  const stopEditing = useCallback(() => {
    setIsEditing(false)
  }, [])

  // Build a time series-like object with current edits for preview
  const previewTimeSeries = useMemo(
    (): TimeSeries => ({
      ...timeSeries,
      title: state.title,
      chartType: state.chartType,
      dateRange: state.dateRange,
      splitBy: state.splitBy,
      filters: state.filters,
    }),
    [timeSeries, state]
  )

  return {
    // State
    isEditing,
    isSaving,
    isDirty,
    state,
    previewTimeSeries,

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
