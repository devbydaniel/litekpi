import { Link } from '@tanstack/react-router'

import { AuthLayout } from '@/layouts/auth'
import { Button } from '@/shared/components/ui/button'
import { Card, CardContent } from '@/shared/components/ui/card'
import { StatusCard } from '@/shared/components/ui/status-card'
import { useEmailVerification } from './hooks/use-email-verification'

interface VerifyEmailPageProps {
  token: string
}

export function VerifyEmailPage({ token }: VerifyEmailPageProps) {
  const { status, error } = useEmailVerification(token)

  return (
    <AuthLayout>
      <Card>
        <CardContent className="p-6">
          {status === 'loading' && (
            <StatusCard
              status="loading"
              title="Verifying your email..."
              description="Please wait a moment."
            />
          )}
          {status === 'success' && (
            <StatusCard
              status="success"
              title="Email verified!"
              description="Your email has been verified successfully. You can now sign in to your account."
              action={
                <Button asChild>
                  <Link to="/login">Sign in</Link>
                </Button>
              }
            />
          )}
          {status === 'error' && (
            <StatusCard
              status="error"
              title="Verification failed"
              description={error || 'Unknown error'}
              action={
                <Link
                  to="/login"
                  className="text-sm text-muted-foreground hover:text-foreground"
                >
                  Return to sign in
                </Link>
              }
            />
          )}
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
