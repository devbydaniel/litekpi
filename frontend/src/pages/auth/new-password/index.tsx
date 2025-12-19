import { Link } from '@tanstack/react-router'

import { AuthLayout } from '@/layouts/auth'
import { Alert, AlertDescription } from '@/shared/components/ui/alert'
import { Button } from '@/shared/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/components/ui/card'
import { StatusCard } from '@/shared/components/ui/status-card'
import { useNewPasswordForm } from './hooks/use-new-password-form'
import { NewPasswordForm } from './ui/new-password-form'

interface NewPasswordPageProps {
  token: string
}

export function NewPasswordPage({ token }: NewPasswordPageProps) {
  const { form, isLoading, error, success, onSubmit } =
    useNewPasswordForm(token)

  if (!token) {
    return (
      <AuthLayout>
        <Card>
          <CardContent className="p-6">
            <StatusCard
              status="error"
              title="Invalid reset link"
              description="This password reset link is invalid or has expired."
              action={
                <Button asChild variant="outline">
                  <Link to="/reset-password">Request a new link</Link>
                </Button>
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
              title="Password reset!"
              description="Your password has been reset successfully. You can now sign in with your new password."
              action={
                <Button asChild>
                  <Link to="/login">Sign in</Link>
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
          <CardTitle className="text-2xl">Set new password</CardTitle>
          <CardDescription>Enter your new password below</CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <NewPasswordForm
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
