import { createRootRoute, Outlet } from '@tanstack/react-router'
import { RootLayout } from '@/layouts/root'

export const Route = createRootRoute({
  component: () => (
    <RootLayout>
      <Outlet />
    </RootLayout>
  ),
})
