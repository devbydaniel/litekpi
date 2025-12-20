import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Card, CardContent } from '@/shared/components/ui/card'
import { EmptyState } from '@/shared/components/ui/empty-state'
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

        <EmptyState
          icon="📊"
          title="Coming soon"
          description="I'm not sure yet what exactly to display here. Suggestions? Let me know!"
          className="min-h-[400px]"
        />
      </div>
    </AuthenticatedLayout>
  )
}
