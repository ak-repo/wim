import * as React from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/Dialog"
import { Button } from "@/components/ui/Button"
import { useDeleteProduct } from "@/features/products/hooks"
import type { Product } from "@/features/products/types"
import { AlertTriangle } from "lucide-react"

interface ProductDeleteDialogProps {
  product: Product | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const ProductDeleteDialog: React.FC<ProductDeleteDialogProps> = ({
  product,
  open,
  onOpenChange,
}) => {
  const deleteProduct = useDeleteProduct()

  const handleDelete = async () => {
    if (!product) return
    try {
      await deleteProduct.mutateAsync(product.id)
      onOpenChange(false)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-destructive" />
            Delete Product
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this product? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        {product && (
          <div className="my-4 p-3 bg-muted rounded-md">
            <p className="text-sm font-medium">{product.name}</p>
            <p className="text-xs text-muted-foreground">SKU: {product.sku}</p>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteProduct.isPending}
          >
            {deleteProduct.isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
