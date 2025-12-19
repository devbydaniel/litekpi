import { Link, useParams } from '@tanstack/react-router'
import { ArrowLeft, BarChart3 } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import { EmptyState } from '@/shared/components/ui/empty-state'
import { Skeleton } from '@/shared/components/ui/skeleton'
import { useProductDetail } from './hooks/use-product-detail'
import { MeasurementChart } from './ui/measurement-chart'

export function ProductDetailPage() {
  const { productId } = useParams({ from: '/_authenticated/products_/$productId' })
  const { product, measurements, isLoading, error } = useProductDetail({ productId })

  if (error) {
    return (
      <AuthenticatedLayout title="Product">
        <div className="flex flex-col items-center justify-center gap-4 py-12">
          <p className="text-muted-foreground">Failed to load product</p>
          <Button asChild variant="outline">
            <Link to="/products">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Products
            </Link>
          </Button>
        </div>
      </AuthenticatedLayout>
    )
  }

  return (
    <AuthenticatedLayout
      title={isLoading ? 'Loading...' : product?.name ?? 'Product'}
      actions={
        <Button asChild variant="outline">
          <Link to="/products">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Products
          </Link>
        </Button>
      }
    >
      {isLoading ? (
        <div className="space-y-6">
          <Skeleton className="h-[400px] w-full" />
          <Skeleton className="h-[400px] w-full" />
        </div>
      ) : measurements.length === 0 ? (
        <EmptyState
          icon={BarChart3}
          title="No measurements yet"
          description="Start sending measurement data to this product using its API key."
        />
      ) : (
        <div className="space-y-6">
          {measurements.map((measurement) => (
            <MeasurementChart
              key={measurement.name}
              productId={productId}
              measurement={measurement}
            />
          ))}
        </div>
      )}
    </AuthenticatedLayout>
  )
}
