import { useState, useMemo, useCallback, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  measurementsApi,
  isSplitResponse,
  type SplitSeries,
  type ChartType,
  type MeasurementPreferences,
} from '@/shared/api/measurements'
import { type DateRangeValue, getDateRangeFromValue } from '../ui/date-range-filter'

interface UseMeasurementChartOptions {
  productId: string
  measurementName: string
}

// Transform split data for Recharts (pivot to { date, key1: sum, key2: sum, ... })
function transformSplitData(series: SplitSeries[]) {
  const seriesKeys = series.map((s) => s.key)

  // Create a map of date -> { date, key1: sum, key2: sum, ... }
  const dateMap = new Map<string, Record<string, number | string>>()

  for (const s of series) {
    for (const dp of s.dataPoints) {
      if (!dateMap.has(dp.date)) {
        dateMap.set(dp.date, { date: dp.date })
      }
      const entry = dateMap.get(dp.date)!
      entry[s.key] = dp.sum
    }
  }

  // Convert to array sorted by date
  const data = Array.from(dateMap.values()).sort((a, b) =>
    (a.date as string).localeCompare(b.date as string)
  )

  return { data, seriesKeys }
}

export function useMeasurementChart({ productId, measurementName }: UseMeasurementChartOptions) {
  const queryClient = useQueryClient()
  const [chartType, setChartType] = useState<ChartType>('area')
  const [dateRange, setDateRange] = useState<DateRangeValue>('last7days')
  const [metadataFilters, setMetadataFilters] = useState<Record<string, string>>({})
  const [splitBy, setSplitBy] = useState<string | undefined>(undefined)
  const [preferencesApplied, setPreferencesApplied] = useState(false)

  // Fetch saved preferences
  const { data: preferencesData } = useQuery({
    queryKey: ['measurements', productId, measurementName, 'preferences'],
    queryFn: () => measurementsApi.getPreferences(productId, measurementName),
  })

  // Apply preferences when they load (only once)
  useEffect(() => {
    if (preferencesData?.preferences && !preferencesApplied) {
      const prefs = preferencesData.preferences
      setChartType(prefs.chartType)
      setDateRange(prefs.dateRange)
      setSplitBy(prefs.splitBy ?? undefined)
      setMetadataFilters(prefs.metadataFilters ?? {})
      setPreferencesApplied(true)
    }
  }, [preferencesData, preferencesApplied])

  // Save preferences mutation
  const savePreferencesMutation = useMutation({
    mutationFn: (prefs: MeasurementPreferences) =>
      measurementsApi.savePreferences(productId, measurementName, prefs),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['measurements', productId, measurementName, 'preferences'],
      })
      toast.success('Default view saved')
    },
  })

  const savePreferences = useCallback(() => {
    savePreferencesMutation.mutate({
      chartType,
      dateRange,
      splitBy: splitBy ?? null,
      metadataFilters,
    })
  }, [chartType, dateRange, splitBy, metadataFilters, savePreferencesMutation])

  // Fetch metadata values for filter dropdowns
  const { data: metadataData, isLoading: isLoadingMetadata } = useQuery({
    queryKey: ['measurements', productId, measurementName, 'metadata'],
    queryFn: () => measurementsApi.getMetadataValues(productId, measurementName),
  })

  // Calculate start and end dates from date range preset
  const { start, end } = useMemo(() => getDateRangeFromValue(dateRange), [dateRange])

  // Build clean metadata filters (only non-empty values)
  const cleanMetadataFilters = useMemo(() => {
    const clean: Record<string, string> = {}
    for (const [key, value] of Object.entries(metadataFilters)) {
      if (value) {
        clean[key] = value
      }
    }
    return clean
  }, [metadataFilters])

  // Fetch chart data based on current filters
  const { data: chartData, isLoading: isLoadingData } = useQuery({
    queryKey: [
      'measurements',
      productId,
      measurementName,
      'data',
      dateRange,
      cleanMetadataFilters,
      splitBy,
    ],
    queryFn: () =>
      measurementsApi.getData(productId, measurementName, {
        start: start.toISOString(),
        end: end.toISOString(),
        metadata: Object.keys(cleanMetadataFilters).length > 0 ? cleanMetadataFilters : undefined,
        splitBy,
      }),
  })

  // Transform data based on whether it's split or not
  const { data, seriesKeys } = useMemo(() => {
    if (!chartData) return { data: [], seriesKeys: [] as string[] }

    if (isSplitResponse(chartData)) {
      return transformSplitData(chartData.series)
    }

    return {
      data: chartData.dataPoints,
      seriesKeys: [] as string[],
    }
  }, [chartData])

  const setMetadataFilter = useCallback((key: string, value: string | undefined) => {
    setMetadataFilters((prev) => {
      const next = { ...prev }
      if (value === undefined) {
        delete next[key]
      } else {
        next[key] = value
      }
      return next
    })
  }, [])

  const clearAllFilters = useCallback(() => setMetadataFilters({}), [])

  // Count active filters
  const activeFilterCount = useMemo(
    () => Object.keys(cleanMetadataFilters).length,
    [cleanMetadataFilters]
  )

  // Check if current config differs from saved preferences
  const isDirty = useMemo(() => {
    const saved = preferencesData?.preferences
    if (!saved) return false
    return (
      chartType !== saved.chartType ||
      dateRange !== saved.dateRange ||
      splitBy !== (saved.splitBy ?? undefined) ||
      JSON.stringify(cleanMetadataFilters) !== JSON.stringify(saved.metadataFilters ?? {})
    )
  }, [preferencesData, chartType, dateRange, splitBy, cleanMetadataFilters])

  return {
    data,
    seriesKeys,
    isSplit: splitBy !== undefined,
    metadata: metadataData?.metadata ?? [],
    chartType,
    dateRange,
    metadataFilters,
    splitBy,
    setChartType,
    setDateRange,
    setMetadataFilter,
    clearAllFilters,
    setSplitBy,
    savePreferences,
    isSaving: savePreferencesMutation.isPending,
    isLoading: isLoadingData || isLoadingMetadata,
    isDirty,
    activeFilterCount,
  }
}
