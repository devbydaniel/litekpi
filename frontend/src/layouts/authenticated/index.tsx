import { Link } from '@tanstack/react-router'
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from '@/shared/components/ui/sidebar'
import { AppSidebar } from '@/widgets/app-sidebar'

interface AuthenticatedLayoutProps {
  children: React.ReactNode
  title?: React.ReactNode
  actions?: React.ReactNode
}

export function AuthenticatedLayout({
  children,
  title,
  actions,
}: AuthenticatedLayoutProps) {
  return (
    <SidebarProvider>
      {/* Accent bar */}
      <div className="hidden w-1 bg-primary md:block" />

      <AppSidebar />

      <SidebarInset>
        {/* Mobile header */}
        <header className="flex h-14 items-center gap-2 border-b px-4 md:hidden">
          <SidebarTrigger />
          <Link to="/" className="flex items-center gap-2 font-semibold">
            <span className="text-xl">📊</span>
            <span>LiteKPI</span>
          </Link>
        </header>

        {/* Desktop top bar */}
        {(title || actions) && (
          <header className="hidden h-14 items-center justify-between border-b px-6 md:flex">
            {title && <h1 className="font-semibold">{title}</h1>}
            {actions && <div className="flex items-center gap-2">{actions}</div>}
          </header>
        )}

        <main className="flex-1 p-4 md:p-6">{children}</main>
      </SidebarInset>
    </SidebarProvider>
  )
}
