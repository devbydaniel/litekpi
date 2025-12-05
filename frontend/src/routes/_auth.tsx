import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'
import { useAuthStore } from '@/shared/stores/auth-store'

export const Route = createFileRoute('/_auth')({
  beforeLoad: () => {
    const { isAuthenticated } = useAuthStore.getState()
    if (isAuthenticated) {
      throw redirect({ to: '/products' })
    }
  },
  component: AuthLayout,
})

function AuthLayout() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/50">
      <div className="w-full max-w-md p-4">
        <Outlet />
      </div>
    </div>
  )
}
