import { useQuery } from '@tanstack/react-query'
import { productsApi } from '@/shared/api/products'
import { measurementsApi } from '@/shared/api/measurements'

interface UseProductDetailOptions {
  productId: string
}

export function useProductDetail({ productId }: UseProductDetailOptions) {
  // Fetch product details
  const {
    data: product,
    isLoading: isLoadingProduct,
    error: productError,
  } = useQuery({
    queryKey: ['products', productId],
    queryFn: () => productsApi.get(productId),
  })

  // Fetch measurement names
  const {
    data: measurementsData,
    isLoading: isLoadingMeasurements,
    error: measurementsError,
  } = useQuery({
    queryKey: ['measurements', productId],
    queryFn: () => measurementsApi.listNames(productId),
    enabled: !!product,
  })

  return {
    product,
    measurements: measurementsData?.measurements ?? [],
    isLoading: isLoadingProduct || isLoadingMeasurements,
    error: productError || measurementsError,
  }
}
