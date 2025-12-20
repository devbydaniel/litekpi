import { useState } from 'react'
import {
  BarChart,
  Bar,
  AreaChart,
  Area,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { BarChart3, AreaChart as AreaChartIcon, TrendingUp } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { ButtonGroup } from '@/shared/components/ui/button-group'
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card'
import { Skeleton } from '@/shared/components/ui/skeleton'
import type { MeasurementSummary } from '@/shared/api/measurements'
import { useMeasurementChart } from '../hooks/use-measurement-chart'
import { DateRangeFilter } from './date-range-filter'
import { MetadataFilters } from './metadata-filters'
import { SplitBySelector } from './split-by-selector'
import { ChartEmptyState } from './chart-empty-state'
import { getSeriesColor } from './chart-colors'

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

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const CustomTooltip = ({ active, payload, label, isSplit }: any) => {
  if (!active || !payload || payload.length === 0) return null

  if (isSplit) {
    return (
      <div className="rounded-lg border bg-background p-2 shadow-sm">
        <div className="text-xs text-muted-foreground mb-1">{label}</div>
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
      <div className="font-medium">Sum: {dataPoint.sum?.toLocaleString()}</div>
      <div className="text-xs text-muted-foreground">Count: {dataPoint.count?.toLocaleString()}</div>
    </div>
  )
}

export function MeasurementChart({ productId, measurement }: MeasurementChartProps) {
  const [chartType, setChartType] = useState<'area' | 'bar' | 'line'>('area')
  const {
    data,
    seriesKeys,
    isSplit,
    metadata,
    dateRange,
    metadataFilters,
    splitBy,
    setDateRange,
    setMetadataFilter,
    setSplitBy,
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
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const tooltipContent = (props: any) => <CustomTooltip {...props} isSplit={isSplit} />

    if (chartType === 'line') {
      return (
        <LineChart data={data} margin={chartMargin}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
          <XAxis dataKey="date" {...commonAxisProps} tickFormatter={formatDate} />
          <YAxis {...commonAxisProps} tickFormatter={formatYAxis} />
          <Tooltip content={tooltipContent} />
          {isSplit && <Legend />}
          {isSplit ? (
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
              dataKey="sum"
              stroke="hsl(var(--primary))"
              strokeWidth={2}
              dot={false}
            />
          )}
        </LineChart>
      )
    }

    if (chartType === 'area') {
      return (
        <AreaChart data={data} margin={chartMargin}>
          <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
          <XAxis dataKey="date" {...commonAxisProps} tickFormatter={formatDate} />
          <YAxis {...commonAxisProps} tickFormatter={formatYAxis} />
          <Tooltip content={tooltipContent} />
          {isSplit && <Legend />}
          {isSplit ? (
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
              dataKey="sum"
              stroke="hsl(var(--primary))"
              fill="hsl(var(--primary) / 0.2)"
              strokeWidth={2}
            />
          )}
        </AreaChart>
      )
    }

    // Bar chart (stacked when split)
    return (
      <BarChart data={data} margin={chartMargin}>
        <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
        <XAxis dataKey="date" {...commonAxisProps} tickFormatter={formatDate} />
        <YAxis {...commonAxisProps} tickFormatter={formatYAxis} />
        <Tooltip content={tooltipContent} />
        {isSplit && <Legend />}
        {isSplit ? (
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
          <Bar dataKey="sum" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} />
        )}
      </BarChart>
    )
  }

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 px-6 py-4">
        <CardTitle className="text-base font-medium">{measurement.name}</CardTitle>
        <div className="flex items-center gap-2">
          <SplitBySelector
            metadata={metadata}
            value={splitBy}
            onChange={setSplitBy}
          />
          <MetadataFilters
            metadata={metadata}
            values={metadataFilters}
            onChange={setMetadataFilter}
          />
          <DateRangeFilter value={dateRange} onChange={setDateRange} />
          <ButtonGroup>
            <Button
              variant={chartType === 'area' ? 'secondary' : 'ghost'}
              size="icon"
              className="h-8 w-8"
              onClick={() => setChartType('area')}
            >
              <AreaChartIcon className="h-4 w-4" />
            </Button>
            <Button
              variant={chartType === 'bar' ? 'secondary' : 'ghost'}
              size="icon"
              className="h-8 w-8"
              onClick={() => setChartType('bar')}
            >
              <BarChart3 className="h-4 w-4" />
            </Button>
            <Button
              variant={chartType === 'line' ? 'secondary' : 'ghost'}
              size="icon"
              className="h-8 w-8"
              onClick={() => setChartType('line')}
            >
              <TrendingUp className="h-4 w-4" />
            </Button>
          </ButtonGroup>
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
