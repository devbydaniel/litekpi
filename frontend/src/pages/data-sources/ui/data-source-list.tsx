import { MoreHorizontal, Trash, RefreshCw, Database } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { EmptyState } from '@/shared/components/ui/empty-state'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
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
import type { DataSource } from '@/shared/api/generated/models'
import { useAuth } from '@/shared/hooks/use-auth'

interface DataSourceListProps {
  dataSources: DataSource[]
  isLoading: boolean
  onDelete: (dataSource: DataSource) => void
  onRegenerateKey: (dataSource: DataSource) => void
}

export function DataSourceList({
  dataSources,
  isLoading,
  onDelete,
  onRegenerateKey,
}: DataSourceListProps) {
  const { user } = useAuth()
  const isAdmin = user?.role === 'admin'

  if (isLoading) {
    return <DataSourceListSkeleton />
  }

  if (dataSources.length === 0) {
    return (
      <EmptyState
        icon={Database}
        title="No data sources yet"
        description="Create your first data source to start tracking metrics."
      />
    )
  }

  return (
    <ItemGroup className="rounded-lg border">
      {dataSources.map((dataSource, index) => (
        <div key={dataSource.id}>
          {index > 0 && <ItemSeparator />}
          <DataSourceListItem
            dataSource={dataSource}
            onDelete={onDelete}
            onRegenerateKey={onRegenerateKey}
            isAdmin={isAdmin}
          />
        </div>
      ))}
    </ItemGroup>
  )
}

interface DataSourceListItemProps {
  dataSource: DataSource
  onDelete: (dataSource: DataSource) => void
  onRegenerateKey: (dataSource: DataSource) => void
  isAdmin: boolean
}

function DataSourceListItem({
  dataSource,
  onDelete,
  onRegenerateKey,
  isAdmin,
}: DataSourceListItemProps) {
  return (
    <Item>
      <ItemContent>
        <ItemTitle>{dataSource.name}</ItemTitle>
        <ItemDescription>
          Created {dataSource.createdAt ? new Date(dataSource.createdAt).toLocaleDateString() : '-'}
        </ItemDescription>
      </ItemContent>

      {isAdmin && (
        <ItemActions>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <MoreHorizontal className="h-4 w-4" />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => onRegenerateKey(dataSource)}>
                <RefreshCw className="mr-2 h-4 w-4" />
                Regenerate API Key
              </DropdownMenuItem>
              <DropdownMenuItem
                className="text-destructive focus:text-destructive"
                onClick={() => onDelete(dataSource)}
              >
                <Trash className="mr-2 h-4 w-4" />
                Delete
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </ItemActions>
      )}
    </Item>
  )
}

function DataSourceListSkeleton() {
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
              <Skeleton className="h-8 w-8" />
            </ItemActions>
          </Item>
        </div>
      ))}
    </ItemGroup>
  )
}
