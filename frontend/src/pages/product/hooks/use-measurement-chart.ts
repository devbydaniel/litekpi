import { useState, useMemo, useCallback } from 'react'
import { useQuery } from '@tanstack/react-query'
import { measurementsApi } from '@/shared/api/measurements'
import { type DateRangeValue, getDateRangeFromValue } from '../ui/date-range-filter'

interface UseMeasurementChartOptions {
  productId: string
  measurementName: string
}

export function useMeasurementChart({ productId, measurementName }: UseMeasurementChartOptions) {
  const [dateRange, setDateRange] = useState<DateRangeValue>('last7days')
  const [metadataFilters, setMetadataFilters] = useState<Record<string, string>>({})

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
    ],
    queryFn: () =>
      measurementsApi.getData(productId, measurementName, {
        start: start.toISOString(),
        end: end.toISOString(),
        metadata: Object.keys(cleanMetadataFilters).length > 0 ? cleanMetadataFilters : undefined,
      }),
  })

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

  return {
    data: chartData?.dataPoints ?? [],
    metadata: metadataData?.metadata ?? [],
    dateRange,
    metadataFilters,
    setDateRange,
    setMetadataFilter,
    isLoading: isLoadingData || isLoadingMetadata,
  }
}
