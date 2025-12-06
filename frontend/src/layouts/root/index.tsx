import { TanStackRouterDevtools } from '@tanstack/router-devtools'

interface RootLayoutProps {
  children: React.ReactNode
}

export function RootLayout({ children }: RootLayoutProps) {
  return (
    <>
      {children}
      {import.meta.env.DEV && <TanStackRouterDevtools />}
    </>
  )
}
