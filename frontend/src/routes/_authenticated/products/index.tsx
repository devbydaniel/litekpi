import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/products/')({
  component: ProductsPage,
})

function ProductsPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold">Products</h1>
          <p className="text-sm text-muted-foreground">
            Manage your products and track their metrics
          </p>
        </div>
        <button className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
          New Product
        </button>
      </div>

      {/* Empty state */}
      <div className="flex min-h-[400px] flex-col items-center justify-center rounded-lg border border-dashed p-8 text-center">
        <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-muted">
          <span className="text-2xl">📦</span>
        </div>
        <h3 className="mt-4 text-lg font-semibold">No products</h3>
        <p className="mt-2 text-sm text-muted-foreground">
          Get started by creating your first product to track metrics.
        </p>
        <button className="mt-4 inline-flex h-10 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
          Create Product
        </button>
      </div>
    </div>
  )
}
