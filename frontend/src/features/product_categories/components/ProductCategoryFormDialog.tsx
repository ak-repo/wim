import * as React from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/Dialog"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { useCreateProductCategory, useUpdateProductCategory } from "@/features/product_categories/hooks"
import type { ProductCategory, CreateProductCategoryRequest } from "@/features/product_categories/types"

interface ProductCategoryFormDialogProps {
  category: ProductCategory | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const ProductCategoryFormDialog: React.FC<ProductCategoryFormDialogProps> = ({
  category,
  open,
  onOpenChange,
}) => {
  const createCategory = useCreateProductCategory()
  const updateCategory = useUpdateProductCategory()
  const isEditing = !!category

  const [formData, setFormData] = React.useState({
    name: category?.name || "",
    description: category?.description || "",
    isActive: category?.isActive ?? true,
  })

  React.useEffect(() => {
    setFormData({
      name: category?.name || "",
      description: category?.description || "",
      isActive: category?.isActive ?? true,
    })
  }, [category, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.checked }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && category) {
        await updateCategory.mutateAsync({
          id: category.id,
          data: {
            name: formData.name,
            description: formData.description || undefined,
            isActive: formData.isActive,
          },
        })
      } else {
        await createCategory.mutateAsync({
          name: formData.name,
          description: formData.description || undefined,
          isActive: formData.isActive,
        } as CreateProductCategoryRequest)
      }
      onOpenChange(false)
    } catch {
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEditing ? "Edit Product Category" : "Create Product Category"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update product category details."
              : "Add a new product category to the system."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Name</label>
            <Input
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Electronics"
              required
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Description</label>
            <textarea
              name="description"
              value={formData.description}
              onChange={handleChange}
              placeholder="Category description..."
              className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              name="isActive"
              id="isActive"
              checked={formData.isActive}
              onChange={handleCheckboxChange}
              className="h-4 w-4 rounded border-input"
            />
            <label htmlFor="isActive" className="text-sm">Active</label>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={createCategory.isPending || updateCategory.isPending}>
              {createCategory.isPending || updateCategory.isPending
                ? "Saving..."
                : isEditing
                ? "Update"
                : "Create"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}