export function LoadingState() {
  return (
    <div className="rounded-lg border bg-card p-6 shadow-sm">
      <div className="text-center">
        <div className="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
        <h1 className="text-xl font-semibold">Verifying your email...</h1>
        <p className="mt-2 text-sm text-muted-foreground">Please wait a moment.</p>
      </div>
    </div>
  )
}
