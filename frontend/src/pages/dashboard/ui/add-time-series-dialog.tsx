import { useState, useEffect } from 'react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog'
import { Button } from '@/shared/components/ui/button'
import { Input } from '@/shared/components/ui/input'
import { Label } from '@/shared/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select'
import { useGetDataSources, useGetDataSourcesDataSourceIdMeasurements } from '@/shared/api/generated/api'
import type { CreateTimeSeriesRequest } from '@/shared/api/generated/models'

interface AddTimeSeriesDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onAdd: (timeSeries: CreateTimeSeriesRequest) => Promise<void>
  isLoading: boolean
}

export function AddTimeSeriesDialog({
  open,
  onOpenChange,
  onAdd,
  isLoading,
}: AddTimeSeriesDialogProps) {
  const [dataSourceId, setDataSourceId] = useState<string>('')
  const [measurementName, setMeasurementName] = useState<string>('')
  const [title, setTitle] = useState<string>('')

  // Fetch data sources
  const { data: dataSourcesData } = useGetDataSources()
  const dataSources = dataSourcesData?.dataSources ?? []

  // Fetch measurements for selected data source
  const { data: measurementsData } = useGetDataSourcesDataSourceIdMeasurements(
    dataSourceId,
    { query: { enabled: !!dataSourceId } }
  )
  const measurements = measurementsData?.measurements ?? []

  // Reset measurement when data source changes
  useEffect(() => {
    setMeasurementName('')
  }, [dataSourceId])

  // Reset form when dialog closes
  useEffect(() => {
    if (!open) {
      setDataSourceId('')
      setMeasurementName('')
      setTitle('')
    }
  }, [open])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!dataSourceId || !measurementName) return

    await onAdd({
      dataSourceId,
      measurementName,
      title: title || undefined,
      chartType: 'area',
      dateRange: 'last_7_days',
    })
    onOpenChange(false)
  }

  const canSubmit = dataSourceId && measurementName

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add Chart</DialogTitle>
          <DialogDescription>
            Select a measurement to display on your dashboard. You can configure
            chart type, filters, and other options after adding the chart.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="dataSource">Data Source</Label>
            <Select value={dataSourceId} onValueChange={setDataSourceId}>
              <SelectTrigger id="dataSource">
                <SelectValue placeholder="Select data source" />
              </SelectTrigger>
              <SelectContent>
                {dataSources.map((ds) => (
                  <SelectItem key={ds.id} value={ds.id ?? ''}>
                    {ds.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="measurement">Measurement</Label>
            <Select
              value={measurementName}
              onValueChange={setMeasurementName}
              disabled={!dataSourceId || measurements.length === 0}
            >
              <SelectTrigger id="measurement">
                <SelectValue placeholder={
                  !dataSourceId
                    ? 'Select a data source first'
                    : measurements.length === 0
                    ? 'No measurements available'
                    : 'Select measurement'
                } />
              </SelectTrigger>
              <SelectContent>
                {measurements.map((m) => (
                  <SelectItem key={m.name} value={m.name ?? ''}>
                    {m.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="title">Title (optional)</Label>
            <Input
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder={measurementName || 'Enter a custom title'}
              maxLength={128}
            />
            <p className="text-xs text-muted-foreground">
              Defaults to measurement name if left empty
            </p>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading || !canSubmit}>
              {isLoading ? 'Adding...' : 'Add Chart'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
