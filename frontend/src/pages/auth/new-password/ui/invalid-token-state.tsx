import { Link } from '@tanstack/react-router'

export function InvalidTokenState() {
  return (
    <div className="rounded-lg border bg-card p-6 shadow-sm">
      <div className="text-center">
        <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-destructive/10">
          <svg
            className="h-6 w-6 text-destructive"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </div>
        <h1 className="text-xl font-semibold">Invalid link</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          This password reset link is invalid or has expired.
        </p>
        <Link
          to="/reset-password"
          className="mt-4 inline-block text-sm font-medium hover:underline"
        >
          Request a new link
        </Link>
      </div>
    </div>
  )
}
