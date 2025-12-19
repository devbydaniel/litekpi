import { CheckCircle2, XCircle } from "lucide-react"

import { cn } from "@/shared/lib/utils"
import { Spinner } from "@/shared/components/ui/spinner"

type StatusType = "success" | "error" | "loading"

interface StatusCardProps {
  status: StatusType
  title: string
  description?: string
  action?: React.ReactNode
  className?: string
}

const statusConfig: Record<
  StatusType,
  { icon: React.ReactNode; iconContainerClass: string }
> = {
  success: {
    icon: <CheckCircle2 className="h-6 w-6 text-green-600" />,
    iconContainerClass: "bg-green-100",
  },
  error: {
    icon: <XCircle className="h-6 w-6 text-destructive" />,
    iconContainerClass: "bg-destructive/10",
  },
  loading: {
    icon: <Spinner size="default" />,
    iconContainerClass: "bg-primary/10",
  },
}

export function StatusCard({
  status,
  title,
  description,
  action,
  className,
}: StatusCardProps) {
  const config = statusConfig[status]

  return (
    <div className={cn("text-center", className)}>
      <div
        className={cn(
          "mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full",
          config.iconContainerClass
        )}
      >
        {config.icon}
      </div>
      <h2 className="mb-2 text-xl font-semibold">{title}</h2>
      {description && (
        <p className="mb-4 text-sm text-muted-foreground">{description}</p>
      )}
      {action && <div className="mt-4">{action}</div>}
    </div>
  )
}
