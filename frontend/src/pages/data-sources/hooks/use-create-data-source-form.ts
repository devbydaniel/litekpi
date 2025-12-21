import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const createDataSourceSchema = z.object({
  name: z.string().min(1, 'Data source name is required').max(255, 'Data source name is too long'),
})

export type CreateDataSourceFormValues = z.infer<typeof createDataSourceSchema>

export function useCreateDataSourceForm() {
  const form = useForm<CreateDataSourceFormValues>({
    resolver: zodResolver(createDataSourceSchema),
    defaultValues: { name: '' },
  })

  const reset = () => {
    form.reset({ name: '' })
  }

  return { form, reset }
}
