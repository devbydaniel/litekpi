import { MoreHorizontal, Pencil, Trash, Key } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { Badge } from '@/shared/components/ui/badge'
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
import type { MCPAPIKey } from '@/shared/api/generated/models'

interface MCPApiKeyListProps {
  keys: MCPAPIKey[]
  isLoading: boolean
  onEdit: (key: MCPAPIKey) => void
  onDelete: (key: MCPAPIKey) => void
}

export function MCPApiKeyList({
  keys,
  isLoading,
  onEdit,
  onDelete,
}: MCPApiKeyListProps) {
  if (isLoading) {
    return <MCPApiKeyListSkeleton />
  }

  if (keys.length === 0) {
    return (
      <EmptyState
        icon={Key}
        title="No MCP API keys yet"
        description="Create an MCP API key to enable external access to your data sources."
      />
    )
  }

  return (
    <ItemGroup className="rounded-lg border">
      {keys.map((key, index) => (
        <div key={key.id}>
          {index > 0 && <ItemSeparator />}
          <MCPApiKeyListItem apiKey={key} onEdit={onEdit} onDelete={onDelete} />
        </div>
      ))}
    </ItemGroup>
  )
}

interface MCPApiKeyListItemProps {
  apiKey: MCPAPIKey
  onEdit: (key: MCPAPIKey) => void
  onDelete: (key: MCPAPIKey) => void
}

function MCPApiKeyListItem({ apiKey, onEdit, onDelete }: MCPApiKeyListItemProps) {
  const dataSourceCount = apiKey.allowedDataSourceIds?.length ?? 0

  return (
    <Item>
      <ItemContent>
        <div className="flex items-center gap-2">
          <ItemTitle>{apiKey.name}</ItemTitle>
          <Badge variant="outline" className="text-xs">
            {dataSourceCount} data source{dataSourceCount !== 1 ? 's' : ''}
          </Badge>
        </div>
        <ItemDescription>
          Created{' '}
          {apiKey.createdAt
            ? new Date(apiKey.createdAt).toLocaleDateString()
            : '-'}
          {apiKey.lastUsedAt && (
            <>
              {' · '}
              Last used {new Date(apiKey.lastUsedAt).toLocaleDateString()}
            </>
          )}
        </ItemDescription>
      </ItemContent>

      <ItemActions>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon">
              <MoreHorizontal className="h-4 w-4" />
              <span className="sr-only">Open menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => onEdit(apiKey)}>
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </DropdownMenuItem>
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={() => onDelete(apiKey)}
            >
              <Trash className="mr-2 h-4 w-4" />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </ItemActions>
    </Item>
  )
}

function MCPApiKeyListSkeleton() {
  return (
    <ItemGroup className="rounded-lg border">
      {[1, 2, 3].map((i) => (
        <div key={i}>
          {i > 1 && <ItemSeparator />}
          <Item>
            <ItemContent>
              <Skeleton className="h-5 w-32" />
              <Skeleton className="h-4 w-48" />
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
