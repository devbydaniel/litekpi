import { Link } from '@tanstack/react-router'

interface SuccessMessageProps {
  email: string
}

export function SuccessMessage({ email }: SuccessMessageProps) {
  return (
    <div className="rounded-lg border bg-card p-6 shadow-sm">
      <div className="mb-6 text-center">
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
        <h1 className="text-2xl font-semibold">Check your email</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          We've sent a verification link to <strong>{email}</strong>. Please
          click the link to verify your account.
        </p>
      </div>
      <div className="text-center text-sm">
        <Link to="/login" className="font-medium hover:underline">
          Return to sign in
        </Link>
      </div>
    </div>
  )
}
