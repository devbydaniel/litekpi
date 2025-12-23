import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const createMCPApiKeySchema = z.object({
  name: z
    .string()
    .min(1, 'API key name is required')
    .max(255, 'Name is too long'),
  dataSourceIds: z
    .array(z.string())
    .min(1, 'At least one data source must be selected'),
})

export type CreateMCPApiKeyFormValues = z.infer<typeof createMCPApiKeySchema>

export function useCreateMCPApiKeyForm() {
  const form = useForm<CreateMCPApiKeyFormValues>({
    resolver: zodResolver(createMCPApiKeySchema),
    defaultValues: {
      name: '',
      dataSourceIds: [],
    },
  })

  const reset = () => {
    form.reset({ name: '', dataSourceIds: [] })
  }

  return { form, reset }
}
