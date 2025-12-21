import { createFileRoute } from '@tanstack/react-router'
import { ReportPage } from '@/pages/report'

export const Route = createFileRoute('/_authenticated/reports/$id')({
  component: ReportRoute,
})

function ReportRoute() {
  const { id } = Route.useParams()
  return <ReportPage reportId={id} />
}
