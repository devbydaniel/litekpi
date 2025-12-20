import { useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import {
  useGetProducts,
  usePostProducts,
  useDeleteProductsId,
  usePostProductsIdRegenerateKey,
  usePostProductsDemo,
  getGetProductsQueryKey,
} from '@/shared/api/generated/api'
import type { Product } from '@/shared/api/generated/models'

export function useProducts() {
  const queryClient = useQueryClient()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [regenerateDialogOpen, setRegenerateDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [productToDelete, setProductToDelete] = useState<Product | null>(null)
  const [apiKey, setApiKey] = useState<string | null>(null)

  const { data, isLoading } = useGetProducts()

  const createMutation = usePostProducts({
    mutation: {
      onSuccess: (response) => {
        queryClient.invalidateQueries({ queryKey: getGetProductsQueryKey() })
        setApiKey(response.apiKey ?? null)
        toast.success('Product created successfully')
      },
      onError: () => {
        toast.error('Failed to create product')
      },
    },
  })

  const deleteMutation = useDeleteProductsId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: getGetProductsQueryKey() })
        setDeleteDialogOpen(false)
        setProductToDelete(null)
        toast.success('Product deleted')
      },
      onError: () => {
        toast.error('Failed to delete product')
      },
    },
  })

  const regenerateMutation = usePostProductsIdRegenerateKey({
    mutation: {
      onSuccess: (response) => {
        setApiKey(response.apiKey ?? null)
        toast.success('API key regenerated')
      },
      onError: () => {
        toast.error('Failed to regenerate API key')
      },
    },
  })

  const createDemoMutation = usePostProductsDemo({
    mutation: {
      onSuccess: (response) => {
        queryClient.invalidateQueries({ queryKey: getGetProductsQueryKey() })
        setApiKey(response.apiKey ?? null)
        setCreateDialogOpen(true)
        toast.success('Demo product created successfully')
      },
      onError: () => {
        toast.error('Failed to create demo product')
      },
    },
  })

  const handleCreateProduct = async (name: string) => {
    await createMutation.mutateAsync({ data: { name } })
  }

  const handleDeleteProduct = (product: Product) => {
    setProductToDelete(product)
    setDeleteDialogOpen(true)
  }

  const confirmDeleteProduct = async () => {
    if (productToDelete?.id) {
      await deleteMutation.mutateAsync({ id: productToDelete.id })
    }
  }

  const handleRegenerateKey = (product: Product) => {
    setSelectedProduct(product)
    setRegenerateDialogOpen(true)
  }

  const confirmRegenerateKey = async () => {
    if (selectedProduct?.id) {
      await regenerateMutation.mutateAsync({ id: selectedProduct.id })
    }
  }

  const handleCreateDemo = async () => {
    await createDemoMutation.mutateAsync()
  }

  const closeDialogs = () => {
    setCreateDialogOpen(false)
    setRegenerateDialogOpen(false)
    setDeleteDialogOpen(false)
    setSelectedProduct(null)
    setProductToDelete(null)
    setApiKey(null)
  }

  return {
    products: data?.products ?? [],
    isLoading,
    createDialogOpen,
    setCreateDialogOpen,
    regenerateDialogOpen,
    deleteDialogOpen,
    selectedProduct,
    productToDelete,
    apiKey,
    isCreating: createMutation.isPending,
    isCreatingDemo: createDemoMutation.isPending,
    isRegenerating: regenerateMutation.isPending,
    isDeleting: deleteMutation.isPending,
    handleCreateProduct,
    handleCreateDemo,
    handleDeleteProduct,
    confirmDeleteProduct,
    handleRegenerateKey,
    confirmRegenerateKey,
    closeDialogs,
  }
}
