import { createFileRoute, useSearch } from '@tanstack/react-router'
import { NewPasswordPage } from '@/pages/auth/new-password'

export const Route = createFileRoute('/_auth/new-password')({
  component: RouteComponent,
  validateSearch: (search: Record<string, unknown>) => ({
    token: (search.token as string) || '',
  }),
})

function RouteComponent() {
  const { token } = useSearch({ from: '/_auth/new-password' })
  return <NewPasswordPage token={token} />
}
