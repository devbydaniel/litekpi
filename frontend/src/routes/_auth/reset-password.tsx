import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/reset-password')({
  component: ResetPasswordPage,
})

function ResetPasswordPage() {
  return (
    <div className="rounded-lg border bg-card p-6 shadow-sm">
      <div className="mb-6 text-center">
        <h1 className="text-2xl font-semibold">Reset password</h1>
        <p className="text-sm text-muted-foreground">
          Enter your email to receive a password reset link
        </p>
      </div>

      <form className="space-y-4">
        <div className="space-y-2">
          <label htmlFor="email" className="text-sm font-medium">
            Email
          </label>
          <input
            id="email"
            type="email"
            placeholder="you@example.com"
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
        </div>

        <button
          type="submit"
          className="inline-flex h-10 w-full items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          Send reset link
        </button>
      </form>

      <div className="mt-6 text-center text-sm">
        Remember your password?{' '}
        <Link to="/login" className="font-medium hover:underline">
          Sign in
        </Link>
      </div>
    </div>
  )
}
