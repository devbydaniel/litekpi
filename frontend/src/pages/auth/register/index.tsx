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
import { OAuthButtons } from '@/shared/components/auth/oauth-buttons'
import { useRegisterForm } from './hooks/use-register-form'
import { RegisterForm } from './ui/register-form'

export function RegisterPage() {
  const { form, isLoading, error, success, email, onSubmit, handleOAuthLogin } =
    useRegisterForm()

  if (success) {
    return (
      <AuthLayout>
        <Card>
          <CardContent className="p-6">
            <StatusCard
              status="success"
              title="Check your email"
              description={`We've sent a verification link to ${email}. Please click the link to verify your account.`}
              action={
                <Link to="/login" className="font-medium hover:underline">
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
          <CardTitle className="text-2xl">Create an account</CardTitle>
          <CardDescription>Get started with LiteKPI</CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          {error && (
            <Alert variant="destructive">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <RegisterForm form={form} isLoading={isLoading} onSubmit={onSubmit} />

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <span className="w-full border-t" />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-card px-2 text-muted-foreground">
                Or continue with
              </span>
            </div>
          </div>

          <OAuthButtons isLoading={isLoading} onOAuthLogin={handleOAuthLogin} />

          <div className="text-center text-sm">
            Already have an account?{' '}
            <Link to="/login" className="font-medium hover:underline">
              Sign in
            </Link>
          </div>
        </CardContent>
      </Card>
    </AuthLayout>
  )
}
