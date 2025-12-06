import { Link } from '@tanstack/react-router'

export function SuccessMessage() {
  return (
    <div className="rounded-lg border bg-card p-6 shadow-sm">
      <div className="text-center">
        <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-green-100">
          <svg
            className="h-6 w-6 text-green-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
            />
          </svg>
        </div>
        <h1 className="text-xl font-semibold">Check your email</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          If an account with that email exists, we've sent a password reset link.
        </p>
        <Link
          to="/login"
          className="mt-4 inline-block text-sm text-muted-foreground hover:text-foreground"
        >
          Return to sign in
        </Link>
      </div>
    </div>
  )
}
