import { Link } from '@tanstack/react-router'
import { AuthLayout } from '@/layouts/auth'
import { useNewPasswordForm } from './hooks/use-new-password-form'
import { NewPasswordForm } from './ui/new-password-form'
import { ErrorAlert } from './ui/error-alert'
import { InvalidTokenState } from './ui/invalid-token-state'
import { SuccessState } from './ui/success-state'

interface NewPasswordPageProps {
  token: string
}

export function NewPasswordPage({ token }: NewPasswordPageProps) {
  const {
    password,
    setPassword,
    confirmPassword,
    setConfirmPassword,
    isLoading,
    error,
    success,
    handleSubmit,
  } = useNewPasswordForm(token)

  if (!token) {
    return (
      <AuthLayout>
        <InvalidTokenState />
      </AuthLayout>
    )
  }

  if (success) {
    return (
      <AuthLayout>
        <SuccessState />
      </AuthLayout>
    )
  }

  return (
    <AuthLayout>
      <div className="rounded-lg border bg-card p-6 shadow-sm">
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-semibold">Set new password</h1>
          <p className="text-sm text-muted-foreground">
            Enter your new password below
          </p>
        </div>

        {error && <ErrorAlert message={error} />}

        <NewPasswordForm
          password={password}
          confirmPassword={confirmPassword}
          isLoading={isLoading}
          onPasswordChange={setPassword}
          onConfirmPasswordChange={setConfirmPassword}
          onSubmit={handleSubmit}
        />

        <div className="mt-6 text-center text-sm">
          Remember your password?{' '}
          <Link to="/login" className="font-medium hover:underline">
            Sign in
          </Link>
        </div>
      </div>
    </AuthLayout>
  )
}
