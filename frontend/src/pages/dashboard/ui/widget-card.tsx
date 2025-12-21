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
import { useRef } from 'react'
import { MoreHorizontal, Trash, Pencil, Download, Copy } from 'lucide-react'
import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/shared/components/ui/card'
import { Button } from '@/shared/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import { Skeleton } from '@/shared/components/ui/skeleton'
import { Input } from '@/shared/components/ui/input'
import type { Widget, UpdateWidgetRequest } from '@/shared/api/generated/models'
import { useWidgetData } from '../hooks/use-widget-data'
import { useWidgetMetadata } from '../hooks/use-widget-metadata'
import { useWidgetEdit } from '../hooks/use-widget-edit'
import { useChartExport } from '../hooks/use-chart-export'
import { WidgetToolbar } from './widget-toolbar'
import { WidgetContextBar } from './widget-context-bar'

interface WidgetCardProps {
  widget: Widget
  canEdit: boolean
  onDelete: () => void
  onUpdate: (widgetId: string, update: UpdateWidgetRequest) => Promise<void>
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

const chartMargin = { top: 5, right: 10, left: 0, bottom: 5 }

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
        <div className="mb-1 text-xs text-muted-foreground">{label}</div>
        {payload.map(
          (item: { name?: string; value?: number; color?: string }) => (
            <div key={item.name} className="flex items-center gap-2 text-sm">
              <div
                className="h-2 w-2 rounded-full"
                style={{ backgroundColor: item.color }}
              />
              <span className="text-muted-foreground">{item.name}:</span>
              <span className="font-medium">
                {item.value?.toLocaleString() ?? 0}
              </span>
            </div>
          )
        )}
      </div>
    )
  }

  const dataPoint = payload[0]?.payload
  if (!dataPoint) return null

  return (
    <div className="rounded-lg border bg-background p-2 shadow-sm">
      <div className="text-xs text-muted-foreground">{label}</div>
      <div className="font-medium">Sum: {dataPoint.sum?.toLocaleString()}</div>
      <div className="text-xs text-muted-foreground">
        Count: {dataPoint.count?.toLocaleString()}
      </div>
    </div>
  )
}

function ChartEmptyState() {
  return (
    <div className="flex h-[300px] flex-col items-center justify-center gap-2 text-muted-foreground">
      <p className="text-sm">No data available</p>
    </div>
  )
}

export function WidgetCard({
  widget,
  canEdit,
  onDelete,
  onUpdate,
}: WidgetCardProps) {
  // Chart ref for export functionality
  const chartRef = useRef<HTMLDivElement>(null)

  // Edit state management
  const editState = useWidgetEdit({
    widget,
    onSave: onUpdate,
  })

  // Chart export functionality
  const chartExport = useChartExport({
    chartRef,
    filename: widget.title || widget.measurementName || 'chart',
  })

  // Fetch metadata for split-by and filter options
  const { metadata } = useWidgetMetadata(
    widget.dataSourceId ?? '',
    widget.measurementName ?? ''
  )

  // Use preview widget for live updates
  const { data, seriesKeys, isSplit, isLoading, chartType } = useWidgetData(
    editState.previewWidget
  )

  const commonAxisProps = {
    tick: { fontSize: 12 },
    tickLine: false,
    axisLine: false,
  }

  const renderChart = () => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const tooltipContent = (props: any) => (
      <CustomTooltip {...props} isSplit={isSplit} />
    )

    if (chartType === 'line') {
      return (
        <LineChart data={data} margin={chartMargin}>
          <XAxis
            dataKey="date"
            {...commonAxisProps}
            tickFormatter={formatDate}
          />
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
          <XAxis
            dataKey="date"
            {...commonAxisProps}
            tickFormatter={formatDate}
          />
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
      <CardHeader className="px-6 py-4">
        {editState.isEditing ? (
          <Input
            value={editState.state.title ?? ''}
            onChange={(e) => editState.setTitle(e.target.value || undefined)}
            placeholder={widget.measurementName ?? 'Widget title'}
            className="text-base font-medium h-8 max-w-xs"
            maxLength={128}
          />
        ) : (
          <CardTitle className="text-base font-medium">
            {widget.title || widget.measurementName}
          </CardTitle>
        )}
        {canEdit && (
          <CardAction className="flex items-center gap-1">
            <Button
              variant={editState.isEditing ? 'secondary' : 'ghost'}
              size="icon"
              className="h-8 w-8"
              onClick={editState.toggleEditing}
            >
              <Pencil className="h-4 w-4" />
              <span className="sr-only">
                {editState.isEditing ? 'Close edit mode' : 'Edit widget'}
              </span>
            </Button>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8">
                  <Download className="h-4 w-4" />
                  <span className="sr-only">Download chart</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={chartExport.copyToClipboard}>
                  <Copy className="mr-2 h-4 w-4" />
                  Copy to clipboard
                </DropdownMenuItem>
                <DropdownMenuItem onClick={chartExport.downloadAsPng}>
                  <Download className="mr-2 h-4 w-4" />
                  Download as PNG
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-8 w-8">
                  <MoreHorizontal className="h-4 w-4" />
                  <span className="sr-only">Widget options</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem
                  className="text-destructive focus:text-destructive"
                  onClick={onDelete}
                >
                  <Trash className="mr-2 h-4 w-4" />
                  Remove
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </CardAction>
        )}
      </CardHeader>

      {/* Toolbar - shown when editing */}
      {editState.isEditing && (
        <WidgetToolbar
          chartType={editState.state.chartType}
          dateRange={editState.state.dateRange}
          splitBy={editState.state.splitBy}
          metadata={metadata}
          filters={editState.state.filters}
          onChartTypeChange={editState.setChartType}
          onDateRangeChange={editState.setDateRange}
          onSplitByChange={editState.setSplitBy}
          onFilterChange={editState.setFilter}
        />
      )}

      <CardContent className="mt-2">
        {isLoading ? (
          <Skeleton className="h-[300px] w-full" />
        ) : data.length === 0 ? (
          <ChartEmptyState />
        ) : (
          <div ref={chartRef}>
            <ResponsiveContainer width="100%" height={300}>
              {renderChart()}
            </ResponsiveContainer>
          </div>
        )}
      </CardContent>

      {/* Context bar - shown when editing */}
      <WidgetContextBar
        isDirty={editState.isDirty}
        isSaving={editState.isSaving}
        isEditing={editState.isEditing}
        onSave={editState.save}
      />
    </Card>
  )
}
