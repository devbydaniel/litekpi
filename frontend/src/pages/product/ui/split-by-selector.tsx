import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import type { MetadataValues } from '@/shared/api/measurements'

interface SplitBySelectorProps {
  metadata: MetadataValues[]
  value: string | undefined
  onChange: (value: string | undefined) => void
}

export function SplitBySelector({ metadata, value, onChange }: SplitBySelectorProps) {
  if (metadata.length === 0) {
    return null
  }

  return (
    <Select
      value={value ?? 'none'}
      onValueChange={(v) => onChange(v === 'none' ? undefined : v)}
    >
      <SelectTrigger className="w-[140px]">
        <SelectValue placeholder="Split by..." />
      </SelectTrigger>
      <SelectContent>
        <SelectItem value="none">No split</SelectItem>
        {metadata.map((meta) => (
          <SelectItem key={meta.key} value={meta.key}>
            {meta.key}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
