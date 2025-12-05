import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/products/$productId/')({
  component: ProductDetailPage,
})

function ProductDetailPage() {
  const { productId } = Route.useParams()

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Link to="/products" className="hover:text-foreground">
          Products
        </Link>
        <span>/</span>
        <span className="text-foreground">Product Details</span>
      </div>

      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Product: {productId}</h1>
          <p className="text-sm text-muted-foreground">
            View and manage your product metrics
          </p>
        </div>
      </div>

      {/* Navigation tabs */}
      <div className="flex gap-4 border-b">
        <Link
          to="/products/$productId"
          params={{ productId }}
          className="border-b-2 border-primary pb-2 text-sm font-medium"
        >
          Overview
        </Link>
        <Link
          to="/products/$productId/dashboards"
          params={{ productId }}
          className="pb-2 text-sm text-muted-foreground hover:text-foreground"
        >
          Dashboards
        </Link>
        <Link
          to="/products/$productId/reports"
          params={{ productId }}
          className="pb-2 text-sm text-muted-foreground hover:text-foreground"
        >
          Reports
        </Link>
        <Link
          to="/products/$productId/team"
          params={{ productId }}
          className="pb-2 text-sm text-muted-foreground hover:text-foreground"
        >
          Team
        </Link>
        <Link
          to="/products/$productId/settings"
          params={{ productId }}
          className="pb-2 text-sm text-muted-foreground hover:text-foreground"
        >
          Settings
        </Link>
      </div>

      {/* Content placeholder */}
      <div className="rounded-lg border p-6">
        <p className="text-muted-foreground">
          Product overview and metrics will be displayed here.
        </p>
      </div>
    </div>
  )
}
