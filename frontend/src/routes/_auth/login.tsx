import { createFileRoute, useSearch } from '@tanstack/react-router'
import { LoginPage } from '@/pages/auth/login'

export const Route = createFileRoute('/_auth/login')({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>): { error?: string } => ({
    error: search.error ? String(search.error) : undefined,
  }),
})

function RouteComponent() {
  const search = useSearch({ from: '/_auth/login' })
  return <LoginPage error={search.error} />
}
