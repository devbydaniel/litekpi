import { Plus } from 'lucide-react'
import { AuthenticatedLayout } from '@/layouts/authenticated'
import { Button } from '@/shared/components/ui/button'
import { useProducts } from './hooks/use-products'
import { ProductList } from './ui/product-list'
import { CreateProductDialog } from './ui/create-product-dialog'
import { RegenerateKeyDialog } from './ui/regenerate-key-dialog'

export function ProductsPage() {
  const {
    products,
    isLoading,
    createDialogOpen,
    setCreateDialogOpen,
    regenerateDialogOpen,
    selectedProduct,
    apiKey,
    isCreating,
    isRegenerating,
    handleCreateProduct,
    handleDeleteProduct,
    handleRegenerateKey,
    confirmRegenerateKey,
    closeDialogs,
  } = useProducts()

  return (
    <AuthenticatedLayout
      title="Products"
      actions={
        <Button onClick={() => setCreateDialogOpen(true)}>
          <Plus className="h-4 w-4" />
          New Product
        </Button>
      }
    >
      <ProductList
        products={products}
        isLoading={isLoading}
        onDelete={handleDeleteProduct}
        onRegenerateKey={handleRegenerateKey}
      />

      <CreateProductDialog
        open={createDialogOpen}
        onOpenChange={setCreateDialogOpen}
        onCreate={handleCreateProduct}
        apiKey={apiKey}
        isLoading={isCreating}
        onClose={closeDialogs}
      />

      <RegenerateKeyDialog
        open={regenerateDialogOpen}
        product={selectedProduct}
        apiKey={apiKey}
        isLoading={isRegenerating}
        onConfirm={confirmRegenerateKey}
        onClose={closeDialogs}
      />
    </AuthenticatedLayout>
  )
}
