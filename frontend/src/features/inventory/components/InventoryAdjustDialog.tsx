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
import { useAdjustInventory } from "@/features/inventory/hooks"
import { useProducts } from "@/features/products/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { useLocations } from "@/features/locations/hooks"

interface InventoryAdjustDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const InventoryAdjustDialog: React.FC<InventoryAdjustDialogProps> = ({
  open,
  onOpenChange,
}) => {
  const adjustInventory = useAdjustInventory()
  const { data: productsData } = useProducts({ limit: 100 })
  const { data: warehousesData } = useWarehouses({ limit: 100 })
  const { data: locationsData } = useLocations({ limit: 100 })

  const [formData, setFormData] = React.useState({
    productId: "",
    warehouseId: "",
    locationId: "",
    batchId: "",
    quantity: "",
    reason: "adjustment",
    notes: "",
  })

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await adjustInventory.mutateAsync({
        productId: parseInt(formData.productId),
        warehouseId: parseInt(formData.warehouseId),
        locationId: parseInt(formData.locationId),
        batchId: formData.batchId ? parseInt(formData.batchId) : undefined,
        quantity: parseInt(formData.quantity),
        reason: formData.reason,
        notes: formData.notes,
      })
      onOpenChange(false)
      setFormData({
        productId: "",
        warehouseId: "",
        locationId: "",
        batchId: "",
        quantity: "",
        reason: "adjustment",
        notes: "",
      })
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>Adjust Inventory</DialogTitle>
          <DialogDescription>
            Manually adjust stock levels for a product at a specific location.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Product *</label>
            <Select
              name="productId"
              value={formData.productId}
              onChange={handleChange}
              required
            >
              <option value="">Select a product</option>
              {productsData?.data?.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name} ({p.sku})
                </option>
              ))}
            </Select>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Warehouse *</label>
              <Select
                name="warehouseId"
                value={formData.warehouseId}
                onChange={handleChange}
                required
              >
                <option value="">Select warehouse</option>
                {warehousesData?.data?.map((w) => (
                  <option key={w.id} value={w.id}>
                    {w.name}
                  </option>
                ))}
              </Select>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Location *</label>
              <Select
                name="locationId"
                value={formData.locationId}
                onChange={handleChange}
                required
              >
                <option value="">Select location</option>
                {locationsData?.data
                  .filter((l) => !formData.warehouseId || l.warehouseId === parseInt(formData.warehouseId))
                  .map((l) => (
                    <option key={l.id} value={l.id}>
                      {l.code}
                    </option>
                  ))}
              </Select>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Quantity *</label>
              <Input
                name="quantity"
                type="number"
                value={formData.quantity}
                onChange={handleChange}
                placeholder="e.g. 10 or -5"
                required
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Reason *</label>
              <Select
                name="reason"
                value={formData.reason}
                onChange={handleChange}
                required
              >
                <option value="adjustment">Manual Adjustment</option>
                <option value="recount">Cycle Count / Recount</option>
                <option value="damage">Damage</option>
                <option value="return">Return</option>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Notes</label>
            <Input
              name="notes"
              value={formData.notes}
              onChange={handleChange}
              placeholder="Reason for adjustment"
            />
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={adjustInventory.isPending}>
              {adjustInventory.isPending ? "Adjusting..." : "Adjust Inventory"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
