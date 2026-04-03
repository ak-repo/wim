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
  useCreateWarehouse,
  useUpdateWarehouse,
} from "@/features/warehouses/hooks"
import type { Warehouse, CreateWarehouseRequest, UpdateWarehouseRequest } from "@/features/warehouses/types"

interface WarehouseFormDialogProps {
  warehouse: Warehouse | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const WarehouseFormDialog: React.FC<WarehouseFormDialogProps> = ({
  warehouse,
  open,
  onOpenChange,
}) => {
  const createWarehouse = useCreateWarehouse()
  const updateWarehouse = useUpdateWarehouse()
  const isEditing = !!warehouse

  const [formData, setFormData] = React.useState({
    code: warehouse?.code || "",
    name: warehouse?.name || "",
    addressLine1: warehouse?.addressLine1 || "",
    addressLine2: warehouse?.addressLine2 || "",
    city: warehouse?.city || "",
    state: warehouse?.state || "",
    postalCode: warehouse?.postalCode || "",
    country: warehouse?.country || "",
  })

  React.useEffect(() => {
    setFormData({
      code: warehouse?.code || "",
      name: warehouse?.name || "",
      addressLine1: warehouse?.addressLine1 || "",
      addressLine2: warehouse?.addressLine2 || "",
      city: warehouse?.city || "",
      state: warehouse?.state || "",
      postalCode: warehouse?.postalCode || "",
      country: warehouse?.country || "",
    })
  }, [warehouse, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && warehouse) {
        const updateData: UpdateWarehouseRequest = {
          name: formData.name,
          addressLine1: formData.addressLine1 || undefined,
          addressLine2: formData.addressLine2 || undefined,
          city: formData.city || undefined,
          state: formData.state || undefined,
          postalCode: formData.postalCode || undefined,
          country: formData.country,
        }
        await updateWarehouse.mutateAsync({ id: warehouse.id, data: updateData })
      } else {
        const createData: CreateWarehouseRequest = {
          code: formData.code,
          name: formData.name,
          addressLine1: formData.addressLine1 || undefined,
          addressLine2: formData.addressLine2 || undefined,
          city: formData.city || undefined,
          state: formData.state || undefined,
          postalCode: formData.postalCode || undefined,
          country: formData.country,
        }
        await createWarehouse.mutateAsync(createData)
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
          <DialogTitle>{isEditing ? "Edit Warehouse" : "Create Warehouse"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update warehouse details."
              : "Add a new warehouse to the system."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Code *</label>
              <Input
                name="code"
                value={formData.code}
                onChange={handleChange}
                placeholder="WH-001"
                disabled={isEditing}
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Country *</label>
              <Input
                name="country"
                value={formData.country}
                onChange={handleChange}
                placeholder="Country"
              />
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Name *</label>
            <Input
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Warehouse name"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Address Line 1</label>
            <Input
              name="addressLine1"
              value={formData.addressLine1}
              onChange={handleChange}
              placeholder="Street address"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Address Line 2</label>
            <Input
              name="addressLine2"
              value={formData.addressLine2}
              onChange={handleChange}
              placeholder="Apt, suite, etc."
            />
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">City</label>
              <Input
                name="city"
                value={formData.city}
                onChange={handleChange}
                placeholder="City"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">State</label>
              <Input
                name="state"
                value={formData.state}
                onChange={handleChange}
                placeholder="State"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Postal Code</label>
              <Input
                name="postalCode"
                value={formData.postalCode}
                onChange={handleChange}
                placeholder="Postal code"
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
              disabled={createWarehouse.isPending || updateWarehouse.isPending}
            >
              {createWarehouse.isPending || updateWarehouse.isPending
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
