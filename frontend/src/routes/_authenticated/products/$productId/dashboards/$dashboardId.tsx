import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute(
  '/_authenticated/products/$productId/dashboards/$dashboardId'
)({
  component: DashboardDetailPage,
})

function DashboardDetailPage() {
  const { productId, dashboardId } = Route.useParams()

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Link to="/products" className="hover:text-foreground">
          Products
        </Link>
        <span>/</span>
        <Link
          to="/products/$productId"
          params={{ productId }}
          className="hover:text-foreground"
        >
          {productId}
        </Link>
        <span>/</span>
        <Link
          to="/products/$productId/dashboards"
          params={{ productId }}
          className="hover:text-foreground"
        >
          Dashboards
        </Link>
        <span>/</span>
        <span className="text-foreground">{dashboardId}</span>
      </div>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Dashboard: {dashboardId}</h1>
        <div className="flex gap-2">
          <button className="inline-flex h-10 items-center justify-center rounded-md border border-input bg-background px-4 py-2 text-sm font-medium ring-offset-background transition-colors hover:bg-accent hover:text-accent-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
            Edit
          </button>
          <button className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
            Add Widget
          </button>
        </div>
      </div>

      {/* Dashboard grid placeholder */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <div className="rounded-lg border p-6">
          <h3 className="font-medium">Widget Placeholder</h3>
          <p className="mt-2 text-sm text-muted-foreground">
            Charts and metrics will be rendered here.
          </p>
        </div>
      </div>
    </div>
  )
}
