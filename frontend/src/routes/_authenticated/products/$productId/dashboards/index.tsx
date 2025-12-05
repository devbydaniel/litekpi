import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute(
  '/_authenticated/products/$productId/dashboards/'
)({
  component: DashboardsPage,
})

function DashboardsPage() {
  const { productId } = Route.useParams()

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
        <span className="text-foreground">Dashboards</span>
      </div>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Dashboards</h1>
        <button className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
          New Dashboard
        </button>
      </div>

      {/* Empty state */}
      <div className="flex min-h-[300px] flex-col items-center justify-center rounded-lg border border-dashed p-8 text-center">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-muted">
          <span className="text-2xl">📊</span>
        </div>
        <h3 className="mt-4 text-lg font-semibold">No dashboards</h3>
        <p className="mt-2 text-sm text-muted-foreground">
          Create a dashboard to visualize your metrics.
        </p>
      </div>
    </div>
  )
}
