import * as React from "react"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/Dialog"
import { useCreateSalesOrder, useUpdateSalesOrder } from "@/features/salesOrders/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { useProducts } from "@/features/products/hooks"
import type { SalesOrder, CreateSalesOrderRequest } from "@/features/salesOrders/types"
import { Loader2, Plus, X } from "lucide-react"

interface SalesOrderFormDialogProps {
  salesOrder?: SalesOrder | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function SalesOrderFormDialog({ salesOrder, open, onOpenChange }: SalesOrderFormDialogProps) {
  const [formData, setFormData] = React.useState<CreateSalesOrderRequest>({
    customerId: 0,
    warehouseId: 0,
    requiredDate: undefined,
    shippingMethod: "",
    shippingAddress: "",
    billingAddress: "",
    notes: "",
    items: [],
  })
  const [newItem, setNewItem] = React.useState({ productId: 0, quantityOrdered: 1 })
  const [showAddItem, setShowAddItem] = React.useState(false)

  const createMutation = useCreateSalesOrder()
  const updateMutation = useUpdateSalesOrder()

  const { data: warehousesData } = useWarehouses({ page: 1, limit: 100 })
  const { data: productsData } = useProducts({ page: 1, limit: 100 })

  React.useEffect(() => {
    if (salesOrder) {
      setFormData({
        customerId: salesOrder.customerId,
        warehouseId: salesOrder.warehouseId,
        requiredDate: salesOrder.requiredDate,
        shippingMethod: salesOrder.shippingMethod || "",
        shippingAddress: salesOrder.shippingAddress || "",
        billingAddress: salesOrder.billingAddress || "",
        notes: salesOrder.notes || "",
        items: salesOrder.items?.map(item => ({
          productId: item.productId,
          quantityOrdered: item.quantityOrdered,
          unitPrice: item.unitPrice,
        })) || [],
      })
    } else {
      setFormData({
        customerId: 0,
        warehouseId: 0,
        requiredDate: undefined,
        shippingMethod: "",
        shippingAddress: "",
        billingAddress: "",
        notes: "",
        items: [],
      })
    }
  }, [salesOrder, open])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (salesOrder) {
        await updateMutation.mutateAsync({ id: salesOrder.id, data: formData })
      } else {
        await createMutation.mutateAsync(formData)
      }
      onOpenChange(false)
    } catch (error) {
      console.error("Failed to save sales order:", error)
    }
  }

  const handleAddItem = () => {
    if (newItem.productId > 0 && newItem.quantityOrdered > 0) {
      setFormData(prev => ({
        ...prev,
        items: [...prev.items, { ...newItem }],
      }))
      setNewItem({ productId: 0, quantityOrdered: 1 })
      setShowAddItem(false)
    }
  }

  const handleRemoveItem = (index: number) => {
    setFormData(prev => ({
      ...prev,
      items: prev.items.filter((_, i) => i !== index),
    }))
  }

  const isPending = createMutation.isPending || updateMutation.isPending
  const warehouses = warehousesData?.data || []
  const products = productsData?.data || []

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{salesOrder ? "Edit Sales Order" : "Create Sales Order"}</DialogTitle>
          <DialogDescription>
            {salesOrder ? "Update the sales order details." : "Create a new sales order for a customer."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium">Customer ID</label>
              <Input
                type="number"
                value={formData.customerId}
                onChange={(e) => setFormData(prev => ({ ...prev, customerId: parseInt(e.target.value) || 0 }))}
                placeholder="Customer ID"
                required
              />
            </div>
            <div>
              <label className="text-sm font-medium">Warehouse</label>
              <select
                className="w-full h-10 px-3 rounded-md border border-input bg-background"
                value={formData.warehouseId}
                onChange={(e) => setFormData(prev => ({ ...prev, warehouseId: parseInt(e.target.value) }))}
                required
              >
                <option value={0}>Select warehouse</option>
                {warehouses.map((w) => (
                  <option key={w.id} value={w.id}>{w.name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-sm font-medium">Required Date</label>
              <Input
                type="date"
                value={formData.requiredDate?.split("T")[0] || ""}
                onChange={(e) => setFormData(prev => ({ ...prev, requiredDate: e.target.value }))}
              />
            </div>
            <div>
              <label className="text-sm font-medium">Shipping Method</label>
              <Input
                value={formData.shippingMethod || ""}
                onChange={(e) => setFormData(prev => ({ ...prev, shippingMethod: e.target.value }))}
                placeholder="e.g., Standard, Express"
              />
            </div>
          </div>

          <div>
            <label className="text-sm font-medium">Shipping Address</label>
            <textarea
              className="w-full h-20 px-3 py-2 rounded-md border border-input bg-background"
              value={formData.shippingAddress || ""}
              onChange={(e) => setFormData(prev => ({ ...prev, shippingAddress: e.target.value }))}
              placeholder="Enter shipping address"
            />
          </div>

          <div>
            <label className="text-sm font-medium">Billing Address</label>
            <textarea
              className="w-full h-20 px-3 py-2 rounded-md border border-input bg-background"
              value={formData.billingAddress || ""}
              onChange={(e) => setFormData(prev => ({ ...prev, billingAddress: e.target.value }))}
              placeholder="Enter billing address"
            />
          </div>

          <div>
            <label className="text-sm font-medium">Notes</label>
            <textarea
              className="w-full h-16 px-3 py-2 rounded-md border border-input bg-background"
              value={formData.notes || ""}
              onChange={(e) => setFormData(prev => ({ ...prev, notes: e.target.value }))}
              placeholder="Order notes"
            />
          </div>

          <div className="border rounded-md p-4">
            <div className="flex items-center justify-between mb-4">
              <h4 className="font-medium">Order Items</h4>
              <Button type="button" variant="outline" size="sm" onClick={() => setShowAddItem(true)}>
                <Plus className="h-4 w-4 mr-1" /> Add Item
              </Button>
            </div>

            {formData.items.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-4">No items added yet.</p>
            ) : (
              <div className="space-y-2">
                {formData.items.map((item, index) => {
                  const product = products.find(p => p.id === item.productId)
                  return (
                    <div key={index} className="flex items-center justify-between p-2 bg-muted rounded">
                      <div>
                        <span className="font-medium">{product?.name || `Product ${item.productId}`}</span>
                        <span className="text-muted-foreground ml-2">x{item.quantityOrdered}</span>
                      </div>
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        onClick={() => handleRemoveItem(index)}
                      >
                        <X className="h-4 w-4" />
                      </Button>
                    </div>
                  )
                })}
              </div>
            )}

            {showAddItem && (
              <div className="mt-4 p-4 border rounded-md bg-muted/50">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-sm font-medium">Product</label>
                    <select
                      className="w-full h-10 px-3 rounded-md border border-input bg-background"
                      value={newItem.productId}
                      onChange={(e) => setNewItem(prev => ({ ...prev, productId: parseInt(e.target.value) }))}
                    >
                      <option value={0}>Select product</option>
                      {products.map((p) => (
                        <option key={p.id} value={p.id}>{p.name} ({p.sku})</option>
                      ))}
                    </select>
                  </div>
                  <div>
                    <label className="text-sm font-medium">Quantity</label>
                    <Input
                      type="number"
                      min="1"
                      value={newItem.quantityOrdered}
                      onChange={(e) => setNewItem(prev => ({ ...prev, quantityOrdered: parseInt(e.target.value) || 1 }))}
                    />
                  </div>
                </div>
                <div className="flex gap-2 mt-4">
                  <Button type="button" onClick={handleAddItem}>Add</Button>
                  <Button type="button" variant="outline" onClick={() => setShowAddItem(false)}>Cancel</Button>
                </div>
              </div>
            )}
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isPending || formData.items.length === 0}>
              {isPending && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
              {salesOrder ? "Update" : "Create"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}