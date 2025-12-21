import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'

interface WidgetContextBarProps {
  isDirty: boolean
  isSaving: boolean
  isEditing: boolean
  onSave: () => void
}

export function WidgetContextBar({
  isDirty,
  isSaving,
  isEditing,
  onSave,
}: WidgetContextBarProps) {
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
