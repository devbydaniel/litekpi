import { createFileRoute } from '@tanstack/react-router'
import { DataSourcesPage } from '@/pages/data-sources'

export const Route = createFileRoute('/_authenticated/data-sources')({
  component: DataSourcesPage,
})
