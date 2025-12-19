import { api } from './client'
import type { User } from '@/shared/types'

export interface AuthResponse {
  user: User
  token: string
}

export interface RegisterResponse {
  message: string
  user: User
}

export interface MessageResponse {
  message: string
}

export interface RegisterData {
  email: string
  password: string
  name: string
  organizationName: string
}

export interface LoginData {
  email: string
  password: string
}

export interface ResetPasswordData {
  token: string
  newPassword: string
}

export interface CompleteOAuthSetupData {
  token: string
  name: string
  organizationName: string
}

export const authApi = {
  register(data: RegisterData): Promise<RegisterResponse> {
    return api.post('/auth/register', data)
  },

  login(data: LoginData): Promise<AuthResponse> {
    return api.post('/auth/login', data)
  },

  verifyEmail(token: string): Promise<MessageResponse> {
    return api.post('/auth/verify-email', { token })
  },

  forgotPassword(email: string): Promise<MessageResponse> {
    return api.post('/auth/forgot-password', { email })
  },

  resetPassword(data: ResetPasswordData): Promise<MessageResponse> {
    return api.post('/auth/reset-password', data)
  },

  resendVerification(email: string): Promise<MessageResponse> {
    return api.post('/auth/resend-verification', { email })
  },

  completeOAuthSetup(data: CompleteOAuthSetupData): Promise<AuthResponse> {
    return api.post('/auth/complete-oauth-setup', data)
  },

  me(): Promise<User> {
    return api.get('/auth/me')
  },

  logout(): Promise<MessageResponse> {
    return api.post('/auth/logout')
  },

  getGoogleAuthUrl(): string {
    return '/api/v1/auth/google'
  },

  getGithubAuthUrl(): string {
    return '/api/v1/auth/github'
  },
}
