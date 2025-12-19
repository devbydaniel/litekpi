import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type Role = 'admin' | 'editor' | 'viewer'

export interface Organization {
  id: string
  name: string
  createdAt: string
  updatedAt: string
}

export interface User {
  id: string
  email: string
  name: string
  emailVerified: boolean
  organizationId: string
  organization?: Organization
  role: Role
  createdAt: string
  updatedAt: string
}

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  setAuth: (user: User, token: string) => void
  setUser: (user: User) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,

      setAuth: (user, token) =>
        set({
          user,
          token,
          isAuthenticated: true,
        }),

      setUser: (user) =>
        set({
          user,
        }),

      logout: () =>
        set({
          user: null,
          token: null,
          isAuthenticated: false,
        }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
)
