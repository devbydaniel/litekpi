import { useAuthStore } from '@/shared/stores/auth-store'

export function useAuth() {
  const { user, token, isAuthenticated, setAuth, setUser, logout } =
    useAuthStore()

  return {
    user,
    token,
    isAuthenticated,
    setAuth,
    setUser,
    logout,
  }
}
