import { AuthLayout } from '@/layouts/auth'
import { Button } from '@/shared/components/ui/button'
import { Card, CardContent } from '@/shared/components/ui/card'
import { StatusCard } from '@/shared/components/ui/status-card'
import { useOAuthCallback } from './hooks/use-oauth-callback'

interface CallbackPageProps {
  token: string
  user: string
  error?: string
}

export function CallbackPage({ token, user: userEncoded, error: initialError }: CallbackPageProps) {
  const { error, navigate } = useOAuthCallback({ token, userEncoded, initialError })

  return (
    <AuthLayout>
      <Card>
        <CardContent className="p-6">
          {error ? (
            <StatusCard
              status="error"
              title="Authentication failed"
              description={error}
              action={
                <Button onClick={() => navigate({ to: '/login' })}>
                  Return to sign in
                </Button>
              }
            />
          ) : (
            <StatusCard
              status="loading"
              title="Completing sign in..."
              description="Please wait a moment."
            />
          )}
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
