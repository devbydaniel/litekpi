import { useCallback, type RefObject } from 'react'
import html2canvas from 'html2canvas'
import { toast } from 'sonner'

interface UseChartExportOptions {
  chartRef: RefObject<HTMLDivElement | null>
  filename?: string
}

export function useChartExport({ chartRef, filename = 'chart' }: UseChartExportOptions) {
  const captureChart = useCallback(async () => {
    if (!chartRef.current) {
      throw new Error('Chart element not found')
    }

    return html2canvas(chartRef.current, {
      backgroundColor: null,
      scale: 2,
    })
  }, [chartRef])

  const copyToClipboard = useCallback(async () => {
    try {
      const canvas = await captureChart()
      const blob = await new Promise<Blob>((resolve, reject) => {
        canvas.toBlob((b) => {
          if (b) resolve(b)
          else reject(new Error('Failed to create blob'))
        }, 'image/png')
      })

      await navigator.clipboard.write([
        new ClipboardItem({ 'image/png': blob }),
      ])
      toast.success('Chart copied to clipboard')
    } catch {
      toast.error('Failed to copy chart')
    }
  }, [captureChart])

  const downloadAsPng = useCallback(async () => {
    try {
      const canvas = await captureChart()
      const link = document.createElement('a')
      link.download = `${filename}.png`
      link.href = canvas.toDataURL('image/png')
      link.click()
      toast.success('Chart downloaded')
    } catch {
      toast.error('Failed to download chart')
    }
  }, [captureChart, filename])

  return {
    copyToClipboard,
    downloadAsPng,
  }
}
