import { createFileRoute } from '@tanstack/react-router'
import { CompleteSetupPage } from '@/pages/auth/complete-setup'

interface CompleteSetupSearch {
  token: string
  email: string
  name?: string
}

export const Route = createFileRoute('/_auth/complete-setup')({
  component: CompleteSetupPage,
  validateSearch: (search: Record<string, unknown>): CompleteSetupSearch => ({
    token: (search.token as string) || '',
    email: (search.email as string) || '',
    name: (search.name as string) || undefined,
  }),
})
