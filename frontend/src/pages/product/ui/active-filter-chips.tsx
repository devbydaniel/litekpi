import { X } from 'lucide-react'
import { Badge } from '@/shared/components/ui/badge'

interface ActiveFilter {
  key: string
  value: string
}

interface ActiveFilterChipsProps {
  filters: ActiveFilter[]
  onRemove: (key: string) => void
}

export function ActiveFilterChips({ filters, onRemove }: ActiveFilterChipsProps) {
  if (filters.length === 0) {
    return null
  }

  return (
    <div className="flex flex-wrap gap-1">
      {filters.map(({ key, value }) => (
        <Badge
          key={key}
          variant="secondary"
          className="gap-1 pr-1 font-normal"
        >
          <span className="text-muted-foreground">{key}:</span>
          <span>{value}</span>
          <button
            type="button"
            onClick={() => onRemove(key)}
            className="ml-1 rounded-full p-0.5 hover:bg-muted-foreground/20"
            aria-label={`Remove ${key} filter`}
          >
            <X className="h-3 w-3" />
          </button>
        </Badge>
      ))}
    </div>
  )
}
