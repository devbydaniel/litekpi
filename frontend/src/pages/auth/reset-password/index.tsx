import { Link } from '@tanstack/react-router'
import { AuthLayout } from '@/layouts/auth'
import { useResetPasswordForm } from './hooks/use-reset-password-form'
import { ResetPasswordForm } from './ui/reset-password-form'
import { ErrorAlert } from './ui/error-alert'
import { SuccessMessage } from './ui/success-message'

export function ResetPasswordPage() {
  const {
    email,
    setEmail,
    isLoading,
    error,
    success,
    handleSubmit,
  } = useResetPasswordForm()

  if (success) {
    return (
      <AuthLayout>
        <SuccessMessage />
      </AuthLayout>
    )
  }

  return (
    <AuthLayout>
      <div className="rounded-lg border bg-card p-6 shadow-sm">
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-semibold">Reset password</h1>
          <p className="text-sm text-muted-foreground">
            Enter your email to receive a password reset link
          </p>
        </div>

        {error && <ErrorAlert message={error} />}

        <ResetPasswordForm
          email={email}
          isLoading={isLoading}
          onEmailChange={setEmail}
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
