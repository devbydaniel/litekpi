import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import type {
  GetMeasurementDataResponse,
  GetMeasurementDataSplitResponse,
  SplitSeries,
  Widget,
  Filter,
} from '@/shared/api/generated/models'
import { customInstance } from '@/shared/api/client'

export type ChartType = 'area' | 'bar' | 'line'
export type DateRangeValue = 'last24h' | 'last7days' | 'last30days'

// Transform split data for Recharts (pivot to { date, key1: sum, key2: sum, ... })
function transformSplitData(series: SplitSeries[]) {
  const seriesKeys = series.map((s) => s.key ?? '')

  // Create a map of date -> { date, key1: sum, key2: sum, ... }
  const dateMap = new Map<string, Record<string, number | string>>()

  for (const s of series) {
    const key = s.key ?? ''
    for (const dp of s.dataPoints ?? []) {
      if (!dp.date) continue
      if (!dateMap.has(dp.date)) {
        dateMap.set(dp.date, { date: dp.date })
      }
      const entry = dateMap.get(dp.date)!
      entry[key] = dp.sum ?? 0
    }
  }

  // Convert to array sorted by date
  const data = Array.from(dateMap.values()).sort((a, b) =>
    (a.date as string).localeCompare(b.date as string)
  )

  return { data, seriesKeys }
}

function getDateRangeFromValue(value: string): { start: Date; end: Date } {
  const now = new Date()
  const end = new Date(now)

  let start: Date
  switch (value) {
    case 'last24h':
    case 'last_24_hours':
      start = new Date(now.getTime() - 24 * 60 * 60 * 1000)
      break
    case 'last7days':
    case 'last_7_days':
      start = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
      break
    case 'last30days':
    case 'last_30_days':
      start = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000)
      break
    default:
      start = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
  }

  return { start, end }
}

// Custom fetch for measurement data with metadata filter support
async function fetchMeasurementData(
  dataSourceId: string,
  name: string,
  params: {
    start: string
    end: string
    metadata?: Record<string, string>
  }
): Promise<GetMeasurementDataResponse> {
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

  return customInstance<GetMeasurementDataResponse>({
    url: `/data-sources/${dataSourceId}/measurements/${encodeURIComponent(name)}/data`,
    method: 'GET',
    params: queryParams,
  })
}

// Custom fetch for split measurement data with metadata filter support
async function fetchMeasurementDataSplit(
  dataSourceId: string,
  name: string,
  params: {
    start: string
    end: string
    splitBy: string
    metadata?: Record<string, string>
  }
): Promise<GetMeasurementDataSplitResponse> {
  const queryParams: Record<string, string> = {
    start: params.start,
    end: params.end,
    splitBy: params.splitBy,
  }

  // Add metadata filters with "metadata." prefix
  if (params.metadata) {
    for (const [key, value] of Object.entries(params.metadata)) {
      queryParams[`metadata.${key}`] = value
    }
  }

  return customInstance<GetMeasurementDataSplitResponse>({
    url: `/data-sources/${dataSourceId}/measurements/${encodeURIComponent(name)}/data/split`,
    method: 'GET',
    params: queryParams,
  })
}

export function useWidgetData(widget: Widget) {
  const dataSourceId = widget.dataSourceId ?? ''
  const measurementName = widget.measurementName ?? ''
  const dateRange = widget.dateRange ?? 'last_7_days'
  const splitBy = widget.splitBy
  const filters = widget.filters ?? []

  // Calculate start and end dates from date range preset
  const { start, end } = useMemo(() => getDateRangeFromValue(dateRange), [dateRange])

  // Build clean metadata filters (from widget filters)
  const cleanMetadataFilters = useMemo(() => {
    const clean: Record<string, string> = {}
    for (const filter of filters as Filter[]) {
      if (filter.key && filter.value) {
        clean[filter.key] = filter.value
      }
    }
    return clean
  }, [filters])

  const hasMetadataFilters = Object.keys(cleanMetadataFilters).length > 0

  // Fetch non-split chart data
  const { data: nonSplitData, isLoading: isLoadingNonSplit } = useQuery({
    queryKey: ['widget', widget.id, 'data', dataSourceId, measurementName, dateRange, cleanMetadataFilters],
    queryFn: () =>
      fetchMeasurementData(dataSourceId, measurementName, {
        start: start.toISOString(),
        end: end.toISOString(),
        metadata: hasMetadataFilters ? cleanMetadataFilters : undefined,
      }),
    enabled: !splitBy && !!dataSourceId && !!measurementName,
  })

  // Fetch split chart data
  const { data: splitData, isLoading: isLoadingSplit } = useQuery({
    queryKey: [
      'widget',
      widget.id,
      'data',
      'split',
      dataSourceId,
      measurementName,
      dateRange,
      cleanMetadataFilters,
      splitBy,
    ],
    queryFn: () =>
      fetchMeasurementDataSplit(dataSourceId, measurementName, {
        start: start.toISOString(),
        end: end.toISOString(),
        splitBy: splitBy!,
        metadata: hasMetadataFilters ? cleanMetadataFilters : undefined,
      }),
    enabled: !!splitBy && !!dataSourceId && !!measurementName,
  })

  const isLoading = splitBy ? isLoadingSplit : isLoadingNonSplit

  // Transform data based on whether it's split or not
  const { data, seriesKeys } = useMemo(() => {
    if (splitBy && splitData?.series) {
      return transformSplitData(splitData.series)
    }

    if (!splitBy && nonSplitData?.dataPoints) {
      return {
        data: nonSplitData.dataPoints,
        seriesKeys: [] as string[],
      }
    }

    return { data: [], seriesKeys: [] as string[] }
  }, [splitBy, splitData, nonSplitData])

  return {
    data,
    seriesKeys,
    isSplit: splitBy !== undefined && splitBy !== null && splitBy !== '',
    isLoading,
    chartType: (widget.chartType ?? 'area') as ChartType,
  }
}
