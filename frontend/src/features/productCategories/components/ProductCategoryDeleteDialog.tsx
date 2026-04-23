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
import { useDeleteProductCategory } from "@/features/productCategories/hooks"
import type { ProductCategory } from "@/features/productCategories/types"
import { AlertTriangle } from "lucide-react"

interface ProductCategoryDeleteDialogProps {
  category: ProductCategory | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const ProductCategoryDeleteDialog: React.FC<ProductCategoryDeleteDialogProps> = ({
  category,
  open,
  onOpenChange,
}) => {
  const deleteCategory = useDeleteProductCategory()

  const handleDelete = async () => {
    if (!category) return
    try {
      await deleteCategory.mutateAsync(category.id)
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
            Delete Category
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this category? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        {category && (
          <div className="my-4 p-3 bg-muted rounded-md">
            <p className="text-sm font-medium">{category.name}</p>
            <p className="text-xs text-muted-foreground">{category.refCode}</p>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteCategory.isPending}
          >
            {deleteCategory.isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
