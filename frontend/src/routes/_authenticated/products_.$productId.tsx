import { createFileRoute } from '@tanstack/react-router'
import { ProductDetailPage } from '@/pages/product'

export const Route = createFileRoute('/_authenticated/products_/$productId')({
  component: ProductDetailPage,
})
