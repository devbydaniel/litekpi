import { useState, useMemo, useCallback, useEffect } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetProductsProductIdMeasurementsNamePreferences,
  useGetProductsProductIdMeasurementsNameMetadata,
  usePostProductsProductIdMeasurementsNamePreferences,
  getGetProductsProductIdMeasurementsNamePreferencesQueryKey,
} from '@/shared/api/generated/api'
import type { AggregatedDataPoint } from '@/shared/api/generated/models'
import { customInstance } from '@/shared/api/client'
import { type DateRangeValue, getDateRangeFromValue } from '../ui/date-range-filter'

interface UseMeasurementChartOptions {
  productId: string
  measurementName: string
}

export type ChartType = 'area' | 'bar' | 'line'

interface SplitSeries {
  key: string
  dataPoints: AggregatedDataPoint[]
}

interface GetMeasurementDataResponse {
  name: string
  dataPoints: AggregatedDataPoint[]
}

interface GetMeasurementDataSplitResponse {
  name: string
  splitBy: string
  series: SplitSeries[]
}

type GetMeasurementDataResult = GetMeasurementDataResponse | GetMeasurementDataSplitResponse

function isSplitResponse(
  response: GetMeasurementDataResult
): response is GetMeasurementDataSplitResponse {
  return 'series' in response
}

// Transform split data for Recharts (pivot to { date, key1: sum, key2: sum, ... })
function transformSplitData(series: SplitSeries[]) {
  const seriesKeys = series.map((s) => s.key)

  // Create a map of date -> { date, key1: sum, key2: sum, ... }
  const dateMap = new Map<string, Record<string, number | string>>()

  for (const s of series) {
    for (const dp of s.dataPoints) {
      if (!dp.date) continue
      if (!dateMap.has(dp.date)) {
        dateMap.set(dp.date, { date: dp.date })
      }
      const entry = dateMap.get(dp.date)!
      entry[s.key] = dp.sum ?? 0
    }
  }

  // Convert to array sorted by date
  const data = Array.from(dateMap.values()).sort((a, b) =>
    (a.date as string).localeCompare(b.date as string)
  )

  return { data, seriesKeys }
}

// Custom fetch for measurement data with metadata filter support
async function fetchMeasurementData(
  productId: string,
  name: string,
  params: {
    start: string
    end: string
    metadata?: Record<string, string>
    splitBy?: string
  }
): Promise<GetMeasurementDataResult> {
  const queryParams: Record<string, string> = {
    start: params.start,
    end: params.end,
  }

  // Add metadata filters with "metadata." prefix
  if (params.metadata) {
    for (const [key, value] of Object.entries(params.metadata)) {
      queryParams[`metadata.${key}`] = value
    }
  }

  if (params.splitBy) {
    queryParams.splitBy = params.splitBy
  }

  return customInstance<GetMeasurementDataResult>({
    url: `/products/${productId}/measurements/${encodeURIComponent(name)}/data`,
    method: 'GET',
    params: queryParams,
  })
}

export function useMeasurementChart({ productId, measurementName }: UseMeasurementChartOptions) {
  const queryClient = useQueryClient()
  const [chartType, setChartType] = useState<ChartType>('area')
  const [dateRange, setDateRange] = useState<DateRangeValue>('last7days')
  const [metadataFilters, setMetadataFilters] = useState<Record<string, string>>({})
  const [splitBy, setSplitBy] = useState<string | undefined>(undefined)
  const [preferencesApplied, setPreferencesApplied] = useState(false)

  // Fetch saved preferences
  const { data: preferencesData } = useGetProductsProductIdMeasurementsNamePreferences(
    productId,
    measurementName
  )

  // Apply preferences when they load (only once)
  useEffect(() => {
    if (preferencesData?.preferences && !preferencesApplied) {
      const prefs = preferencesData.preferences
      setChartType((prefs.chartType as ChartType) ?? 'area')
      setDateRange((prefs.dateRange as DateRangeValue) ?? 'last7days')
      setSplitBy(prefs.splitBy ?? undefined)
      setMetadataFilters((prefs.metadataFilters as Record<string, string>) ?? {})
      setPreferencesApplied(true)
    }
  }, [preferencesData, preferencesApplied])

  // Save preferences mutation
  const savePreferencesMutation = usePostProductsProductIdMeasurementsNamePreferences({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetProductsProductIdMeasurementsNamePreferencesQueryKey(
            productId,
            measurementName
          ),
        })
        toast.success('Default view saved')
      },
    },
  })

  const savePreferences = useCallback(() => {
    savePreferencesMutation.mutate({
      productId,
      name: measurementName,
      data: {
        preferences: {
          chartType,
          dateRange,
          splitBy: splitBy ?? undefined,
          metadataFilters,
        },
      },
    })
  }, [productId, measurementName, chartType, dateRange, splitBy, metadataFilters, savePreferencesMutation])

  // Fetch metadata values for filter dropdowns
  const { data: metadataData, isLoading: isLoadingMetadata } =
    useGetProductsProductIdMeasurementsNameMetadata(productId, measurementName)

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

  // Fetch chart data based on current filters (using custom fetch for metadata filter support)
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
      fetchMeasurementData(productId, measurementName, {
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
