import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Card, CardContent } from '@/shared/components/ui/card'
import { useAuthStore } from '@/shared/stores/auth-store'

export function DashboardPage() {
  const user = useAuthStore((state) => state.user)

  return (
    <AuthenticatedLayout title="Dashboard">
      <div className="space-y-6">
        <Card>
          <CardContent className="p-4">
            <p className="text-sm text-muted-foreground">
              Welcome back,{' '}
              <span className="font-medium text-foreground">{user?.email}</span>
            </p>
          </CardContent>
        </Card>

        <Card className="border-dashed">
          <CardContent className="flex min-h-[400px] flex-col items-center justify-center p-8 text-center">
            <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-muted">
              <span className="text-2xl">📊</span>
            </div>
            <h3 className="mt-4 text-lg font-semibold">Coming soon</h3>
            <p className="mt-2 text-sm text-muted-foreground">
              Product management and KPI tracking features are under
              development.
            </p>
          </CardContent>
        </Card>
      </div>
    </AuthenticatedLayout>
  )
}
