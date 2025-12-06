import { Link } from '@tanstack/react-router'
import { AuthLayout } from '@/layouts/auth'
import { useRegisterForm } from './hooks/use-register-form'
import { RegisterForm } from './ui/register-form'
import { OAuthButtons } from './ui/oauth-buttons'
import { ErrorAlert } from './ui/error-alert'
import { SuccessMessage } from './ui/success-message'

export function RegisterPage() {
  const {
    email,
    setEmail,
    password,
    setPassword,
    confirmPassword,
    setConfirmPassword,
    isLoading,
    error,
    success,
    handleSubmit,
    handleOAuthLogin,
  } = useRegisterForm()

  if (success) {
    return (
      <AuthLayout>
        <SuccessMessage email={email} />
      </AuthLayout>
    )
  }

  return (
    <AuthLayout>
      <div className="rounded-lg border bg-card p-6 shadow-sm">
        <div className="mb-6 text-center">
          <h1 className="text-2xl font-semibold">Create an account</h1>
          <p className="text-sm text-muted-foreground">Get started with Trackable</p>
        </div>

        {error && <ErrorAlert message={error} />}

        <RegisterForm
          email={email}
          password={password}
          confirmPassword={confirmPassword}
          isLoading={isLoading}
          onEmailChange={setEmail}
          onPasswordChange={setPassword}
          onConfirmPasswordChange={setConfirmPassword}
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

        <div className="mt-6 text-center text-sm">
          Already have an account?{' '}
          <Link to="/login" className="font-medium hover:underline">
            Sign in
          </Link>
        </div>
      </div>
    </AuthLayout>
  )
}
