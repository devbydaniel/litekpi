import { AuthLayout } from '@/layouts/auth'
import { useOAuthCallback } from './hooks/use-oauth-callback'
import { LoadingState } from './ui/loading-state'
import { ErrorState } from './ui/error-state'

interface CallbackPageProps {
  token: string
  user: string
  error?: string
}

export function CallbackPage({ token, user: userEncoded, error: initialError }: CallbackPageProps) {
  const { error, navigate } = useOAuthCallback({ token, userEncoded, initialError })

  return (
    <AuthLayout>
      {error ? (
        <ErrorState error={error} onReturnToLogin={() => navigate({ to: '/login' })} />
      ) : (
        <LoadingState />
      )}
    </AuthLayout>
  )
}
