import { TanStackRouterDevtools } from '@tanstack/router-devtools'
import { Toaster } from '@/shared/components/ui/sonner'

interface RootLayoutProps {
  children: React.ReactNode
}

export function RootLayout({ children }: RootLayoutProps) {
  return (
    <>
      {children}
      <Toaster />
      {false && import.meta.env.DEV && <TanStackRouterDevtools />}
    </>
  )
}
