import { TrendingUp, TrendingDown, Minus } from 'lucide-react'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/shared/components/ui/card'
import { cn } from '@/shared/lib/utils'
import type { ComputedKPI } from '@/shared/api/generated/models'

interface KpiCardProps {
  kpi: ComputedKPI
  className?: string
}

const TIMEFRAME_LABELS: Record<string, string> = {
  last_7_days: 'Last 7 days',
  last_30_days: 'Last 30 days',
  this_month: 'This month',
  last_month: 'Last month',
}

function formatValue(value: number | undefined): string {
  if (value === undefined || value === null) return '-'
  if (Math.abs(value) >= 1000000) {
    return `${(value / 1000000).toFixed(1)}M`
  }
  if (Math.abs(value) >= 1000) {
    return `${(value / 1000).toFixed(1)}K`
  }
  return value.toLocaleString(undefined, { maximumFractionDigits: 2 })
}

function formatChange(
  change: number | undefined,
  changePercent: number | undefined,
  displayType: string | undefined
): string {
  if (change === undefined || change === null) return ''
  if (displayType === 'percent' && changePercent !== undefined) {
    const sign = changePercent > 0 ? '+' : ''
    return `${sign}${changePercent.toFixed(1)}%`
  }
  const sign = change > 0 ? '+' : ''
  return `${sign}${formatValue(change)}`
}

export function KpiCard({ kpi, className }: KpiCardProps) {
  const hasComparison = kpi.comparisonEnabled && kpi.change !== undefined
  const isPositive = hasComparison && (kpi.change ?? 0) > 0
  const isNegative = hasComparison && (kpi.change ?? 0) < 0
  const isNeutral = hasComparison && kpi.change === 0

  const changeText = formatChange(
    kpi.change,
    kpi.changePercent,
    kpi.comparisonDisplayType
  )

  return (
    <Card className={cn('relative', className)}>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm text-muted-foreground">
          {kpi.label}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-1">
          <div className="text-2xl font-bold">{formatValue(kpi.value)}</div>
          {hasComparison && (
            <div
              className={cn(
                'flex items-center gap-1 text-sm',
                isPositive && 'text-emerald-600',
                isNegative && 'text-red-600',
                isNeutral && 'text-muted-foreground'
              )}
            >
              {isPositive && <TrendingUp className="h-4 w-4" />}
              {isNegative && <TrendingDown className="h-4 w-4" />}
              {isNeutral && <Minus className="h-4 w-4" />}
              <span>{changeText}</span>
              <span className="text-muted-foreground">vs previous period</span>
            </div>
          )}
          <div className="text-xs text-muted-foreground">
            {TIMEFRAME_LABELS[kpi.timeframe ?? ''] ?? kpi.timeframe}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

interface KpiCardSkeletonProps {
  className?: string
}

export function KpiCardSkeleton({ className }: KpiCardSkeletonProps) {
  return (
    <Card className={cn('relative', className)}>
      <CardHeader className="pb-2">
        <div className="h-4 w-24 animate-pulse rounded bg-muted" />
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div className="h-8 w-20 animate-pulse rounded bg-muted" />
          <div className="h-4 w-32 animate-pulse rounded bg-muted" />
        </div>
      </CardContent>
    </Card>
  )
}
