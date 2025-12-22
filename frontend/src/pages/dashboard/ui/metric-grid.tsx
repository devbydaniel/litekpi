import { useState, useMemo } from 'react'
import { MoreHorizontal, Trash, Pencil, TrendingUp, TrendingDown, Minus } from 'lucide-react'
import {
  BarChart,
  Bar,
  AreaChart,
  Area,
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { Button } from '@/shared/components/ui/button'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardAction,
} from '@/shared/components/ui/card'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog'
import { Skeleton } from '@/shared/components/ui/skeleton'
import { cn } from '@/shared/lib/utils'
import { MetricForm } from '@/widgets/metric-form'
import { DisplayMode } from '@/shared/api/generated/models'
import type { ComputedMetric, Metric, UpdateMetricRequest } from '@/shared/api/generated/models'

interface MetricGridProps {
  metrics: Metric[]
  computedMetrics: ComputedMetric[]
  isLoading: boolean
  canEdit: boolean
  onUpdate: (metricId: string, metric: UpdateMetricRequest) => Promise<void>
  onDelete: (metricId: string) => Promise<void>
  isUpdating: boolean
}

const TIMEFRAME_LABELS: Record<string, string> = {
  last_7_days: 'Last 7 days',
  last_30_days: 'Last 30 days',
  this_month: 'This month',
  last_month: 'Last month',
}

function formatValue(value: number | undefined): string {
  if (value === undefined || value === null) return '-'
  if (Math.abs(value) >= 1000000) {
    return `${(value / 1000000).toFixed(1)}M`
  }
  if (Math.abs(value) >= 1000) {
    return `${(value / 1000).toFixed(1)}K`
  }
  return value.toLocaleString(undefined, { maximumFractionDigits: 2 })
}

function formatChange(
  change: number | undefined,
  changePercent: number | undefined,
  displayType: string | undefined
): string {
  if (change === undefined || change === null) return ''
  if (displayType === 'percent' && changePercent !== undefined) {
    const sign = changePercent > 0 ? '+' : ''
    return `${sign}${changePercent.toFixed(1)}%`
  }
  const sign = change > 0 ? '+' : ''
  return `${sign}${formatValue(change)}`
}

const formatDate = (value: string) => {
  const date = new Date(value)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const formatYAxis = (value: number) => {
  if (value >= 1000000) return `${(value / 1000000).toFixed(1)}M`
  if (value >= 1000) return `${(value / 1000).toFixed(1)}K`
  return value.toLocaleString()
}

const SERIES_COLORS = [
  'hsl(172, 66%, 45%)',
  'hsl(265, 55%, 55%)',
  'hsl(195, 70%, 48%)',
  'hsl(235, 50%, 55%)',
  'hsl(160, 50%, 45%)',
  'hsl(280, 45%, 55%)',
  'hsl(205, 65%, 52%)',
  'hsl(250, 45%, 58%)',
  'hsl(185, 55%, 42%)',
  'hsl(220, 55%, 50%)',
  'hsl(215, 15%, 50%)',
] as const

function getSeriesColor(index: number, key: string): string {
  if (key === 'Other') {
    return SERIES_COLORS[SERIES_COLORS.length - 1]
  }
  return SERIES_COLORS[index % (SERIES_COLORS.length - 1)]
}

export function MetricGrid({
  metrics,
  computedMetrics,
  isLoading,
  canEdit,
  onUpdate,
  onDelete,
  isUpdating,
}: MetricGridProps) {
  const [editingMetric, setEditingMetric] = useState<Metric | null>(null)

  // Separate scalar and time series metrics
  const scalarMetrics = useMemo(
    () => computedMetrics.filter((m) => m.displayMode === DisplayMode.DisplayModeScalar),
    [computedMetrics]
  )
  const timeSeriesMetrics = useMemo(
    () => computedMetrics.filter((m) => m.displayMode === DisplayMode.DisplayModeTimeSeries),
    [computedMetrics]
  )

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
          {[1, 2, 3, 4].map((i) => (
            <ScalarMetricCardSkeleton key={i} />
          ))}
        </div>
        <Skeleton className="h-[400px] w-full" />
      </div>
    )
  }

  if (computedMetrics.length === 0) {
    return null
  }

  const handleEdit = async (values: UpdateMetricRequest) => {
    if (!editingMetric?.id) return
    await onUpdate(editingMetric.id, values)
    setEditingMetric(null)
  }

  return (
    <>
      <div className="space-y-6">
        {/* Scalar metrics in a grid */}
        {scalarMetrics.length > 0 && (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {scalarMetrics.map((metric) => (
              <div key={metric.id} className="group relative">
                <ScalarMetricCard metric={metric} />
                {canEdit && (
                  <MetricActions
                    onEdit={() => {
                      const originalMetric = metrics.find((m) => m.id === metric.id)
                      if (originalMetric) setEditingMetric(originalMetric)
                    }}
                    onDelete={() => metric.id && onDelete(metric.id)}
                  />
                )}
              </div>
            ))}
          </div>
        )}

        {/* Time series metrics as cards */}
        {timeSeriesMetrics.map((metric) => (
          <div key={metric.id} className="group relative">
            <TimeSeriesMetricCard metric={metric} />
            {canEdit && (
              <MetricActions
                onEdit={() => {
                  const originalMetric = metrics.find((m) => m.id === metric.id)
                  if (originalMetric) setEditingMetric(originalMetric)
                }}
                onDelete={() => metric.id && onDelete(metric.id)}
                position="header"
              />
            )}
          </div>
        ))}
      </div>

      <Dialog open={!!editingMetric} onOpenChange={(open) => !open && setEditingMetric(null)}>
        <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Edit Metric</DialogTitle>
          </DialogHeader>
          {editingMetric && (
            <MetricForm
              initialValues={editingMetric}
              onSubmit={handleEdit}
              onCancel={() => setEditingMetric(null)}
              isLoading={isUpdating}
              submitLabel="Save Changes"
            />
          )}
        </DialogContent>
      </Dialog>
    </>
  )
}

