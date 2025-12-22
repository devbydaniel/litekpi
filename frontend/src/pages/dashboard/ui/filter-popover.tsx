import { useState } from 'react'
import { Check, X } from 'lucide-react'
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
import { Button } from '@/shared/components/ui/button'
import { Label } from '@/shared/components/ui/label'
import type { Filter } from '@/shared/api/generated/models'
import type { TimeSeriesMetadata } from '../hooks/use-time-series-metadata'

interface FilterPopoverProps {
  metadata: TimeSeriesMetadata
  filters: Filter[]
  onFilterChange: (key: string, value: string | undefined) => void
  children: React.ReactNode
}

export function FilterPopover({
  metadata,
  filters,
  onFilterChange,
  children,
}: FilterPopoverProps) {
  const filterMap = new Map(
    filters.map((f) => [f.key ?? '', f.value ?? ''])
  )

  return (
    <Popover>
      <PopoverTrigger asChild>{children}</PopoverTrigger>
      <PopoverContent className="w-80 p-0" align="start">
        <div className="border-b px-3 py-2">
          <h4 className="text-sm font-medium">Filters</h4>
        </div>
        <div className="max-h-[400px] overflow-y-auto p-3">
          <div className="grid gap-4">
            {metadata.keys.map((key) => (
              <SearchableFilter
                key={key}
                filterKey={key}
                values={metadata.values[key] ?? []}
                selectedValue={filterMap.get(key)}
                onValueChange={(value) => onFilterChange(key, value)}
              />
            ))}
            {metadata.keys.length === 0 && (
              <p className="text-sm text-muted-foreground">
                No metadata available for filtering.
              </p>
            )}
          </div>
        </div>
      </PopoverContent>
    </Popover>
  )
}

interface SearchableFilterProps {
  filterKey: string
  values: string[]
  selectedValue: string | undefined
  onValueChange: (value: string | undefined) => void
}

function SearchableFilter({
  filterKey,
  values,
  selectedValue,
  onValueChange,
}: SearchableFilterProps) {
  const [open, setOpen] = useState(false)

  return (
    <div className="space-y-1.5">
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
              {selectedValue || `Select ${filterKey}...`}
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
