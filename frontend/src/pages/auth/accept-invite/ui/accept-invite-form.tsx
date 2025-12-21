import { useState } from 'react'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import { Badge } from '@/shared/components/ui/badge'
import type { ValidateInviteResponse } from '@/shared/api/generated/models'

interface AcceptInviteFormProps {
  inviteInfo: ValidateInviteResponse
  isLoading: boolean
  onSubmit: (name: string, password: string) => Promise<void>
}

const roleBadgeVariant = {
  admin: 'default',
  editor: 'secondary',
  viewer: 'outline',
} as const

export function AcceptInviteForm({
  inviteInfo,
  isLoading,
  onSubmit,
}: AcceptInviteFormProps) {
  const [name, setName] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (password !== confirmPassword) {
      setError('Passwords do not match')
      return
    }

    if (password.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }

    await onSubmit(name, password)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2 rounded-lg border bg-muted/50 p-4">
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Organization</span>
          <span className="font-medium">{inviteInfo.organizationName}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Your email</span>
          <span className="font-medium">{inviteInfo.email}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Role</span>
          <Badge variant={roleBadgeVariant[inviteInfo.role as keyof typeof roleBadgeVariant] ?? 'outline'}>
            {inviteInfo.role}
          </Badge>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-sm text-muted-foreground">Invited by</span>
          <span className="font-medium">{inviteInfo.inviterName}</span>
        </div>
      </div>

      {error && (
        <p className="text-sm text-destructive">{error}</p>
      )}

      <div className="space-y-2">
        <Label htmlFor="name">Your name</Label>
        <Input
          id="name"
          placeholder="Enter your name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="password">Password</Label>
        <Input
          id="password"
          type="password"
          placeholder="Create a password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={8}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="confirmPassword">Confirm password</Label>
        <Input
          id="confirmPassword"
          type="password"
          placeholder="Confirm your password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
          minLength={8}
        />
      </div>

      <Button type="submit" className="w-full" disabled={isLoading || !name || !password || !confirmPassword}>
        {isLoading ? 'Creating account...' : 'Accept Invite'}
      </Button>
    </form>
  )
}
