import { useState } from 'react'
import { Check, ChevronsUpDown, X } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { Badge } from '@/shared/components/ui/badge'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/shared/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/shared/components/ui/popover'
import { cn } from '@/shared/lib/utils'
import type { DataSource } from '@/shared/api/generated/models'

interface DataSourceMultiSelectProps {
  dataSources: DataSource[]
  selectedIds: string[]
  onChange: (ids: string[]) => void
  disabled?: boolean
}

export function DataSourceMultiSelect({
  dataSources,
  selectedIds,
  onChange,
  disabled,
}: DataSourceMultiSelectProps) {
  const [open, setOpen] = useState(false)

  const selectedDataSources = dataSources.filter(
    (ds) => ds.id && selectedIds.includes(ds.id)
  )

  const toggleDataSource = (id: string) => {
    if (selectedIds.includes(id)) {
      onChange(selectedIds.filter((selectedId) => selectedId !== id))
    } else {
      onChange([...selectedIds, id])
    }
  }

  const removeDataSource = (id: string) => {
    onChange(selectedIds.filter((selectedId) => selectedId !== id))
  }

  return (
    <div className="space-y-2">
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between"
            disabled={disabled}
          >
            {selectedIds.length === 0
              ? 'Select data sources...'
              : `${selectedIds.length} selected`}
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-full p-0" align="start">
          <Command>
            <CommandInput placeholder="Search data sources..." />
            <CommandList>
              <CommandEmpty>No data sources found.</CommandEmpty>
              <CommandGroup>
                {dataSources.map((ds) => (
                  <CommandItem
                    key={ds.id}
                    value={ds.name}
                    onSelect={() => ds.id && toggleDataSource(ds.id)}
                  >
                    <Check
                      className={cn(
                        'mr-2 h-4 w-4',
                        ds.id && selectedIds.includes(ds.id)
                          ? 'opacity-100'
                          : 'opacity-0'
                      )}
                    />
                    {ds.name}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      {selectedDataSources.length > 0 && (
        <div className="flex flex-wrap gap-1">
          {selectedDataSources.map((ds) => (
            <Badge key={ds.id} variant="secondary" className="gap-1">
              {ds.name}
              <button
                type="button"
                onClick={() => ds.id && removeDataSource(ds.id)}
                className="ml-1 rounded-full hover:bg-muted"
                disabled={disabled}
              >
                <X className="h-3 w-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  )
}
