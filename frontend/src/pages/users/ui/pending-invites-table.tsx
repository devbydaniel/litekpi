import { X, Clock, Mail } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { Badge } from '@/shared/components/ui/badge'
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
import type { InviteWithInviter } from '@/shared/api/generated/models'

interface PendingInvitesTableProps {
  invites: InviteWithInviter[]
  isLoading: boolean
  onCancel: (invite: InviteWithInviter) => void
}

const roleBadgeVariant = {
  admin: 'default',
  editor: 'secondary',
  viewer: 'outline',
} as const

function formatRelativeTime(date: string) {
  const now = new Date()
  const expires = new Date(date)
  const diffMs = expires.getTime() - now.getTime()
  const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays <= 0) return 'Expired'
  if (diffDays === 1) return 'Expires tomorrow'
  return `Expires in ${diffDays} days`
}

export function PendingInvitesTable({
  invites,
  isLoading,
  onCancel,
}: PendingInvitesTableProps) {
  if (isLoading) {
    return <InvitesTableSkeleton />
  }

  if (invites.length === 0) {
    return null
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
        <Mail className="h-4 w-4" />
        Pending Invites
      </div>
      <ItemGroup className="rounded-lg border border-dashed">
        {invites.map((invite, index) => (
          <div key={invite.id}>
            {index > 0 && <ItemSeparator />}
            <InviteListItem invite={invite} onCancel={onCancel} />
          </div>
        ))}
      </ItemGroup>
    </div>
  )
}

interface InviteListItemProps {
  invite: InviteWithInviter
  onCancel: (invite: InviteWithInviter) => void
}

function InviteListItem({ invite, onCancel }: InviteListItemProps) {
  return (
    <Item>
      <ItemContent>
        <div className="flex items-center gap-2">
          <ItemTitle className="text-muted-foreground">{invite.email}</ItemTitle>
          <Badge variant={roleBadgeVariant[invite.role as keyof typeof roleBadgeVariant] ?? 'outline'}>
            {invite.role}
          </Badge>
        </div>
        <ItemDescription className="flex items-center gap-2">
          <span>Invited by {invite.inviterName}</span>
          <span className="text-muted-foreground/50">-</span>
          <span className="flex items-center gap-1">
            <Clock className="h-3 w-3" />
            {invite.expiresAt ? formatRelativeTime(invite.expiresAt) : 'Unknown'}
          </span>
        </ItemDescription>
      </ItemContent>

      <ItemActions>
        <Button
          variant="ghost"
          size="icon"
          onClick={() => onCancel(invite)}
          title="Cancel invite"
        >
          <X className="h-4 w-4" />
          <span className="sr-only">Cancel invite</span>
        </Button>
      </ItemActions>
    </Item>
  )
}

function InvitesTableSkeleton() {
  return (
    <div className="space-y-3">
      <Skeleton className="h-4 w-32" />
      <ItemGroup className="rounded-lg border border-dashed">
        {[1, 2].map((i) => (
          <div key={i}>
            {i > 1 && <ItemSeparator />}
            <Item>
              <ItemContent>
                <Skeleton className="h-5 w-48" />
                <Skeleton className="h-4 w-32" />
              </ItemContent>
              <ItemActions>
                <Skeleton className="h-8 w-8" />
              </ItemActions>
            </Item>
          </div>
        ))}
      </ItemGroup>
    </div>
  )
}
