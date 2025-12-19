import { Link, useSearch } from '@tanstack/react-router'

import { AuthLayout } from '@/layouts/auth'
import { Alert, AlertDescription } from '@/shared/components/ui/alert'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/shared/components/ui/card'
import { useCompleteSetupForm } from './hooks/use-complete-setup-form'
import { CompleteSetupForm } from './ui/complete-setup-form'

export function CompleteSetupPage() {
  const { token, email, name } = useSearch({ from: '/_auth/complete-setup' })
  const { form, isLoading, error, onSubmit } = useCompleteSetupForm({
    token,
    email,
    initialName: name,
  })

  if (!token) {
    return (
      <AuthLayout>
        <Card>
          <CardContent className="p-6 text-center">
            <p className="text-muted-foreground mb-4">
              Invalid or expired setup link. Please try signing in again.
            </p>
            <Link to="/login" className="font-medium hover:underline">
              Return to sign in
            </Link>
          </CardContent>
        </Card>
      </AuthLayout>
    )
  }

  return (
    <AuthLayout>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Complete your account</CardTitle>
          <CardDescription>
            Just a few more details to get started
            {email && (
              <span className="block mt-1 text-foreground font-medium">{email}</span>
            )}
          </CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <CompleteSetupForm form={form} isLoading={isLoading} onSubmit={onSubmit} />

          <div className="text-center text-sm text-muted-foreground">
            Want to use a different account?{' '}
            <Link to="/login" className="font-medium hover:underline">
              Sign in
            </Link>
          </div>
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
