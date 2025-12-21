import { createFileRoute } from '@tanstack/react-router'
import { z } from 'zod'
import { AcceptInvitePage } from '@/pages/auth/accept-invite'

const searchSchema = z.object({
  token: z.string().optional(),
})

export const Route = createFileRoute('/_auth/accept-invite')({
  component: AcceptInvitePage,
  validateSearch: searchSchema,
})
