import { useState, useEffect, useMemo } from 'react'
import { Check, X, Hash, TrendingUp } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shared/components/ui/popover'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/shared/components/ui/command'
import { Switch } from '@/shared/components/ui/switch'
import { cn } from '@/shared/lib/utils'
import {
  useGetDataSources,
  useGetDataSourcesDataSourceIdMeasurements,
  useGetDataSourcesDataSourceIdMeasurementsNameMetadata,
} from '@/shared/api/generated/api'
import type {
  CreateMetricRequest,
  UpdateMetricRequest,
  Metric,
  Filter,
  Aggregation,
  ComparisonDisplayType,
  DisplayMode,
  Granularity,
  ChartType,
} from '@/shared/api/generated/models'

const TIMEFRAME_OPTIONS = [
  { value: 'last_7_days', label: 'Last 7 days' },
  { value: 'last_30_days', label: 'Last 30 days' },
  { value: 'this_month', label: 'This month' },
  { value: 'last_month', label: 'Last month' },
]

const AGGREGATION_OPTIONS = [
  { value: 'sum', label: 'Sum', description: 'Total of all values' },
  { value: 'average', label: 'Average', description: 'Mean of all values' },
  { value: 'count', label: 'Count', description: 'Number of data points' },
  { value: 'count_unique', label: 'Count Unique', description: 'Unique values of a field' },
]

const GRANULARITY_OPTIONS = [
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: 'Monthly' },
]

const CHART_TYPE_OPTIONS = [
  { value: 'area', label: 'Area' },
  { value: 'bar', label: 'Bar' },
  { value: 'line', label: 'Line' },
]

const COMPARISON_DISPLAY_OPTIONS = [
  { value: 'percent', label: 'Percentage' },
  { value: 'absolute', label: 'Absolute' },
]

interface MetricFormProps {
  initialValues?: Metric
  onSubmit: (values: CreateMetricRequest | UpdateMetricRequest) => Promise<void>
  onCancel: () => void
  isLoading: boolean
  submitLabel?: string
}

