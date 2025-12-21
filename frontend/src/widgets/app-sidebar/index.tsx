import { Link, useNavigate } from '@tanstack/react-router'
import { Database, LayoutDashboard, LogOut, Users } from 'lucide-react'
import { useAuthStore } from '@/shared/stores/auth-store'
import { postAuthLogout } from '@/shared/api/generated/api'
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
      await postAuthLogout()
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
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 32 32"
            className="size-5 text-primary"
          >
            <rect x="4" y="20" width="6" height="8" rx="1" fill="currentColor" />
            <rect
              x="13"
              y="12"
              width="6"
              height="16"
              rx="1"
              fill="currentColor"
            />
            <rect x="22" y="4" width="6" height="24" rx="1" fill="currentColor" />
          </svg>
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
                <SidebarMenuButton asChild tooltip="Data Sources">
                  <Link to="/data-sources">
                    <Database className="size-4" />
                    <span>Data Sources</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
              {user?.role === 'admin' && (
                <SidebarMenuItem>
                  <SidebarMenuButton asChild tooltip="Users">
                    <Link to="/users">
                      <Users className="size-4" />
                      <span>Users</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              )}
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
