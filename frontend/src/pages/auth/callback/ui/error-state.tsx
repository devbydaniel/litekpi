interface ErrorStateProps {
  error: string
  onReturnToLogin: () => void
}

export function ErrorState({ error, onReturnToLogin }: ErrorStateProps) {
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
        <h1 className="text-xl font-semibold">Authentication failed</h1>
        <p className="mt-2 text-sm text-muted-foreground">{error}</p>
        <button
          onClick={onReturnToLogin}
          className="mt-4 inline-flex h-10 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          Return to sign in
        </button>
      </div>
    </div>
  )
}
