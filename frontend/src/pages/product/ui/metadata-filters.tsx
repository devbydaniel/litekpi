import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import type { MetadataValues } from '@/shared/api/measurements'

interface MetadataFiltersProps {
  metadata: MetadataValues[]
  values: Record<string, string>
  onChange: (key: string, value: string | undefined) => void
}

export function MetadataFilters({ metadata, values, onChange }: MetadataFiltersProps) {
  if (metadata.length === 0) {
    return null
  }

  return (
    <>
      {metadata.map((meta) => (
        <Select
          key={meta.key}
          value={values[meta.key] ?? 'all'}
          onValueChange={(v) => onChange(meta.key, v === 'all' ? undefined : v)}
        >
          <SelectTrigger className="w-[150px]">
            <SelectValue placeholder={meta.key} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All {meta.key}</SelectItem>
            {meta.values.map((value) => (
              <SelectItem key={value} value={value}>
                {value}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      ))}
    </>
  )
}
