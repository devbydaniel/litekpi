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
import { useLoginForm } from './hooks/use-login-form'
import { LoginForm } from './ui/login-form'

interface LoginPageProps {
  error?: string
}

export function LoginPage({ error: initialError }: LoginPageProps) {
  const { form, isLoading, error, onSubmit } = useLoginForm({
    initialError,
  })

  return (
    <AuthLayout>
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-2xl">Welcome back</CardTitle>
          <CardDescription>Sign in to your account to continue</CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <LoginForm form={form} isLoading={isLoading} onSubmit={onSubmit} />

          <div className="text-center text-sm">
            <Link
              to="/reset-password"
              className="text-muted-foreground hover:text-foreground"
            >
              Forgot password?
            </Link>
          </div>

          <div className="text-center text-sm">
            Don't have an account?{' '}
            <Link to="/register" className="font-medium hover:underline">
              Sign up
            </Link>
          </div>
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