export function MetricForm({
  initialValues,
  onSubmit,
  onCancel,
  isLoading,
  submitLabel = 'Save',
}: MetricFormProps) {
  // Data source and measurement
  const [dataSourceId, setDataSourceId] = useState(initialValues?.dataSourceId ?? '')
  const [measurementName, setMeasurementName] = useState(initialValues?.measurementName ?? '')
  const [label, setLabel] = useState(initialValues?.label ?? '')

  // Display mode
  const [displayMode, setDisplayMode] = useState<string>(initialValues?.displayMode ?? 'scalar')

  // Aggregation
  const [aggregation, setAggregation] = useState<string>(initialValues?.aggregation ?? 'sum')
  const [aggregationKey, setAggregationKey] = useState(initialValues?.aggregationKey ?? '')
  const [granularity, setGranularity] = useState<string>(initialValues?.granularity ?? 'daily')

  // Scalar options
  const [timeframe, setTimeframe] = useState<string>(initialValues?.timeframe ?? 'last_30_days')
  const [comparisonEnabled, setComparisonEnabled] = useState(initialValues?.comparisonEnabled ?? false)
  const [comparisonDisplayType, setComparisonDisplayType] = useState<string>(
    initialValues?.comparisonDisplayType ?? 'percent'
  )

  // Time series options
  const [chartType, setChartType] = useState<string>(initialValues?.chartType ?? 'area')
  const [splitBy, setSplitBy] = useState<string>(initialValues?.splitBy ?? '')

  // Filters
  const [filters, setFilters] = useState<Filter[]>(initialValues?.filters ?? [])

  // Fetch data sources
  const { data: dataSourcesData } = useGetDataSources()
  const dataSources = dataSourcesData?.dataSources ?? []

  // Fetch measurements for selected data source
  const { data: measurementsData } = useGetDataSourcesDataSourceIdMeasurements(dataSourceId, {
    query: { enabled: !!dataSourceId },
  })
  const measurements = measurementsData?.measurements ?? []

  // Fetch metadata for selected measurement
  const { data: metadataData } = useGetDataSourcesDataSourceIdMeasurementsNameMetadata(
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

  // Reset measurement when data source changes (only if not editing)
  useEffect(() => {
    if (!initialValues) {
      setMeasurementName('')
      setFilters([])
      setAggregationKey('')
      setSplitBy('')
    }
  }, [dataSourceId, initialValues])

  const handleFilterChange = (key: string, value: string | undefined) => {
    if (value === undefined) {
      setFilters(filters.filter((f) => f.key !== key))
    } else {
      const existingIndex = filters.findIndex((f) => f.key === key)
      if (existingIndex >= 0) {
        const newFilters = [...filters]
        newFilters[existingIndex] = { key, value }
        setFilters(newFilters)
      } else {
        setFilters([...filters, { key, value }])
      }
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    await onSubmit({
      dataSourceId,
      measurementName,
      label: label || measurementName,
      displayMode: displayMode as DisplayMode,
      aggregation: aggregation as Aggregation,
      aggregationKey: aggregation === 'count_unique' ? aggregationKey : undefined,
      granularity: granularity as Granularity,
      timeframe,
      filters: filters.length > 0 ? filters : undefined,
      comparisonEnabled: displayMode === 'scalar' ? comparisonEnabled : false,
      comparisonDisplayType:
        displayMode === 'scalar' && comparisonEnabled
          ? (comparisonDisplayType as ComparisonDisplayType)
          : undefined,
      chartType: displayMode === 'time_series' ? (chartType as ChartType) : undefined,
      splitBy: displayMode === 'time_series' && splitBy ? splitBy : undefined,
    })
  }

  const showAggregationKey = aggregation === 'count_unique'
  const canSubmit =
    dataSourceId && measurementName && (!showAggregationKey || aggregationKey)

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Data Source */}
      <div className="space-y-2">
        <Label htmlFor="dataSource">Data Source</Label>
        <Select value={dataSourceId} onValueChange={setDataSourceId} disabled={!!initialValues}>
          <SelectTrigger id="dataSource">
            <SelectValue placeholder="Select data source" />
          </SelectTrigger>
          <SelectContent>
            {dataSources.map((ds) => (
              <SelectItem key={ds.id} value={ds.id ?? ''}>
                {ds.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Measurement */}
      <div className="space-y-2">
        <Label htmlFor="measurement">Measurement</Label>
        <Select
          value={measurementName}
          onValueChange={setMeasurementName}
          disabled={!dataSourceId || measurements.length === 0 || !!initialValues}
        >
          <SelectTrigger id="measurement">
            <SelectValue
              placeholder={
                !dataSourceId
                  ? 'Select a data source first'
                  : measurements.length === 0
                    ? 'No measurements available'
                    : 'Select measurement'
              }
            />
          </SelectTrigger>
          <SelectContent>
            {measurements.map((m) => (
              <SelectItem key={m.name} value={m.name ?? ''}>
                {m.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Label */}
      <div className="space-y-2">
        <Label htmlFor="label">Label</Label>
        <Input
          id="label"
          value={label}
          onChange={(e) => setLabel(e.target.value)}
          placeholder={measurementName || 'Enter a label'}
          maxLength={255}
        />
        <p className="text-xs text-muted-foreground">Defaults to measurement name if left empty</p>
      </div>

      {/* Display Mode */}
      <div className="space-y-2">
        <Label>Display Mode</Label>
        <div className="grid grid-cols-2 gap-3">
          <button
            type="button"
            className={cn(
              'flex flex-col items-center gap-2 rounded-lg border p-3 transition-colors',
              displayMode === 'scalar'
                ? 'border-primary bg-primary/5'
                : 'border-border hover:border-muted-foreground'
            )}
            onClick={() => setDisplayMode('scalar')}
          >
            <Hash className="h-5 w-5" />
            <div className="text-center">
              <div className="text-sm font-medium">Scalar</div>
              <div className="text-xs text-muted-foreground">Single value</div>
            </div>
          </button>
          <button
            type="button"
            className={cn(
              'flex flex-col items-center gap-2 rounded-lg border p-3 transition-colors',
              displayMode === 'time_series'
                ? 'border-primary bg-primary/5'
                : 'border-border hover:border-muted-foreground'
            )}
            onClick={() => setDisplayMode('time_series')}
          >
            <TrendingUp className="h-5 w-5" />
            <div className="text-center">
              <div className="text-sm font-medium">Time Series</div>
              <div className="text-xs text-muted-foreground">Chart over time</div>
            </div>
          </button>
        </div>
      </div>

      {/* Aggregation */}
      <div className="space-y-2">
        <Label htmlFor="aggregation">Aggregation</Label>
        <Select value={aggregation} onValueChange={setAggregation}>
          <SelectTrigger id="aggregation">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {AGGREGATION_OPTIONS.map((opt) => (
              <SelectItem key={opt.value} value={opt.value}>
                <div className="flex flex-col">
                  <span>{opt.label}</span>
                  <span className="text-xs text-muted-foreground">{opt.description}</span>
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Aggregation Key (for count_unique) */}
      {showAggregationKey && (
        <div className="space-y-2">
          <Label htmlFor="aggregationKey">
            Count Unique Field <span className="text-destructive">*</span>
          </Label>
          <Select
            value={aggregationKey}
            onValueChange={setAggregationKey}
            disabled={metadata.keys.length === 0}
          >
            <SelectTrigger id="aggregationKey">
              <SelectValue
                placeholder={
                  metadata.keys.length === 0
                    ? 'No metadata fields available'
                    : 'Select field to count unique values'
                }
              />
            </SelectTrigger>
            <SelectContent>
              {metadata.keys.map((key) => (
                <SelectItem key={key} value={key}>
                  {key}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <p className="text-xs text-muted-foreground">
            The metadata field to count unique values of (e.g., user_id for MAU)
          </p>
        </div>
      )}

      {/* Granularity */}
      <div className="space-y-2">
        <Label htmlFor="granularity">Granularity</Label>
        <Select value={granularity} onValueChange={setGranularity}>
          <SelectTrigger id="granularity">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {GRANULARITY_OPTIONS.map((opt) => (
              <SelectItem key={opt.value} value={opt.value}>
                {opt.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Scalar-specific options */}
      {displayMode === 'scalar' && (
        <div className="space-y-4 rounded-lg border p-4">
          <h4 className="text-sm font-medium">Scalar Options</h4>

          <div className="space-y-2">
            <Label htmlFor="timeframe">Timeframe</Label>
            <Select value={timeframe} onValueChange={setTimeframe}>
              <SelectTrigger id="timeframe">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {TIMEFRAME_OPTIONS.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-center justify-between">
            <div className="space-y-0.5">
              <Label htmlFor="comparison">Compare to previous period</Label>
              <p className="text-xs text-muted-foreground">Show change from the previous period</p>
            </div>
            <Switch
              id="comparison"
              checked={comparisonEnabled}
              onCheckedChange={setComparisonEnabled}
            />
          </div>

          {comparisonEnabled && (
            <div className="space-y-2">
              <Label htmlFor="comparisonDisplay">Comparison display</Label>
              <Select value={comparisonDisplayType} onValueChange={setComparisonDisplayType}>
                <SelectTrigger id="comparisonDisplay">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {COMPARISON_DISPLAY_OPTIONS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}
        </div>
      )}

      {/* Time series-specific options */}
      {displayMode === 'time_series' && (
        <div className="space-y-4 rounded-lg border p-4">
          <h4 className="text-sm font-medium">Chart Options</h4>

          <div className="space-y-2">
            <Label htmlFor="chartType">Chart Type</Label>
            <Select value={chartType} onValueChange={setChartType}>
              <SelectTrigger id="chartType">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {CHART_TYPE_OPTIONS.map((opt) => (
                  <SelectItem key={opt.value} value={opt.value}>
                    {opt.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {metadata.keys.length > 0 && (
            <div className="space-y-2">
              <Label htmlFor="splitBy">Split By (optional)</Label>
              <Select
                value={splitBy || '__none__'}
                onValueChange={(v) => setSplitBy(v === '__none__' ? '' : v)}
              >
                <SelectTrigger id="splitBy">
                  <SelectValue placeholder="No split" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="__none__">No split</SelectItem>
                  {metadata.keys.map((key) => (
                    <SelectItem key={key} value={key}>
                      {key}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}
        </div>
      )}

      {/* Filters */}
      {metadata.keys.length > 0 && (
        <div className="space-y-2">
          <Label>Filters</Label>
          <div className="grid gap-3 rounded-lg border p-3">
            {metadata.keys.map((key) => (
              <FilterSelect
                key={key}
                filterKey={key}
                values={metadata.values[key] ?? []}
                selectedValue={filters.find((f) => f.key === key)?.value}
                onValueChange={(value) => handleFilterChange(key, value)}
              />
            ))}
          </div>
        </div>
      )}

      {/* Actions */}
      <div className="flex justify-end gap-2 pt-2">
        <Button type="button" variant="outline" onClick={onCancel} disabled={isLoading}>
          Cancel
        </Button>
        <Button type="submit" disabled={isLoading || !canSubmit}>
          {isLoading ? 'Saving...' : submitLabel}
        </Button>
      </div>
    </form>
  )
}

interface FilterSelectProps {
  filterKey: string
  values: string[]
  selectedValue: string | undefined
  onValueChange: (value: string | undefined) => void
}

function FilterSelect({ filterKey, values, selectedValue, onValueChange }: FilterSelectProps) {
  const [open, setOpen] = useState(false)

  return (
    <div className="space-y-1">
      <Label className="text-xs text-muted-foreground">{filterKey}</Label>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="h-8 w-full justify-between"
          >
            <span className="truncate">{selectedValue || `All ${filterKey}`}</span>
            {selectedValue && (
              <X
                className="ml-2 h-3 w-3 shrink-0 opacity-50 hover:opacity-100"
                onClick={(e) => {
                  e.stopPropagation()
                  onValueChange(undefined)
                }}
              />
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0" align="start">
          <Command>
            <CommandInput placeholder={`Search ${filterKey}...`} />
            <CommandList>
              <CommandEmpty>No results found.</CommandEmpty>
              <CommandGroup>
                {values.map((value) => (
                  <CommandItem
                    key={value}
                    value={value}
                    onSelect={() => {
                      onValueChange(value === selectedValue ? undefined : value)
                      setOpen(false)
                    }}
                  >
                    <Check
                      className={`mr-2 h-4 w-4 ${
                        value === selectedValue ? 'opacity-100' : 'opacity-0'
                      }`}
                    />
                    {value}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  )
}
