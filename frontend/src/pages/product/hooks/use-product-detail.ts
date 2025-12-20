import {
  useGetProductsId,
  useGetProductsProductIdMeasurements,
} from '@/shared/api/generated/api'

interface UseProductDetailOptions {
  productId: string
}

export function useProductDetail({ productId }: UseProductDetailOptions) {
  // Fetch product details
  const {
    data: product,
    isLoading: isLoadingProduct,
    error: productError,
  } = useGetProductsId(productId)

  // Fetch measurement names
  const {
    data: measurementsData,
    isLoading: isLoadingMeasurements,
    error: measurementsError,
  } = useGetProductsProductIdMeasurements(productId, {
    query: {
      enabled: !!product,
    },
  })

  return {
    product,
    measurements: measurementsData?.measurements ?? [],
    isLoading: isLoadingProduct || isLoadingMeasurements,
    error: productError || measurementsError,
  }
}
