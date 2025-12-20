import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { productsApi } from '@/shared/api/products'
import type { Product } from '@/shared/types'

export function useProducts() {
  const queryClient = useQueryClient()
  const [createDialogOpen, setCreateDialogOpen] = useState(false)
  const [regenerateDialogOpen, setRegenerateDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [productToDelete, setProductToDelete] = useState<Product | null>(null)
  const [apiKey, setApiKey] = useState<string | null>(null)

  const { data, isLoading } = useQuery({
    queryKey: ['products'],
    queryFn: () => productsApi.list(),
  })

  const createMutation = useMutation({
    mutationFn: productsApi.create,
    onSuccess: (response) => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      setApiKey(response.apiKey)
      toast.success('Product created successfully')
    },
    onError: () => {
      toast.error('Failed to create product')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: productsApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      setDeleteDialogOpen(false)
      setProductToDelete(null)
      toast.success('Product deleted')
    },
    onError: () => {
      toast.error('Failed to delete product')
    },
  })

  const regenerateMutation = useMutation({
    mutationFn: productsApi.regenerateKey,
    onSuccess: (response) => {
      setApiKey(response.apiKey)
      toast.success('API key regenerated')
    },
    onError: () => {
      toast.error('Failed to regenerate API key')
    },
  })

  const createDemoMutation = useMutation({
    mutationFn: productsApi.createDemo,
    onSuccess: (response) => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
      setApiKey(response.apiKey)
      setCreateDialogOpen(true)
      toast.success('Demo product created successfully')
    },
    onError: () => {
      toast.error('Failed to create demo product')
    },
  })

  const handleCreateProduct = async (name: string) => {
    await createMutation.mutateAsync({ name })
  }

  const handleDeleteProduct = (product: Product) => {
    setProductToDelete(product)
    setDeleteDialogOpen(true)
  }

  const confirmDeleteProduct = async () => {
    if (productToDelete) {
      await deleteMutation.mutateAsync(productToDelete.id)
    }
  }

  const handleRegenerateKey = (product: Product) => {
    setSelectedProduct(product)
    setRegenerateDialogOpen(true)
  }

  const confirmRegenerateKey = async () => {
    if (selectedProduct) {
      await regenerateMutation.mutateAsync(selectedProduct.id)
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
