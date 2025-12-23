import { useState } from 'react'
import { Check, X, Hash, TrendingUp } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/shared/components/ui/form'
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
import type {
  Metric,
  CreateMetricRequest,
  UpdateMetricRequest,
} from '@/shared/api/generated/models'
import {
  useMetricForm,
  formValuesToRequest,
  type MetricFormValues,
} from '../hooks/use-metric-form'

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
  {
    value: 'count_unique',
    label: 'Count Unique',
    description: 'Unique values of a field',
  },
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
  onSubmit: (
    values: CreateMetricRequest | UpdateMetricRequest
  ) => Promise<void>
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
  const { form, dataSources, measurements, metadata, isEditing } =
    useMetricForm({ initialValues })

  // Watch fields needed for conditional rendering
  const displayMode = form.watch('displayMode')
  const aggregation = form.watch('aggregation')
  const comparisonEnabled = form.watch('comparisonEnabled')
  const measurementName = form.watch('measurementName')
  const dataSourceId = form.watch('dataSourceId')
  const filters = form.watch('filters')

  const handleFormSubmit = async (values: MetricFormValues) => {
    await onSubmit(formValuesToRequest(values))
  }

  const handleFilterChange = (key: string, value: string | undefined) => {
    const currentFilters = form.getValues('filters')
    if (value === undefined) {
      form.setValue(
        'filters',
        currentFilters.filter((f) => f.key !== key)
      )
    } else {
      const existingIndex = currentFilters.findIndex((f) => f.key === key)
      if (existingIndex >= 0) {
        const newFilters = [...currentFilters]
        newFilters[existingIndex] = { key, value }
        form.setValue('filters', newFilters)
      } else {
        form.setValue('filters', [...currentFilters, { key, value }])
      }
    }
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(handleFormSubmit)}
        className="space-y-4"
      >
        {/* Data Source */}
        <FormField
          control={form.control}
          name="dataSourceId"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Data Source</FormLabel>
              <Select
                value={field.value}
                onValueChange={field.onChange}
                disabled={isEditing}
              >
                <FormControl>
                  <SelectTrigger>
                    <SelectValue placeholder="Select data source" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {dataSources.map((ds) => (
                    <SelectItem key={ds.id} value={ds.id ?? ''}>
                      {ds.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Measurement */}
        <FormField
          control={form.control}
          name="measurementName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Measurement</FormLabel>
              <Select
                value={field.value}
                onValueChange={field.onChange}
                disabled={!dataSourceId || measurements.length === 0 || isEditing}
              >
                <FormControl>
                  <SelectTrigger>
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
                </FormControl>
                <SelectContent>
                  {measurements.map((m) => (
                    <SelectItem key={m.name} value={m.name ?? ''}>
                      {m.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Label */}
        <FormField
          control={form.control}
          name="label"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Label</FormLabel>
              <FormControl>
                <Input
                  placeholder={measurementName || 'Enter a label'}
                  maxLength={255}
                  {...field}
                />
              </FormControl>
              <FormDescription>
                Defaults to measurement name if left empty
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Aggregation */}
        <FormField
          control={form.control}
          name="aggregation"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Aggregation</FormLabel>
              <Select value={field.value} onValueChange={field.onChange}>
                <FormControl>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  {AGGREGATION_OPTIONS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      <div className="flex flex-col">
                        <span>{opt.label}</span>
                        <span className="text-xs text-muted-foreground">
                          {opt.description}
                        </span>
                      </div>
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Aggregation Key (for count_unique) */}
        {aggregation === 'count_unique' && (
          <FormField
            control={form.control}
            name="aggregationKey"
            render={({ field }) => (
              <FormItem>
                <FormLabel>
                  Count Unique Field <span className="text-destructive">*</span>
                </FormLabel>
                <Select
                  value={field.value}
                  onValueChange={field.onChange}
                  disabled={metadata.keys.length === 0}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue
                        placeholder={
                          metadata.keys.length === 0
                            ? 'No metadata fields available'
                            : 'Select field to count unique values'
                        }
                      />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    {metadata.keys.map((key) => (
                      <SelectItem key={key} value={key}>
                        {key}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
                <FormDescription>
                  The metadata field to count unique values of (e.g., user_id
                  for MAU)
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        )}

        {/* Display Mode */}
        <FormField
          control={form.control}
          name="displayMode"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Display Mode</FormLabel>
              <div className="grid grid-cols-2 gap-3">
                <button
                  type="button"
                  className={cn(
                    'flex flex-col items-center gap-2 rounded-lg border p-3 transition-colors',
                    field.value === 'scalar'
                      ? 'border-primary bg-primary/5'
                      : 'border-border hover:border-muted-foreground'
                  )}
                  onClick={() => field.onChange('scalar')}
                >
                  <Hash className="h-5 w-5" />
                  <div className="text-center">
                    <div className="text-sm font-medium">Scalar</div>
                    <div className="text-xs text-muted-foreground">
                      Single value
                    </div>
                  </div>
                </button>
                <button
                  type="button"
                  className={cn(
                    'flex flex-col items-center gap-2 rounded-lg border p-3 transition-colors',
                    field.value === 'time_series'
                      ? 'border-primary bg-primary/5'
                      : 'border-border hover:border-muted-foreground'
                  )}
                  onClick={() => field.onChange('time_series')}
                >
                  <TrendingUp className="h-5 w-5" />
                  <div className="text-center">
                    <div className="text-sm font-medium">Time Series</div>
                    <div className="text-xs text-muted-foreground">
                      Chart over time
                    </div>
                  </div>
                </button>
              </div>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* Scalar-specific options */}
        {displayMode === 'scalar' && (
          <div className="space-y-4 rounded-lg border p-4">
            <h4 className="text-sm font-medium">Scalar Options</h4>

            <FormField
              control={form.control}
              name="timeframe"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Timeframe</FormLabel>
                  <Select value={field.value} onValueChange={field.onChange}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {TIMEFRAME_OPTIONS.map((opt) => (
                        <SelectItem key={opt.value} value={opt.value}>
                          {opt.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="comparisonEnabled"
              render={({ field }) => (
                <FormItem className="flex items-center justify-between">
                  <div className="space-y-0.5">
                    <FormLabel>Compare to previous period</FormLabel>
                    <FormDescription>
                      Show change from the previous period
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            {comparisonEnabled && (
              <FormField
                control={form.control}
                name="comparisonDisplayType"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Comparison display</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {COMPARISON_DISPLAY_OPTIONS.map((opt) => (
                          <SelectItem key={opt.value} value={opt.value}>
                            {opt.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}
          </div>
        )}

        {/* Time series-specific options */}
        {displayMode === 'time_series' && (
          <div className="space-y-4 rounded-lg border p-4">
            <h4 className="text-sm font-medium">Time Series Options</h4>

            <FormField
              control={form.control}
              name="granularity"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Granularity</FormLabel>
                  <Select
                    value={field.value}
                    onValueChange={field.onChange}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {GRANULARITY_OPTIONS.map((opt) => (
                        <SelectItem key={opt.value} value={opt.value}>
                          {opt.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    How data points are grouped on the chart
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="chartType"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Chart Type</FormLabel>
                  <Select value={field.value} onValueChange={field.onChange}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {CHART_TYPE_OPTIONS.map((opt) => (
                        <SelectItem key={opt.value} value={opt.value}>
                          {opt.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            {metadata.keys.length > 0 && (
              <FormField
                control={form.control}
                name="splitBy"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Split By (optional)</FormLabel>
                    <Select
                      value={field.value || '__none__'}
                      onValueChange={(v) =>
                        field.onChange(v === '__none__' ? '' : v)
                      }
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="No split" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="__none__">No split</SelectItem>
                        {metadata.keys.map((key) => (
                          <SelectItem key={key} value={key}>
                            {key}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
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
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
            disabled={isLoading}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={isLoading}>
            {isLoading ? 'Saving...' : submitLabel}
          </Button>
        </div>
      </form>
    </Form>
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
            <span className="truncate">
              {selectedValue || `All ${filterKey}`}
            </span>
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
