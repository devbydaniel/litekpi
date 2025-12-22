import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'

interface TimeSeriesContextBarProps {
  isDirty: boolean
  isSaving: boolean
  isEditing: boolean
  onSave: () => void
}

export function TimeSeriesContextBar({
  isDirty,
  isSaving,
  isEditing,
  onSave,
}: TimeSeriesContextBarProps) {
  if (!isEditing) return null

  return (
    <div className="flex items-center justify-end px-6 py-3">
      <Button
        size="sm"
        onClick={onSave}
        disabled={isSaving || !isDirty}
      >
        {isSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Save
      </Button>
    </div>
  )
}
