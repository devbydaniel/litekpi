import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute(
  '/_authenticated/products/$productId/settings'
)({
  component: ProductSettingsPage,
})

function ProductSettingsPage() {
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
        <span className="text-foreground">Settings</span>
      </div>

      <h1 className="text-2xl font-semibold">Settings</h1>

      <div className="space-y-6">
        {/* General settings */}
        <div className="rounded-lg border p-6">
          <h2 className="text-lg font-medium">General</h2>
          <p className="text-sm text-muted-foreground">
            Manage your product settings
          </p>
          <div className="mt-4 space-y-4">
            <div className="space-y-2">
              <label htmlFor="name" className="text-sm font-medium">
                Product Name
              </label>
              <input
                id="name"
                type="text"
                placeholder="My Product"
                className="flex h-10 w-full max-w-md rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
              />
            </div>
          </div>
        </div>

        {/* API Key */}
        <div className="rounded-lg border p-6">
          <h2 className="text-lg font-medium">API Key</h2>
          <p className="text-sm text-muted-foreground">
            Use this key to send data to your product
          </p>
          <div className="mt-4">
            <code className="rounded bg-muted px-2 py-1 text-sm">
              ••••••••••••••••
            </code>
            <button className="ml-2 text-sm text-primary hover:underline">
              Regenerate
            </button>
          </div>
        </div>

        {/* Danger zone */}
        <div className="rounded-lg border border-destructive/50 p-6">
          <h2 className="text-lg font-medium text-destructive">Danger Zone</h2>
          <p className="text-sm text-muted-foreground">
            Irreversible actions for your product
          </p>
          <div className="mt-4">
            <button className="inline-flex h-10 items-center justify-center rounded-md bg-destructive px-4 py-2 text-sm font-medium text-destructive-foreground ring-offset-background transition-colors hover:bg-destructive/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
              Delete Product
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
