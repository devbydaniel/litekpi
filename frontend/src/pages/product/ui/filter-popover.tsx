import { Filter } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { Badge } from '@/shared/components/ui/badge'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shared/components/ui/popover'
import type { MetadataValues } from '@/shared/api/measurements'
import { SearchableFilter } from './searchable-filter'

interface FilterPopoverProps {
  metadata: MetadataValues[]
  filters: Record<string, string>
  activeCount: number
  onFilterChange: (key: string, value: string | undefined) => void
}

export function FilterPopover({
  metadata,
  filters,
  activeCount,
  onFilterChange,
}: FilterPopoverProps) {
  if (metadata.length === 0) {
    return null
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="ghost" size="sm" className="gap-2">
          <Filter className="h-4 w-4" />
          <span>Filters</span>
          {activeCount > 0 && (
            <Badge variant="secondary" className="ml-1 h-5 px-1.5">
              {activeCount}
            </Badge>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-[400px] p-4">
        <div className="grid grid-cols-2 gap-3">
          {metadata.map((meta) => (
            <SearchableFilter
              key={meta.key}
              label={meta.key}
              options={meta.values}
              value={filters[meta.key]}
              onChange={(value) => onFilterChange(meta.key, value)}
            />
          ))}
        </div>
      </PopoverContent>
    </Popover>
  )
}
