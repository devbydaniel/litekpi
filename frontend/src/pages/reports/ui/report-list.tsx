import { Link } from '@tanstack/react-router'
import { ChevronRight, FileText } from 'lucide-react'
import { EmptyState } from '@/shared/components/ui/empty-state'
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemGroup,
  ItemSeparator,
  ItemTitle,
} from '@/shared/components/ui/item'
import { Skeleton } from '@/shared/components/ui/skeleton'
import type { Report } from '@/shared/api/generated/models'

interface ReportListProps {
  reports: Report[]
  isLoading: boolean
}

export function ReportList({ reports, isLoading }: ReportListProps) {
  if (isLoading) {
    return <ReportListSkeleton />
  }

  if (reports.length === 0) {
    return (
      <EmptyState
        icon={FileText}
        title="No reports yet"
        description="Create your first report to track KPI metrics."
      />
    )
  }

  return (
    <ItemGroup className="rounded-lg border">
      {reports.map((report, index) => (
        <div key={report.id}>
          {index > 0 && <ItemSeparator />}
          <Link to="/reports/$id" params={{ id: report.id! }}>
            <Item className="cursor-pointer">
              <ItemContent>
                <ItemTitle>{report.name}</ItemTitle>
                <ItemDescription>
                  Created {report.createdAt ? new Date(report.createdAt).toLocaleDateString() : '-'}
                </ItemDescription>
              </ItemContent>
              <ItemActions>
                <ChevronRight className="h-5 w-5 text-muted-foreground" />
              </ItemActions>
            </Item>
          </Link>
        </div>
      ))}
    </ItemGroup>
  )
}

function ReportListSkeleton() {
  return (
    <ItemGroup className="rounded-lg border">
      {[1, 2, 3].map((i) => (
        <div key={i}>
          {i > 1 && <ItemSeparator />}
          <Item>
            <ItemContent>
              <Skeleton className="h-5 w-32" />
              <Skeleton className="h-4 w-24" />
            </ItemContent>
            <ItemActions>
              <Skeleton className="h-5 w-5" />
            </ItemActions>
          </Item>
        </div>
      ))}
    </ItemGroup>
  )
}
