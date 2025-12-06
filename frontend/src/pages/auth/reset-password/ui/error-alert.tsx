interface ErrorAlertProps {
  message: string
}

export function ErrorAlert({ message }: ErrorAlertProps) {
  return (
    <div className="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
      {message}
    </div>
  )
}
