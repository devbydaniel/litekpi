import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
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

export function MeasurementChart({ productId, measurement }: MeasurementChartProps) {
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

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-base font-medium">{measurement.name}</CardTitle>
        <div className="flex items-center gap-2">
          <MetadataFilters
            metadata={metadata}
            values={metadataFilters}
            onChange={setMetadataFilter}
          />
          <DateRangeFilter value={dateRange} onChange={setDateRange} />
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <Skeleton className="h-[300px] w-full" />
        ) : data.length === 0 ? (
          <ChartEmptyState />
        ) : (
          <ResponsiveContainer width="100%" height={300}>
            <LineChart
              data={data}
              margin={{
                top: 5,
                right: 10,
                left: 10,
                bottom: 5,
              }}
            >
              <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
              <XAxis
                dataKey="date"
                tick={{ fontSize: 12 }}
                tickLine={false}
                axisLine={false}
                tickFormatter={(value) => {
                  const date = new Date(value)
                  return date.toLocaleDateString('en-US', {
                    month: 'short',
                    day: 'numeric',
                  })
                }}
              />
              <YAxis
                tick={{ fontSize: 12 }}
                tickLine={false}
                axisLine={false}
                tickFormatter={(value) => {
                  if (value >= 1000000) {
                    return `${(value / 1000000).toFixed(1)}M`
                  }
                  if (value >= 1000) {
                    return `${(value / 1000).toFixed(1)}K`
                  }
                  return value.toLocaleString()
                }}
              />
              <Tooltip
                content={({ active, payload, label }) => {
                  if (!active || !payload || payload.length === 0) {
                    return null
                  }
                  const dataPoint = payload[0].payload
                  return (
                    <div className="rounded-lg border bg-background p-2 shadow-sm">
                      <div className="text-xs text-muted-foreground">{label}</div>
                      <div className="font-medium">
                        Sum: {dataPoint.sum.toLocaleString()}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Count: {dataPoint.count.toLocaleString()}
                      </div>
                    </div>
                  )
                }}
              />
              <Line
                type="monotone"
                dataKey="sum"
                stroke="hsl(var(--primary))"
                strokeWidth={2}
                dot={false}
                activeDot={{ r: 4, strokeWidth: 0 }}
              />
            </LineChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  )
}
