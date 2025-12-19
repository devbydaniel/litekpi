import { Link } from '@tanstack/react-router'

import { AuthLayout } from '@/layouts/auth'
import { Alert, AlertDescription } from '@/shared/components/ui/alert'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/components/ui/card'
import { StatusCard } from '@/shared/components/ui/status-card'
import { useResetPasswordForm } from './hooks/use-reset-password-form'
import { ResetPasswordForm } from './ui/reset-password-form'

export function ResetPasswordPage() {
  const { form, isLoading, error, success, onSubmit } = useResetPasswordForm()

  if (success) {
    return (
      <AuthLayout>
        <Card>
          <CardContent className="p-6">
            <StatusCard
              status="success"
              title="Check your email"
              description="If an account with that email exists, we've sent a password reset link."
              action={
                <Link
                  to="/login"
                  className="text-sm text-muted-foreground hover:text-foreground"
                >
                  Return to sign in
                </Link>
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
          <CardTitle className="text-2xl">Reset password</CardTitle>
          <CardDescription>
            Enter your email to receive a password reset link
          </CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <ResetPasswordForm
            form={form}
            isLoading={isLoading}
            onSubmit={onSubmit}
          />

          <div className="text-center text-sm">
            Remember your password?{' '}
            <Link to="/login" className="font-medium hover:underline">
              Sign in
            </Link>
          </div>
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
