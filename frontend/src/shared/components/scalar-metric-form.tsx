import { useState, useEffect, useMemo } from 'react'
import { Check, X } from 'lucide-react'
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
import {
  useGetDataSources,
  useGetDataSourcesDataSourceIdMeasurements,
  useGetDataSourcesDataSourceIdMeasurementsNameMetadata,
} from '@/shared/api/generated/api'
import type {
  CreateScalarMetricRequest,
  UpdateScalarMetricRequest,
  ScalarMetric,
  Filter,
  Aggregation,
  ComparisonDisplayType,
} from '@/shared/api/generated/models'

const TIMEFRAME_OPTIONS = [
  { value: 'last_7_days', label: 'Last 7 days' },
  { value: 'last_30_days', label: 'Last 30 days' },
  { value: 'this_month', label: 'This month' },
  { value: 'last_month', label: 'Last month' },
]

const AGGREGATION_OPTIONS = [
  { value: 'sum', label: 'Sum' },
  { value: 'average', label: 'Average' },
]

const COMPARISON_DISPLAY_OPTIONS = [
  { value: 'percent', label: 'Percentage' },
  { value: 'absolute', label: 'Absolute' },
]

export interface ScalarMetricFormValues {
  dataSourceId: string
  measurementName: string
  label: string
  timeframe: string
  aggregation: string
  filters: Filter[]
  comparisonEnabled: boolean
  comparisonDisplayType: string
}

interface ScalarMetricFormProps {
  initialValues?: ScalarMetric
  onSubmit: (values: CreateScalarMetricRequest | UpdateScalarMetricRequest) => Promise<void>
  onCancel: () => void
  isLoading: boolean
  submitLabel?: string
}

export function ScalarMetricForm({
  initialValues,
  onSubmit,
  onCancel,
  isLoading,
  submitLabel = 'Save',
}: ScalarMetricFormProps) {
  const [dataSourceId, setDataSourceId] = useState(initialValues?.dataSourceId ?? '')
  const [measurementName, setMeasurementName] = useState(initialValues?.measurementName ?? '')
  const [label, setLabel] = useState(initialValues?.label ?? '')
  const [timeframe, setTimeframe] = useState<string>(initialValues?.timeframe ?? 'last_30_days')
  const [aggregation, setAggregation] = useState<string>(initialValues?.aggregation ?? 'sum')
  const [filters, setFilters] = useState<Filter[]>(initialValues?.filters ?? [])
  const [comparisonEnabled, setComparisonEnabled] = useState(initialValues?.comparisonEnabled ?? false)
  const [comparisonDisplayType, setComparisonDisplayType] = useState<string>(initialValues?.comparisonDisplayType ?? 'percent')

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
      timeframe,
      aggregation: aggregation as Aggregation,
      filters: filters.length > 0 ? filters : undefined,
      comparisonEnabled,
      comparisonDisplayType: comparisonEnabled ? comparisonDisplayType as ComparisonDisplayType : undefined,
    })
  }

  const canSubmit = dataSourceId && measurementName

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
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

      <div className="space-y-2">
        <Label htmlFor="label">Label</Label>
        <Input
          id="label"
          value={label}
          onChange={(e) => setLabel(e.target.value)}
          placeholder={measurementName || 'Enter a label'}
          maxLength={255}
        />
        <p className="text-xs text-muted-foreground">
          Defaults to measurement name if left empty
        </p>
      </div>

      <div className="grid grid-cols-2 gap-4">
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

        <div className="space-y-2">
          <Label htmlFor="aggregation">Aggregation</Label>
          <Select value={aggregation} onValueChange={setAggregation}>
            <SelectTrigger id="aggregation">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {AGGREGATION_OPTIONS.map((opt) => (
                <SelectItem key={opt.value} value={opt.value}>
                  {opt.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

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

      <div className="space-y-4 rounded-lg border p-4">
        <div className="flex items-center justify-between">
          <div className="space-y-0.5">
            <Label htmlFor="comparison">Compare to previous period</Label>
            <p className="text-xs text-muted-foreground">
              Show change from the previous time period
            </p>
          </div>
          <Switch
            id="comparison"
            checked={comparisonEnabled}
            onCheckedChange={setComparisonEnabled}
          />
        </div>

        {comparisonEnabled && (
          <div className="space-y-2">
            <Label htmlFor="comparisonDisplay">Display type</Label>
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

function FilterSelect({
  filterKey,
  values,
  selectedValue,
  onValueChange,
}: FilterSelectProps) {
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
