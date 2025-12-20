import { BarChart3, AreaChart as AreaChartIcon, TrendingUp } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { ButtonGroup } from '@/shared/components/ui/button-group'
import type { MetadataValues } from '@/shared/api/generated/models'
import type { ChartType } from '../hooks/use-measurement-chart'
import type { DateRangeValue } from './date-range-filter'
import { DateRangeFilter } from './date-range-filter'
import { SplitBySelector } from './split-by-selector'
import { FilterPopover } from './filter-popover'

interface ChartToolbarProps {
  metadata: MetadataValues[]
  chartType: ChartType
  dateRange: DateRangeValue
  splitBy: string | undefined
  filters: Record<string, string>
  activeFilterCount: number
  onChartTypeChange: (type: ChartType) => void
  onDateRangeChange: (range: DateRangeValue) => void
  onSplitByChange: (splitBy: string | undefined) => void
  onFilterChange: (key: string, value: string | undefined) => void
}

export function ChartToolbar({
  metadata,
  chartType,
  dateRange,
  splitBy,
  filters,
  activeFilterCount,
  onChartTypeChange,
  onDateRangeChange,
  onSplitByChange,
  onFilterChange,
}: ChartToolbarProps) {
  return (
    <div className="flex flex-wrap items-center gap-2 px-6 py-2">
      <DateRangeFilter value={dateRange} onChange={onDateRangeChange} />
      <SplitBySelector
        metadata={metadata}
        value={splitBy}
        onChange={onSplitByChange}
      />
      <FilterPopover
        metadata={metadata}
        filters={filters}
        activeCount={activeFilterCount}
        onFilterChange={onFilterChange}
      />
      <div className="flex-1" />
      <ButtonGroup>
        <Button
          variant={chartType === 'area' ? 'secondary' : 'ghost'}
          size="icon"
          className="h-8 w-8"
          onClick={() => onChartTypeChange('area')}
        >
          <AreaChartIcon className="h-4 w-4" />
        </Button>
        <Button
          variant={chartType === 'bar' ? 'secondary' : 'ghost'}
          size="icon"
          className="h-8 w-8"
          onClick={() => onChartTypeChange('bar')}
        >
          <BarChart3 className="h-4 w-4" />
        </Button>
        <Button
          variant={chartType === 'line' ? 'secondary' : 'ghost'}
          size="icon"
          className="h-8 w-8"
          onClick={() => onChartTypeChange('line')}
        >
          <TrendingUp className="h-4 w-4" />
        </Button>
      </ButtonGroup>
    </div>
  )
}
