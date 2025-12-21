import { createFileRoute } from '@tanstack/react-router'
import { ReportsPage } from '@/pages/reports'

export const Route = createFileRoute('/_authenticated/reports')({
  component: ReportsPage,
})
