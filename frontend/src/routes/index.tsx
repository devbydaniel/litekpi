import { createFileRoute, redirect } from '@tanstack/react-router'
import { useAuthStore } from '@/shared/stores/auth-store'

export const Route = createFileRoute('/')({
  beforeLoad: () => {
    const { isAuthenticated } = useAuthStore.getState()
    if (isAuthenticated) {
      throw redirect({ to: '/products' })
    } else {
      throw redirect({ to: '/login' })
    }
  },
})
