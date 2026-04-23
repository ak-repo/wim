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
import { Select } from "@/components/ui/Select"
import { useCreateProduct, useUpdateProduct } from "@/features/products/hooks"
import { useProductCategories } from "@/features/productCategories/hooks"
import type { Product, CreateProductRequest, UpdateProductRequest } from "@/features/products/types"

interface ProductFormDialogProps {
  product: Product | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const ProductFormDialog: React.FC<ProductFormDialogProps> = ({
  product,
  open,
  onOpenChange,
}) => {
  const createProduct = useCreateProduct()
  const updateProduct = useUpdateProduct()
  const isEditing = !!product
  const { data: categoriesData } = useProductCategories({ page: 1, limit: 100 })

  const [formData, setFormData] = React.useState({
    sku: product?.sku || "",
    name: product?.name || "",
    description: product?.description || "",
    category: product?.category || "",
    unitOfMeasure: product?.unitOfMeasure || "unit",
    weight: product?.weight?.toString() || "",
    length: product?.length?.toString() || "",
    width: product?.width?.toString() || "",
    height: product?.height?.toString() || "",
    barcode: product?.barcode || "",
  })

  React.useEffect(() => {
    setFormData({
      sku: product?.sku || "",
      name: product?.name || "",
      description: product?.description || "",
    category: product?.category || "",
      unitOfMeasure: product?.unitOfMeasure || "unit",
      weight: product?.weight?.toString() || "",
      length: product?.length?.toString() || "",
      width: product?.width?.toString() || "",
      height: product?.height?.toString() || "",
      barcode: product?.barcode || "",
    })
  }, [product, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && product) {
        const updateData: UpdateProductRequest = {
          name: formData.name,
          description: formData.description || undefined,
          category: formData.category || undefined,
          unitOfMeasure: formData.unitOfMeasure,
          weight: formData.weight ? parseFloat(formData.weight) : undefined,
          length: formData.length ? parseFloat(formData.length) : undefined,
          width: formData.width ? parseFloat(formData.width) : undefined,
          height: formData.height ? parseFloat(formData.height) : undefined,
          barcode: formData.barcode || undefined,
        }
        await updateProduct.mutateAsync({ id: product.id, data: updateData })
      } else {
        const createData: CreateProductRequest = {
          sku: formData.sku,
          name: formData.name,
          description: formData.description || undefined,
          category: formData.category || undefined,
          unitOfMeasure: formData.unitOfMeasure,
          weight: formData.weight ? parseFloat(formData.weight) : undefined,
          length: formData.length ? parseFloat(formData.length) : undefined,
          width: formData.width ? parseFloat(formData.width) : undefined,
          height: formData.height ? parseFloat(formData.height) : undefined,
          barcode: formData.barcode || undefined,
        }
        await createProduct.mutateAsync(createData)
      }
      onOpenChange(false)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{isEditing ? "Edit Product" : "Create Product"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update product details."
              : "Add a new product to the catalog."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">SKU *</label>
              <Input
                name="sku"
                value={formData.sku}
                onChange={handleChange}
                placeholder="PROD-001"
                disabled={isEditing}
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Unit of Measure *</label>
              <Input
                name="unitOfMeasure"
                value={formData.unitOfMeasure}
                onChange={handleChange}
                placeholder="unit"
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Name *</label>
            <Input
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Product name"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Description</label>
            <Input
              name="description"
              value={formData.description}
              onChange={handleChange}
              placeholder="Product description"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Category</label>
            <Select
              name="category"
              value={formData.category}
              onChange={handleChange}
            >
              <option value="">Uncategorized</option>
              {categoriesData?.data?.map((category) => (
                <option key={category.id} value={category.name}>
                  {category.name}
                </option>
              ))}
            </Select>
          </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Barcode</label>
              <Input
                name="barcode"
                value={formData.barcode}
                onChange={handleChange}
                placeholder="Barcode"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Weight (kg)</label>
              <Input
                name="weight"
                type="number"
                step="0.01"
                value={formData.weight}
                onChange={handleChange}
                placeholder="0.00"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Length (cm)</label>
              <Input
                name="length"
                type="number"
                step="0.1"
                value={formData.length}
                onChange={handleChange}
                placeholder="0.0"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Width (cm)</label>
              <Input
                name="width"
                type="number"
                step="0.1"
                value={formData.width}
                onChange={handleChange}
                placeholder="0.0"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Height (cm)</label>
              <Input
                name="height"
                type="number"
                step="0.1"
                value={formData.height}
                onChange={handleChange}
                placeholder="0.0"
              />
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={createProduct.isPending || updateProduct.isPending}
            >
              {createProduct.isPending || updateProduct.isPending
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
