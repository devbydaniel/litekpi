import { createFileRoute, Outlet, redirect, Link } from '@tanstack/react-router'
import { useAuthStore } from '@/shared/stores/auth-store'

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: () => {
    const { isAuthenticated } = useAuthStore.getState()
    if (!isAuthenticated) {
      throw redirect({ to: '/login' })
    }
  },
  component: AuthenticatedLayout,
})

function AuthenticatedLayout() {
  const { user, logout } = useAuthStore()

  return (
    <div className="flex min-h-screen">
      {/* Sidebar */}
      <aside className="hidden w-64 border-r bg-muted/30 lg:block">
        <div className="flex h-full flex-col">
          {/* Logo */}
          <div className="flex h-14 items-center border-b px-4">
            <Link to="/products" className="flex items-center gap-2 font-semibold">
              <span className="text-xl">📊</span>
              <span>Trackable</span>
            </Link>
          </div>

          {/* Navigation */}
          <nav className="flex-1 space-y-1 p-4">
            <Link
              to="/products"
              className="flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium hover:bg-muted"
              activeProps={{ className: 'bg-muted' }}
            >
              Products
            </Link>
          </nav>

          {/* User menu */}
          <div className="border-t p-4">
            <div className="flex items-center gap-2">
              <div className="flex-1 truncate">
                <p className="text-sm font-medium">{user?.email}</p>
              </div>
              <button
                onClick={logout}
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1">
        {/* Mobile header */}
        <header className="flex h-14 items-center gap-4 border-b px-4 lg:hidden">
          <Link to="/products" className="flex items-center gap-2 font-semibold">
            <span className="text-xl">📊</span>
            <span>Trackable</span>
          </Link>
        </header>

        <div className="p-4 lg:p-6">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
