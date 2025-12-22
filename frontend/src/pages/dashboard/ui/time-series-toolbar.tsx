import { AreaChart, BarChart3, TrendingUp, Filter } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { Badge } from '@/shared/components/ui/badge'
import { ButtonGroup } from '@/shared/components/ui/button-group'
import type { ChartType, DateRangeValue } from '../hooks/use-time-series-edit'
import type { TimeSeriesMetadata } from '../hooks/use-time-series-metadata'
import type { Filter as FilterType } from '@/shared/api/generated/models'
import { FilterPopover } from './filter-popover'

interface TimeSeriesToolbarProps {
  chartType: ChartType
  dateRange: DateRangeValue
  splitBy: string | undefined
  metadata: TimeSeriesMetadata
  filters: FilterType[]
  onChartTypeChange: (type: ChartType) => void
  onDateRangeChange: (range: DateRangeValue) => void
  onSplitByChange: (key: string | undefined) => void
  onFilterChange: (key: string, value: string | undefined) => void
}

export function TimeSeriesToolbar({
  chartType,
  dateRange,
  splitBy,
  metadata,
  filters,
  onChartTypeChange,
  onDateRangeChange,
  onSplitByChange,
  onFilterChange,
}: TimeSeriesToolbarProps) {
  const metadataKeys = metadata.keys
  const activeFilterCount = filters.length

  return (
    <div className="flex flex-wrap items-center gap-2 px-6 pb-3">
      {/* Chart Type Buttons */}
      <ButtonGroup>
        <Button
          variant={chartType === 'area' ? 'secondary' : 'ghost'}
          size="icon"
          className="h-8 w-8"
          onClick={() => onChartTypeChange('area')}
          title="Area chart"
        >
          <AreaChart className="h-4 w-4" />
        </Button>
        <Button
          variant={chartType === 'bar' ? 'secondary' : 'ghost'}
          size="icon"
          className="h-8 w-8"
          onClick={() => onChartTypeChange('bar')}
          title="Bar chart"
        >
          <BarChart3 className="h-4 w-4" />
        </Button>
        <Button
          variant={chartType === 'line' ? 'secondary' : 'ghost'}
          size="icon"
          className="h-8 w-8"
          onClick={() => onChartTypeChange('line')}
          title="Line chart"
        >
          <TrendingUp className="h-4 w-4" />
        </Button>
      </ButtonGroup>

      <div className="h-6 w-px bg-border" />

      {/* Date Range Selector */}
      <Select
        value={dateRange}
        onValueChange={(v) => onDateRangeChange(v as DateRangeValue)}
      >
        <SelectTrigger className="h-8 w-[180px]">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="last_24_hours">Last 24 hours</SelectItem>
          <SelectItem value="last_7_days">Last 7 days</SelectItem>
          <SelectItem value="last_30_days">Last 30 days</SelectItem>
        </SelectContent>
      </Select>

      {/* Split By Selector */}
      {metadataKeys.length > 0 && (
        <>
          <div className="h-6 w-px bg-border" />
          <Select
            value={splitBy ?? '__none__'}
            onValueChange={(v) =>
              onSplitByChange(v === '__none__' ? undefined : v)
            }
          >
            <SelectTrigger className="h-8 w-[140px]">
              <SelectValue placeholder="Split by..." />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="__none__">No split</SelectItem>
              {metadataKeys.map((key) => (
                <SelectItem key={key} value={key}>
                  {key}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </>
      )}

      {/* Filter Button with Popover */}
      {metadataKeys.length > 0 && (
        <>
          <div className="h-6 w-px bg-border" />
          <FilterPopover
            metadata={metadata}
            filters={filters}
            onFilterChange={onFilterChange}
          >
            <Button variant="ghost" size="sm" className="h-8 gap-1.5">
              <Filter className="h-4 w-4" />
              Filter
              {activeFilterCount > 0 && (
                <Badge variant="secondary" className="ml-1 h-5 min-w-5 px-1.5">
                  {activeFilterCount}
                </Badge>
              )}
            </Button>
          </FilterPopover>
        </>
      )}
    </div>
  )
}
