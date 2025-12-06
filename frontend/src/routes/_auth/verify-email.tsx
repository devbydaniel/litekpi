import { createFileRoute, useSearch } from '@tanstack/react-router'
import { VerifyEmailPage } from '@/pages/auth/verify-email'

export const Route = createFileRoute('/_auth/verify-email')({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>) => ({
    token: (search.token as string) || '',
  }),
})

function RouteComponent() {
  const { token } = useSearch({ from: '/_auth/verify-email' })
  return <VerifyEmailPage token={token} />
}
