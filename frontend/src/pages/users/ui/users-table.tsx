import { MoreHorizontal, Trash, Shield, Users } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { EmptyState } from '@/shared/components/ui/empty-state'
import { Badge } from '@/shared/components/ui/badge'
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
import { useAuthStore } from '@/shared/stores/auth-store'
import type { User } from '@/shared/api/generated/models'

interface UsersTableProps {
  users: User[]
  isLoading: boolean
  onChangeRole: (user: User) => void
  onRemove: (user: User) => void
}

const roleBadgeVariant = {
  admin: 'default',
  editor: 'secondary',
  viewer: 'outline',
} as const

export function UsersTable({
  users,
  isLoading,
  onChangeRole,
  onRemove,
}: UsersTableProps) {
  const currentUser = useAuthStore((state) => state.user)

  if (isLoading) {
    return <UsersTableSkeleton />
  }

  if (users.length === 0) {
    return (
      <EmptyState
        icon={Users}
        title="No users yet"
        description="Invite users to join your organization."
      />
    )
  }

  return (
    <ItemGroup className="rounded-lg border">
      {users.map((user, index) => (
        <div key={user.id}>
          {index > 0 && <ItemSeparator />}
          <UserListItem
            user={user}
            isCurrentUser={user.id === currentUser?.id}
            onChangeRole={onChangeRole}
            onRemove={onRemove}
          />
        </div>
      ))}
    </ItemGroup>
  )
}

interface UserListItemProps {
  user: User
  isCurrentUser: boolean
  onChangeRole: (user: User) => void
  onRemove: (user: User) => void
}

function UserListItem({
  user,
  isCurrentUser,
  onChangeRole,
  onRemove,
}: UserListItemProps) {
  return (
    <Item>
      <ItemContent>
        <div className="flex items-center gap-2">
          <ItemTitle>{user.name}</ItemTitle>
          {isCurrentUser && (
            <Badge variant="outline" className="text-xs">
              You
            </Badge>
          )}
        </div>
        <ItemDescription className="flex items-center gap-2">
          <span>{user.email}</span>
          <span className="text-muted-foreground/50">-</span>
          <Badge variant={roleBadgeVariant[user.role as keyof typeof roleBadgeVariant] ?? 'outline'}>
            {user.role}
          </Badge>
        </ItemDescription>
      </ItemContent>

      {!isCurrentUser && (
        <ItemActions>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon">
                <MoreHorizontal className="h-4 w-4" />
                <span className="sr-only">Open menu</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => onChangeRole(user)}>
                <Shield className="mr-2 h-4 w-4" />
                Change Role
              </DropdownMenuItem>
              <DropdownMenuItem
                className="text-destructive focus:text-destructive"
                onClick={() => onRemove(user)}
              >
                <Trash className="mr-2 h-4 w-4" />
                Remove
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </ItemActions>
      )}
    </Item>
  )
}

function UsersTableSkeleton() {
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
