import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'

export type DateRangeValue = 'last24h' | 'last7days' | 'last30days'

interface DateRangeFilterProps {
  value: DateRangeValue
  onChange: (value: DateRangeValue) => void
}

const dateRangeOptions: { value: DateRangeValue; label: string }[] = [
  { value: 'last24h', label: 'Last 24 hours' },
  { value: 'last7days', label: 'Last 7 days' },
  { value: 'last30days', label: 'Last 30 days' },
]

export function DateRangeFilter({ value, onChange }: DateRangeFilterProps) {
  return (
    <Select value={value} onValueChange={(v) => onChange(v as DateRangeValue)}>
      <SelectTrigger className="w-[150px]">
        <SelectValue placeholder="Select range" />
      </SelectTrigger>
      <SelectContent>
        {dateRangeOptions.map((option) => (
          <SelectItem key={option.value} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}

export function getDateRangeFromValue(value: DateRangeValue): { start: Date; end: Date } {
  const now = new Date()
  const end = new Date(now)

  let start: Date
  switch (value) {
    case 'last24h':
      start = new Date(now.getTime() - 24 * 60 * 60 * 1000)
      break
    case 'last7days':
      start = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
      break
    case 'last30days':
      start = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000)
      break
    default:
      start = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
  }

  return { start, end }
}
