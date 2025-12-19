import { BarChart3 } from 'lucide-react'

export function ChartEmptyState() {
  return (
    <div className="flex h-[300px] flex-col items-center justify-center gap-2 text-muted-foreground">
      <BarChart3 className="h-10 w-10" />
      <p className="text-sm">No data for selected filters</p>
    </div>
  )
}
