import { useNavigate } from '@tanstack/react-router'
import { ChevronDown } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import type { Dashboard } from '@/shared/api/generated/models'

interface DashboardSwitcherProps {
  dashboards: Dashboard[]
  currentDashboard?: Dashboard
}

export function DashboardSwitcher({
  dashboards,
  currentDashboard,
}: DashboardSwitcherProps) {
  const navigate = useNavigate()

  if (dashboards.length <= 1) {
    return null
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="gap-1">
          {currentDashboard?.name ?? 'Dashboard'}
          <ChevronDown className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        {dashboards.map((dashboard) => (
          <DropdownMenuItem
            key={dashboard.id}
            onClick={() => {
              if (dashboard.isDefault) {
                navigate({ to: '/' })
              } else {
                navigate({ to: '/dashboards/$id', params: { id: dashboard.id ?? '' } })
              }
            }}
          >
            {dashboard.name}
            {dashboard.isDefault && (
              <span className="ml-2 text-xs text-muted-foreground">(Default)</span>
            )}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
