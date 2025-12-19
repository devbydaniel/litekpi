import { MoreHorizontal, Trash, RefreshCw, Package } from 'lucide-react'
import { Button } from '@/shared/components/ui/button'
import { EmptyState } from '@/shared/components/ui/empty-state'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/components/ui/dropdown-menu'
import {
  Item,
  ItemActions,
  ItemContent,
  ItemDescription,
  ItemGroup,
  ItemSeparator,
  ItemTitle,
} from '@/shared/components/ui/item'
import { Skeleton } from '@/shared/components/ui/skeleton'
import type { Product } from '@/shared/types'

interface ProductListProps {
  products: Product[]
  isLoading: boolean
  onDelete: (id: string) => void
  onRegenerateKey: (product: Product) => void
}

export function ProductList({
  products,
  isLoading,
  onDelete,
  onRegenerateKey,
}: ProductListProps) {
  if (isLoading) {
    return <ProductListSkeleton />
  }

  if (products.length === 0) {
    return (
      <EmptyState
        icon={Package}
        title="No products yet"
        description="Create your first product to start tracking KPIs."
      />
    )
  }

  return (
    <ItemGroup className="rounded-lg border">
      {products.map((product, index) => (
        <div key={product.id}>
          {index > 0 && <ItemSeparator />}
          <ProductListItem
            product={product}
            onDelete={onDelete}
            onRegenerateKey={onRegenerateKey}
          />
        </div>
      ))}
    </ItemGroup>
  )
}

interface ProductListItemProps {
  product: Product
  onDelete: (id: string) => void
  onRegenerateKey: (product: Product) => void
}

function ProductListItem({
  product,
  onDelete,
  onRegenerateKey,
}: ProductListItemProps) {
  return (
    <Item>
      <ItemContent>
        <ItemTitle>{product.name}</ItemTitle>
        <ItemDescription>
          Created {new Date(product.createdAt).toLocaleDateString()}
        </ItemDescription>
      </ItemContent>

      <ItemActions>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon">
              <MoreHorizontal className="h-4 w-4" />
              <span className="sr-only">Open menu</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={() => onRegenerateKey(product)}>
              <RefreshCw className="mr-2 h-4 w-4" />
              Regenerate API Key
            </DropdownMenuItem>
            <DropdownMenuItem
              className="text-destructive focus:text-destructive"
              onClick={() => onDelete(product.id)}
            >
              <Trash className="mr-2 h-4 w-4" />
              Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </ItemActions>
    </Item>
  )
}

function ProductListSkeleton() {
  return (
    <ItemGroup className="rounded-lg border">
      {[1, 2, 3].map((i) => (
        <div key={i}>
          {i > 1 && <ItemSeparator />}
          <Item>
            <ItemContent>
              <Skeleton className="h-5 w-32" />
              <Skeleton className="h-4 w-24" />
            </ItemContent>
            <ItemActions>
              <Skeleton className="h-8 w-8" />
            </ItemActions>
          </Item>
        </div>
      ))}
    </ItemGroup>
  )
}

