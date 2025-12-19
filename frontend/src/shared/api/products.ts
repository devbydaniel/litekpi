import { api } from './client'
import type { Product } from '@/shared/types'

export interface CreateProductRequest {
  name: string
}

export interface CreateProductResponse {
  product: Product
  apiKey: string
}

export interface RegenerateKeyResponse {
  apiKey: string
}

export interface ListProductsResponse {
  products: Product[]
}

export const productsApi = {
  list(): Promise<ListProductsResponse> {
    return api.get('/products')
  },

  get(id: string): Promise<Product> {
    return api.get(`/products/${id}`)
  },

  create(data: CreateProductRequest): Promise<CreateProductResponse> {
    return api.post('/products', data)
  },

  createDemo(): Promise<CreateProductResponse> {
    return api.post('/products/demo')
  },

  delete(id: string): Promise<void> {
    return api.delete(`/products/${id}`)
  },

  regenerateKey(id: string): Promise<RegenerateKeyResponse> {
    return api.post(`/products/${id}/regenerate-key`)
  },
}