interface MetricActionsProps {
  onEdit: () => void
  onDelete: () => void
  position?: 'corner' | 'header'
}

function MetricActions({ onEdit, onDelete, position = 'corner' }: MetricActionsProps) {
  const positionClass =
    position === 'corner'
      ? 'absolute right-2 top-2 opacity-0 transition-opacity group-hover:opacity-100'
      : 'absolute right-4 top-4 opacity-0 transition-opacity group-hover:opacity-100'

  return (
    <div className={positionClass}>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon" className="h-8 w-8">
            <MoreHorizontal className="h-4 w-4" />
            <span className="sr-only">Metric options</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem onClick={onEdit}>
            <Pencil className="mr-2 h-4 w-4" />
            Edit
          </DropdownMenuItem>
          <DropdownMenuItem
            className="text-destructive focus:text-destructive"
            onClick={onDelete}
          >
            <Trash className="mr-2 h-4 w-4" />
            Remove
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}

interface ScalarMetricCardProps {
  metric: ComputedMetric
}

function ScalarMetricCard({ metric }: ScalarMetricCardProps) {
  const hasComparison = metric.comparisonEnabled && metric.change !== undefined
  const isPositive = hasComparison && (metric.change ?? 0) > 0
  const isNegative = hasComparison && (metric.change ?? 0) < 0
  const isNeutral = hasComparison && metric.change === 0

  const changeText = formatChange(
    metric.change,
    metric.changePercent,
    metric.comparisonDisplayType
  )

  return (
    <Card className="relative">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm text-muted-foreground">{metric.label}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-1">
          <div className="text-2xl font-bold">{formatValue(metric.value)}</div>
          {hasComparison && (
            <div
              className={cn(
                'flex items-center gap-1 text-sm',
                isPositive && 'text-emerald-600',
                isNegative && 'text-red-600',
                isNeutral && 'text-muted-foreground'
              )}
            >
              {isPositive && <TrendingUp className="h-4 w-4" />}
              {isNegative && <TrendingDown className="h-4 w-4" />}
              {isNeutral && <Minus className="h-4 w-4" />}
              <span>{changeText}</span>
              <span className="text-muted-foreground">vs previous period</span>
            </div>
          )}
          <div className="text-xs text-muted-foreground">
            {TIMEFRAME_LABELS[metric.timeframe ?? ''] ?? metric.timeframe}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

function ScalarMetricCardSkeleton() {
  return (
    <Card className="relative">
      <CardHeader className="pb-2">
        <div className="h-4 w-24 animate-pulse rounded bg-muted" />
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="h-8 w-20 animate-pulse rounded bg-muted" />
          <div className="h-4 w-32 animate-pulse rounded bg-muted" />
        </div>
      </CardContent>
    </Card>
  )
}

interface TimeSeriesMetricCardProps {
  metric: ComputedMetric
}

function TimeSeriesMetricCard({ metric }: TimeSeriesMetricCardProps) {
  const hasSplit = metric.series && metric.series.length > 0
  const chartType = metric.chartType ?? 'area'

  // Transform data for chart
  const { data, seriesKeys } = useMemo(() => {
    if (hasSplit && metric.series) {
      // Merge all series into a single dataset with each key as a column
      const dateMap = new Map<string, Record<string, string | number>>()
      const keys: string[] = []

      for (const series of metric.series) {
        if (series.key) keys.push(series.key)
        for (const point of series.dataPoints ?? []) {
          if (!point.date) continue
          const existing = dateMap.get(point.date) || { date: point.date }
          if (series.key) {
            existing[series.key] = point.value ?? 0
          }
          dateMap.set(point.date, existing)
        }
      }

      return {
        data: Array.from(dateMap.values()).sort(
          (a, b) => new Date(String(a.date)).getTime() - new Date(String(b.date)).getTime()
        ),
        seriesKeys: keys,
      }
    }

    // Single series data
    return {
      data: (metric.dataPoints ?? []).map((p) => ({
        date: p.date,
        value: p.value ?? 0,
      })),
      seriesKeys: [],
    }
  }, [metric.dataPoints, metric.series, hasSplit])

  const commonAxisProps = {
    tick: { fontSize: 12 },
    tickLine: false,
    axisLine: false,
  }

  const chartMargin = { top: 5, right: 10, left: 0, bottom: 5 }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const CustomTooltip = ({ active, payload, label }: any) => {
    if (!active || !payload || payload.length === 0) return null

    if (hasSplit) {
      return (
        <div className="rounded-lg border bg-background p-2 shadow-sm">
          <div className="mb-1 text-xs text-muted-foreground">{label}</div>
          {payload.map((item: { name?: string; value?: number; color?: string }) => (
            <div key={item.name} className="flex items-center gap-2 text-sm">
              <div
                className="h-2 w-2 rounded-full"
                style={{ backgroundColor: item.color }}
              />
              <span className="text-muted-foreground">{item.name}:</span>
              <span className="font-medium">{item.value?.toLocaleString() ?? 0}</span>
            </div>
          ))}
        </div>
      )
    }

    const dataPoint = payload[0]?.payload
    if (!dataPoint) return null

    return (
      <div className="rounded-lg border bg-background p-2 shadow-sm">
        <div className="text-xs text-muted-foreground">{label}</div>
        <div className="font-medium">Value: {dataPoint.value?.toLocaleString()}</div>
      </div>
    )
  }

  const renderChart = () => {
    if (data.length === 0) {
      return (
        <div className="flex h-[300px] flex-col items-center justify-center gap-2 text-muted-foreground">
          <p className="text-sm">No data available</p>
        </div>
      )
    }

    if (chartType === 'line') {
      return (
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={data} margin={chartMargin}>
            <XAxis dataKey="date" {...commonAxisProps} tickFormatter={formatDate} />
            <YAxis {...commonAxisProps} tickFormatter={formatYAxis} />
            <Tooltip content={CustomTooltip} />
            {hasSplit && <Legend />}
            {hasSplit ? (
              seriesKeys.map((key, i) => (
                <Line
                  key={key}
                  type="monotone"
                  dataKey={key}
                  stroke={getSeriesColor(i, key)}
                  strokeWidth={2}
                  dot={false}
                />
              ))
            ) : (
              <Line
                type="monotone"
                dataKey="value"
                stroke="hsl(var(--primary))"
                strokeWidth={2}
                dot={false}
              />
            )}
          </LineChart>
        </ResponsiveContainer>
      )
    }

    if (chartType === 'area') {
      return (
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={data} margin={chartMargin}>
            <XAxis dataKey="date" {...commonAxisProps} tickFormatter={formatDate} />
            <YAxis {...commonAxisProps} tickFormatter={formatYAxis} />
            <Tooltip content={CustomTooltip} />
            {hasSplit && <Legend />}
            {hasSplit ? (
              seriesKeys.map((key, i) => (
                <Area
                  key={key}
                  type="monotone"
                  dataKey={key}
                  stackId="1"
                  stroke={getSeriesColor(i, key)}
                  fill={getSeriesColor(i, key)}
                  fillOpacity={0.6}
                />
              ))
            ) : (
              <Area
                type="monotone"
                dataKey="value"
                stroke="hsl(var(--primary))"
                fill="hsl(var(--primary) / 0.2)"
                strokeWidth={2}
              />
            )}
          </AreaChart>
        </ResponsiveContainer>
      )
    }

    // Bar chart
    return (
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={data} margin={chartMargin}>
          <XAxis dataKey="date" {...commonAxisProps} tickFormatter={formatDate} />
          <YAxis {...commonAxisProps} tickFormatter={formatYAxis} />
          <Tooltip content={CustomTooltip} />
          {hasSplit && <Legend />}
          {hasSplit ? (
            seriesKeys.map((key, i) => (
              <Bar
                key={key}
                dataKey={key}
                stackId="1"
                fill={getSeriesColor(i, key)}
                radius={i === seriesKeys.length - 1 ? [4, 4, 0, 0] : [0, 0, 0, 0]}
              />
            ))
          ) : (
            <Bar dataKey="value" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} />
          )}
        </BarChart>
      </ResponsiveContainer>
    )
  }

  return (
    <Card>
      <CardHeader className="px-6 py-4">
        <CardTitle className="text-base font-medium">
          {metric.label || metric.measurementName}
        </CardTitle>
        <CardAction />
      </CardHeader>
      <CardContent className="mt-2">{renderChart()}</CardContent>
    </Card>
  )
}
