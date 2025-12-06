import { AuthLayout } from '@/layouts/auth'
import { useEmailVerification } from './hooks/use-email-verification'
import { LoadingState } from './ui/loading-state'
import { SuccessState } from './ui/success-state'
import { ErrorState } from './ui/error-state'

interface VerifyEmailPageProps {
  token: string
}

export function VerifyEmailPage({ token }: VerifyEmailPageProps) {
  const { status, error } = useEmailVerification(token)

  return (
    <AuthLayout>
      {status === 'loading' && <LoadingState />}
      {status === 'success' && <SuccessState />}
      {status === 'error' && <ErrorState error={error || 'Unknown error'} />}
    </AuthLayout>
  )
}
