import { useState } from 'react'
import {
  BarChart,
  Bar,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import { BarChart3, AreaChart as AreaChartIcon } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card'
import { Skeleton } from '@/shared/components/ui/skeleton'
import type { MeasurementSummary } from '@/shared/api/measurements'
import { useMeasurementChart } from '../hooks/use-measurement-chart'
import { DateRangeFilter } from './date-range-filter'
import { MetadataFilters } from './metadata-filters'
import { ChartEmptyState } from './chart-empty-state'

interface MeasurementChartProps {
  productId: string
  measurement: MeasurementSummary
}

const chartMargin = { top: 5, right: 10, left: 10, bottom: 5 }

const formatDate = (value: string) => {
  const date = new Date(value)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const formatYAxis = (value: number) => {
  if (value >= 1000000) return `${(value / 1000000).toFixed(1)}M`
  if (value >= 1000) return `${(value / 1000).toFixed(1)}K`
  return value.toLocaleString()
}

const renderTooltip = ({ active, payload, label }: { active?: boolean; payload?: Array<{ payload?: { sum: number; count: number } }>; label?: string }) => {
  if (!active || !payload || payload.length === 0 || !payload[0].payload) return null
  const dataPoint = payload[0].payload
  return (
    <div className="rounded-lg border bg-background p-2 shadow-sm">
      <div className="text-xs text-muted-foreground">{label}</div>
      <div className="font-medium">Sum: {dataPoint.sum.toLocaleString()}</div>
      <div className="text-xs text-muted-foreground">Count: {dataPoint.count.toLocaleString()}</div>
    </div>
  )
}

export function MeasurementChart({ productId, measurement }: MeasurementChartProps) {
  const [chartType, setChartType] = useState<'area' | 'bar'>('area')
  const {
    data,
    metadata,
    dateRange,
    metadataFilters,
    setDateRange,
    setMetadataFilter,
    isLoading,
  } = useMeasurementChart({
    productId,
    measurementName: measurement.name,
  })

  const commonAxisProps = {
    tick: { fontSize: 12 },
    tickLine: false,
    axisLine: false,
  }

  const renderChart = () => {
    if (chartType === 'area') {
      return (
        <AreaChart data={data} margin={chartMargin}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
          <XAxis dataKey="date" {...commonAxisProps} tickFormatter={formatDate} />
          <YAxis {...commonAxisProps} tickFormatter={formatYAxis} />
          <Tooltip content={renderTooltip} />
          <Area type="monotone" dataKey="sum" stroke="hsl(var(--primary))" fill="hsl(var(--primary) / 0.2)" strokeWidth={2} />
        </AreaChart>
      )
    }
    return (
      <BarChart data={data} margin={chartMargin}>
        <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
        <XAxis dataKey="date" {...commonAxisProps} tickFormatter={formatDate} />
        <YAxis {...commonAxisProps} tickFormatter={formatYAxis} />
        <Tooltip content={renderTooltip} />
        <Bar dataKey="sum" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} />
      </BarChart>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 px-6 py-4">
        <CardTitle className="text-base font-medium">{measurement.name}</CardTitle>
        <div className="flex items-center gap-2">
          <MetadataFilters
            metadata={metadata}
            values={metadataFilters}
            onChange={setMetadataFilter}
          />
          <DateRangeFilter value={dateRange} onChange={setDateRange} />
          <div className="flex">
            <Button
              variant={chartType === 'area' ? 'secondary' : 'ghost'}
              size="icon"
              className="h-8 w-8 rounded-r-none"
              onClick={() => setChartType('area')}
            >
              <AreaChartIcon className="h-4 w-4" />
            </Button>
            <Button
              variant={chartType === 'bar' ? 'secondary' : 'ghost'}
              size="icon"
              className="h-8 w-8 rounded-l-none"
              onClick={() => setChartType('bar')}
            >
              <BarChart3 className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <Skeleton className="h-[300px] w-full" />
        ) : data.length === 0 ? (
          <ChartEmptyState />
        ) : (
          <ResponsiveContainer width="100%" height={300}>
            {renderChart()}
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  )
}
