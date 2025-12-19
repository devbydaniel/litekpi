import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const createProductSchema = z.object({
  name: z.string().min(1, 'Product name is required').max(255, 'Product name is too long'),
})

export type CreateProductFormValues = z.infer<typeof createProductSchema>

export function useCreateProductForm() {
  const form = useForm<CreateProductFormValues>({
    resolver: zodResolver(createProductSchema),
    defaultValues: { name: '' },
  })

  const reset = () => {
    form.reset({ name: '' })
  }

  return { form, reset }
}
