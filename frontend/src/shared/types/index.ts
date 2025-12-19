// Re-export common types
export type { User } from '@/shared/stores/auth-store'

// API response types
export interface ApiResponse<T> {
  data: T
  message?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// Product types
export interface Product {
  id: string
  name: string
  organizationId: string
  createdAt: string
  updatedAt: string
}

// Measurement types
export interface Measurement {
  id: string
  productId: string
  name: string
  value: number
  timestamp: string
  metadata: Record<string, string> | null
  createdAt: string
}

// Data point types (deprecated - use Measurement)
export interface DataPoint {
  id: string
  productId: string
  metric: string
  value: number
  timestamp: string
  tags: Record<string, string>
  createdAt: string
}

// Dashboard types
export interface Dashboard {
  id: string
  productId: string
  name: string
  config: DashboardConfig
  createdAt: string
  updatedAt: string
}

export interface DashboardConfig {
  widgets: Widget[]
  layout: LayoutItem[]
}

export interface Widget {
  id: string
  type: 'chart' | 'metric' | 'table'
  title: string
  config: Record<string, unknown>
}

export interface LayoutItem {
  widgetId: string
  x: number
  y: number
  w: number
  h: number
}

// Team types
export interface TeamMember {
  id: string
  userId: string
  productId: string
  role: 'owner' | 'admin' | 'viewer'
  user: {
    id: string
    email: string
  }
  createdAt: string
}
