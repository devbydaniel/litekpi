import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute(
  '/_authenticated/products/$productId/team'
)({
  component: TeamPage,
})

function TeamPage() {
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
        <span className="text-foreground">Team</span>
      </div>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Team</h1>
        <button className="inline-flex h-10 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground ring-offset-background transition-colors hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring">
          Invite Member
        </button>
      </div>

      {/* Team members table */}
      <div className="rounded-lg border">
        <div className="border-b px-4 py-3">
          <h3 className="font-medium">Members</h3>
        </div>
        <div className="divide-y">
          <div className="flex items-center justify-between px-4 py-3">
            <div>
              <p className="font-medium">owner@example.com</p>
              <p className="text-sm text-muted-foreground">Owner</p>
            </div>
            <span className="rounded-full bg-primary/10 px-2 py-1 text-xs font-medium text-primary">
              Owner
            </span>
          </div>
        </div>
      </div>
    </div>
  )
}
