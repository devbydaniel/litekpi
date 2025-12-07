import { useAuthStore } from '@/shared/stores/auth-store'

const API_BASE_URL = '/api/v1'

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public data?: unknown
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

interface RequestOptions extends RequestInit {
  params?: Record<string, string>
}

async function request<T>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<T> {
  const { params, ...fetchOptions } = options
  const token = useAuthStore.getState().token

  // Build URL with query params
  let url = `${API_BASE_URL}${endpoint}`
  if (params) {
    const searchParams = new URLSearchParams(params)
    url += `?${searchParams.toString()}`
  }

  // Set default headers
  const headers = new Headers(fetchOptions.headers)
  if (!headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json')
  }
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(url, {
    ...fetchOptions,
    headers,
  })

  // Handle non-OK responses
  if (!response.ok) {
    let errorData: unknown
    try {
      errorData = await response.json()
    } catch {
      errorData = await response.text()
    }

    // Handle 401 by logging out
    if (response.status === 401) {
      useAuthStore.getState().logout()
    }

    throw new ApiError(
      response.status,
      `HTTP ${response.status}: ${response.statusText}`,
      errorData
    )
  }

  // Handle empty responses
  if (response.status === 204) {
    return undefined as T
  }

  return response.json()
}

export const api = {
  get<T>(endpoint: string, options?: RequestOptions): Promise<T> {
    return request<T>(endpoint, { ...options, method: 'GET' })
  },

  post<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
    return request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  },

  put<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
    return request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  },

  patch<T>(endpoint: string, data?: unknown, options?: RequestOptions): Promise<T> {
    return request<T>(endpoint, {
      ...options,
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    })
  },

  delete<T>(endpoint: string, options?: RequestOptions): Promise<T> {
    return request<T>(endpoint, { ...options, method: 'DELETE' })
  },
}

// Custom instance for Orval-generated clients
export type CustomInstanceConfig = {
  url: string
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
  params?: Record<string, string>
  data?: unknown
  headers?: Record<string, string>
  signal?: AbortSignal
}

export const customInstance = async <T>({
  url,
  method,
  params,
  data,
  headers,
  signal,
}: CustomInstanceConfig): Promise<T> => {
  const token = useAuthStore.getState().token

  // Build URL with query params
  let fullUrl = `${API_BASE_URL}${url}`
  if (params) {
    const searchParams = new URLSearchParams(params)
    fullUrl += `?${searchParams.toString()}`
  }

  // Set headers
  const requestHeaders = new Headers(headers)
  if (!requestHeaders.has('Content-Type') && data) {
    requestHeaders.set('Content-Type', 'application/json')
  }
  if (token) {
    requestHeaders.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(fullUrl, {
    method,
    headers: requestHeaders,
    body: data ? JSON.stringify(data) : undefined,
    signal,
  })

  // Handle non-OK responses
  if (!response.ok) {
    let errorData: unknown
    try {
      errorData = await response.json()
    } catch {
      errorData = await response.text()
    }

    // Handle 401 by logging out
    if (response.status === 401) {
      useAuthStore.getState().logout()
    }

    throw new ApiError(
      response.status,
      `HTTP ${response.status}: ${response.statusText}`,
      errorData
    )
  }

  // Handle empty responses
  if (response.status === 204) {
    return undefined as T
  }

  return response.json()
}
