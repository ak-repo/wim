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
import {
  useCreateProductCategory,
  useUpdateProductCategory,
} from "@/features/productCategories/hooks"
import type {
  ProductCategory,
  CreateProductCategoryRequest,
  UpdateProductCategoryRequest,
} from "@/features/productCategories/types"

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
    isActive: category?.isActive ?? true,
  })

  React.useEffect(() => {
    setFormData({
      name: category?.name || "",
      isActive: category?.isActive ?? true,
    })
  }, [category, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.checked }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && category) {
        const updateData: UpdateProductCategoryRequest = {
          name: formData.name,
          isActive: formData.isActive,
        }
        await updateCategory.mutateAsync({ id: category.id, data: updateData })
      } else {
        const createData: CreateProductCategoryRequest = {
          name: formData.name,
          isActive: formData.isActive,
        }
        await createCategory.mutateAsync(createData)
      }
      onOpenChange(false)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEditing ? "Edit Category" : "Create Category"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update product category details."
              : "Add a new product category."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Name</label>
            <Input
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Category name"
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              name="isActive"
              id="categoryActive"
              checked={formData.isActive}
              onChange={handleCheckboxChange}
              className="h-4 w-4 rounded border-input"
            />
            <label htmlFor="categoryActive" className="text-sm">Active</label>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
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
