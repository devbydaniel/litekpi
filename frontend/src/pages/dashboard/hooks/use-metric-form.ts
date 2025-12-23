import { useEffect, useMemo } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  useGetDataSources,
  useGetDataSourcesDataSourceIdMeasurements,
  useGetDataSourcesDataSourceIdMeasurementsNameMetadata,
} from '@/shared/api/generated/api'
import type {
  CreateMetricRequest,
  UpdateMetricRequest,
  Metric,
} from '@/shared/api/generated/models'

const metricFormSchema = z
  .object({
    dataSourceId: z.string().min(1, 'Data source is required'),
    measurementName: z.string().min(1, 'Measurement is required'),
    label: z.string().max(255, 'Label is too long').optional(),
    displayMode: z.enum(['scalar', 'time_series']),
    aggregation: z.enum(['sum', 'average', 'count', 'count_unique']),
    aggregationKey: z.string().optional(),
    granularity: z.enum(['daily', 'weekly', 'monthly']).optional(),
    timeframe: z.enum([
      'last_7_days',
      'last_30_days',
      'this_month',
      'last_month',
    ]),
    comparisonEnabled: z.boolean(),
    comparisonDisplayType: z.enum(['percent', 'absolute']),
    chartType: z.enum(['area', 'bar', 'line']),
    splitBy: z.string().optional(),
    filters: z.array(z.object({ key: z.string(), value: z.string() })),
  })
  .refine(
    (data) => {
      if (data.aggregation === 'count_unique') {
        return !!data.aggregationKey
      }
      return true
    },
    {
      message: 'Aggregation key is required for Count Unique',
      path: ['aggregationKey'],
    }
  )
  .refine(
    (data) => {
      if (data.displayMode === 'time_series') {
        return !!data.granularity
      }
      return true
    },
    {
      message: 'Granularity is required for Time Series',
      path: ['granularity'],
    }
  )

export type MetricFormValues = z.infer<typeof metricFormSchema>

const defaultValues: MetricFormValues = {
  dataSourceId: '',
  measurementName: '',
  label: '',
  displayMode: 'scalar',
  aggregation: 'sum',
  aggregationKey: '',
  granularity: 'daily',
  timeframe: 'last_30_days',
  comparisonEnabled: false,
  comparisonDisplayType: 'percent',
  chartType: 'area',
  splitBy: '',
  filters: [],
}

function metricToFormValues(metric: Metric): MetricFormValues {
  const validTimeframes = [
    'last_7_days',
    'last_30_days',
    'this_month',
    'last_month',
  ] as const
  const timeframe = validTimeframes.includes(
    metric.timeframe as (typeof validTimeframes)[number]
  )
    ? (metric.timeframe as MetricFormValues['timeframe'])
    : 'last_30_days'

  const filters = (metric.filters ?? []).filter(
    (f): f is { key: string; value: string } => !!f.key && !!f.value
  )

  return {
    dataSourceId: metric.dataSourceId ?? '',
    measurementName: metric.measurementName ?? '',
    label: metric.label ?? '',
    displayMode: metric.displayMode ?? 'scalar',
    aggregation: metric.aggregation ?? 'sum',
    aggregationKey: metric.aggregationKey ?? '',
    granularity: metric.granularity ?? 'daily',
    timeframe,
    comparisonEnabled: metric.comparisonEnabled ?? false,
    comparisonDisplayType: metric.comparisonDisplayType ?? 'percent',
    chartType: metric.chartType ?? 'area',
    splitBy: metric.splitBy ?? '',
    filters,
  }
}

export function formValuesToRequest(
  values: MetricFormValues
): CreateMetricRequest | UpdateMetricRequest {
  return {
    dataSourceId: values.dataSourceId,
    measurementName: values.measurementName,
    label: values.label || values.measurementName,
    displayMode: values.displayMode,
    aggregation: values.aggregation,
    aggregationKey:
      values.aggregation === 'count_unique' ? values.aggregationKey : undefined,
    granularity:
      values.displayMode === 'time_series' ? values.granularity : undefined,
    timeframe: values.timeframe,
    filters: values.filters.length > 0 ? values.filters : undefined,
    comparisonEnabled:
      values.displayMode === 'scalar' ? values.comparisonEnabled : false,
    comparisonDisplayType:
      values.displayMode === 'scalar' && values.comparisonEnabled
        ? values.comparisonDisplayType
        : undefined,
    chartType:
      values.displayMode === 'time_series' ? values.chartType : undefined,
    splitBy:
      values.displayMode === 'time_series' && values.splitBy
        ? values.splitBy
        : undefined,
  }
}

interface UseMetricFormOptions {
  initialValues?: Metric
}

export function useMetricForm({ initialValues }: UseMetricFormOptions = {}) {
  const form = useForm<MetricFormValues>({
    resolver: zodResolver(metricFormSchema),
    defaultValues: initialValues
      ? metricToFormValues(initialValues)
      : defaultValues,
  })

  const dataSourceId = form.watch('dataSourceId')
  const measurementName = form.watch('measurementName')

  // Fetch data sources
  const { data: dataSourcesData } = useGetDataSources()
  const dataSources = dataSourcesData?.dataSources ?? []

  // Fetch measurements for selected data source
  const { data: measurementsData } = useGetDataSourcesDataSourceIdMeasurements(
    dataSourceId,
    { query: { enabled: !!dataSourceId } }
  )
  const measurements = measurementsData?.measurements ?? []

  // Fetch metadata for selected measurement
  const { data: metadataData } =
    useGetDataSourcesDataSourceIdMeasurementsNameMetadata(
      dataSourceId,
      measurementName,
      { query: { enabled: !!dataSourceId && !!measurementName } }
    )

  const metadata = useMemo(() => {
    const keys: string[] = []
    const values: Record<string, string[]> = {}
    for (const item of metadataData?.metadata ?? []) {
      if (item.key) {
        keys.push(item.key)
        values[item.key] = item.values ?? []
      }
    }
    return { keys, values }
  }, [metadataData])

  // Reset dependent fields when data source changes (only if not editing)
  useEffect(() => {
    if (!initialValues && dataSourceId) {
      form.setValue('measurementName', '')
      form.setValue('filters', [])
      form.setValue('aggregationKey', '')
      form.setValue('splitBy', '')
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [dataSourceId])

  const reset = () => {
    form.reset(defaultValues)
  }

  return {
    form,
    dataSources,
    measurements,
    metadata,
    reset,
    isEditing: !!initialValues,
  }
}
