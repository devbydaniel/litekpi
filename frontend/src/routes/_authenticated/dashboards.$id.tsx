import { createFileRoute } from '@tanstack/react-router'
import { DashboardPage } from '@/pages/dashboard'

export const Route = createFileRoute('/_authenticated/dashboards/$id')({
  component: DashboardRoute,
})

function DashboardRoute() {
  const { id } = Route.useParams()
  return <DashboardPage dashboardId={id} />
}
