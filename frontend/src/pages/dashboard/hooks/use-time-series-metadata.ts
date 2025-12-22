import { useMemo } from 'react'
import { useGetDataSourcesDataSourceIdMeasurementsNameMetadata } from '@/shared/api/generated/api'

export interface TimeSeriesMetadata {
  keys: string[]
  values: Record<string, string[]>
}

export function useTimeSeriesMetadata(dataSourceId: string, measurementName: string) {
  const { data, isLoading } = useGetDataSourcesDataSourceIdMeasurementsNameMetadata(
    dataSourceId,
    measurementName,
    { query: { enabled: !!dataSourceId && !!measurementName } }
  )

  const metadata = useMemo<TimeSeriesMetadata>(() => {
    const keys: string[] = []
    const values: Record<string, string[]> = {}

    for (const item of data?.metadata ?? []) {
      if (item.key) {
        keys.push(item.key)
        values[item.key] = item.values ?? []
      }
    }

    return { keys, values }
  }, [data])

  return {
    metadata,
    isLoading,
  }
}
