import { Link } from '@tanstack/react-router'

export function SuccessState() {
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
              d="M5 13l4 4L19 7"
            />
          </svg>
        </div>
        <h1 className="text-xl font-semibold">Email verified!</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Your email has been verified successfully. You can now sign in to your account.
        </p>
        <Link
          to="/login"
          className="mt-4 inline-flex h-10 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          Sign in
        </Link>
      </div>
    </div>
  )
}
