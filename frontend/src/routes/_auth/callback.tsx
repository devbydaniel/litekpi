import { createFileRoute, useSearch } from '@tanstack/react-router'
import { CallbackPage } from '@/pages/auth/callback'

export const Route = createFileRoute('/_auth/callback')({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>) => ({
    token: (search.token as string) || '',
    user: (search.user as string) || '',
    error: (search.error as string) || '',
  }),
})

function RouteComponent() {
  const search = useSearch({ from: '/_auth/callback' })
  return <CallbackPage token={search.token} user={search.user} error={search.error} />
}
