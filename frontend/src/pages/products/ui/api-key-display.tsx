import { useState } from 'react'
import { Copy, Check } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'

interface ApiKeyDisplayProps {
  apiKey: string
}

export function ApiKeyDisplay({ apiKey }: ApiKeyDisplayProps) {
  const [copied, setCopied] = useState(false)

  const copyToClipboard = () => {
    navigator.clipboard.writeText(apiKey)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="flex gap-2">
      <Input
        value={apiKey}
        readOnly
        className="font-mono text-sm"
      />
      <Button
        type="button"
        variant="outline"
        size="icon"
        onClick={copyToClipboard}
      >
        {copied ? (
          <Check className="h-4 w-4" />
        ) : (
          <Copy className="h-4 w-4" />
        )}
      </Button>
    </div>
  )
}
