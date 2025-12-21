import { Link, useSearch } from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import { AuthLayout } from '@/layouts/auth'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/components/ui/card'
import { StatusCard } from '@/shared/components/ui/status-card'
import { Button } from '@/shared/components/ui/button'
import { useAcceptInvite } from './hooks/use-accept-invite'
import { AcceptInviteForm } from './ui/accept-invite-form'

export function AcceptInvitePage() {
  const { token } = useSearch({ from: '/_auth/accept-invite' })
  const {
    isValidating,
    isValid,
    inviteInfo,
    success,
    isAccepting,
    handleAccept,
    goToLogin,
  } = useAcceptInvite(token)

  if (!token) {
    return (
      <AuthLayout>
        <Card>
          <CardContent className="p-6">
            <StatusCard
              status="error"
              title="Invalid invite link"
              description="This invite link is missing the required token. Please check the link and try again."
              action={
                <Link to="/login" className="font-medium hover:underline">
                  Go to sign in
                </Link>
              }
            />
          </CardContent>
        </Card>
      </AuthLayout>
    )
  }

  if (isValidating) {
    return (
      <AuthLayout>
        <Card>
          <CardContent className="flex items-center justify-center p-12">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </CardContent>
        </Card>
      </AuthLayout>
    )
  }

  if (!isValid) {
    return (
      <AuthLayout>
        <Card>
          <CardContent className="p-6">
            <StatusCard
              status="error"
              title="Invalid or expired invite"
              description="This invitation link is invalid or has expired. Please ask for a new invite."
              action={
                <Link to="/login" className="font-medium hover:underline">
                  Go to sign in
                </Link>
              }
            />
          </CardContent>
        </Card>
      </AuthLayout>
    )
  }

  if (success) {
    return (
      <AuthLayout>
        <Card>
          <CardContent className="p-6">
            <StatusCard
              status="success"
              title="Account created!"
              description="Your account has been created successfully. You can now sign in."
              action={
                <Button onClick={goToLogin}>
                  Go to sign in
                </Button>
              }
            />
          </CardContent>
        </Card>
      </AuthLayout>
    )
  }

  return (
    <AuthLayout>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Join {inviteInfo?.organizationName}</CardTitle>
          <CardDescription>
            Create your account to accept this invitation
          </CardDescription>
        </CardHeader>

        <CardContent>
          {inviteInfo && (
            <AcceptInviteForm
              inviteInfo={inviteInfo}
              isLoading={isAccepting}
              onSubmit={handleAccept}
            />
          )}
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
