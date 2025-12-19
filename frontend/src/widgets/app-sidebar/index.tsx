import { Link, useNavigate } from '@tanstack/react-router'
import { LayoutDashboard, LogOut, Package } from 'lucide-react'
import { useAuthStore } from '@/shared/stores/auth-store'
import { authApi } from '@/shared/api/auth'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/shared/components/ui/sidebar'

export function AppSidebar() {
  const { user, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = async () => {
    try {
      await authApi.logout()
    } catch {
      // Continue with logout even if API call fails
    }
    logout()
    navigate({ to: '/login' })
  }

  return (
    <Sidebar>
      <SidebarHeader>
        <Link
          to="/"
          className="flex items-center gap-2 px-2 py-1 font-semibold"
        >
          <span className="text-xl">📊</span>
          <span className="font-heading">LiteKPI</span>
        </Link>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton asChild tooltip="Dashboard">
                  <Link to="/">
                    <LayoutDashboard className="size-4" />
                    <span>Dashboard</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              <SidebarMenuItem>
                <SidebarMenuButton asChild tooltip="Products">
                  <Link to="/products">
                    <Package className="size-4" />
                    <span>Products</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="border-t border-sidebar-border">
        <SidebarMenu>
          <SidebarMenuItem>
            <div className="flex items-center gap-2 px-2 py-1">
              <span className="flex-1 truncate text-sm">{user?.email}</span>
              <button
                onClick={handleLogout}
                className="text-muted-foreground transition-colors hover:text-foreground"
                title="Logout"
              >
                <LogOut className="size-4" />
              </button>
            </div>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  )
}
