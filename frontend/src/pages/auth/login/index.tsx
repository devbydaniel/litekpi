import { Link } from '@tanstack/react-router'
import { AuthLayout } from '@/layouts/auth'
import { useLoginForm } from './hooks/use-login-form'
import { LoginForm } from './ui/login-form'
import { OAuthButtons } from './ui/oauth-buttons'
import { ErrorAlert } from './ui/error-alert'

interface LoginPageProps {
  error?: string
}

export function LoginPage({ error: initialError }: LoginPageProps) {
  const {
    email,
    setEmail,
    password,
    setPassword,
    isLoading,
    error,
    handleSubmit,
    handleOAuthLogin,
  } = useLoginForm({ initialError })

  return (
    <AuthLayout>
      <div className="rounded-lg border bg-card p-6 shadow-sm">
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-semibold">Welcome back</h1>
          <p className="text-sm text-muted-foreground">
            Sign in to your account to continue
          </p>
        </div>

        {error && <ErrorAlert message={error} />}

        <LoginForm
          email={email}
          password={password}
          isLoading={isLoading}
          onEmailChange={setEmail}
          onPasswordChange={setPassword}
          onSubmit={handleSubmit}
        />

        <div className="relative my-6">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t" />
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-card px-2 text-muted-foreground">Or continue with</span>
          </div>
        </div>

        <OAuthButtons isLoading={isLoading} onOAuthLogin={handleOAuthLogin} />

        <div className="mt-4 text-center text-sm">
          <Link
            to="/reset-password"
            className="text-muted-foreground hover:text-foreground"
          >
            Forgot password?
          </Link>
        </div>

        <div className="mt-6 text-center text-sm">
          Don't have an account?{' '}
          <Link to="/register" className="font-medium hover:underline">
            Sign up
          </Link>
        </div>
      </div>
    </AuthLayout>
  )
}
