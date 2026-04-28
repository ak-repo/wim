import * as React from "react"
import { ArrowRight, Package, MapPin, CheckCircle } from "lucide-react"
import { Button } from "@/components/ui/Button"
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/Dialog"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/Table"
import { Badge } from "@/components/ui/Badge"
import { useLocations } from "@/features/locations/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { usePutAwayPurchaseOrder } from "@/features/purchaseOrders/hooks"
import { Input } from "@/components/ui/Input"
import { Label } from "@/components/ui/Label"
import type { PurchaseOrderResponse } from "@/types/purchaseOrder"
import { DialogFooter } from "@/components/ui/Dialog"
import { Skeleton } from "@/components/ui/skeleton"

interface PutAwayDialogProps {
  order: PurchaseOrderResponse | null
  open: boolean
  onClose: VoidFunction
  onSuccess: VoidFunction
}

export function PutAwayDialog({ order, open, onClose, onSuccess }: PutAwayDialogProps) {
  const [formData, setFormData] = React.useState<any>({ notes: "", items: [] })
  const [selectedWarehouse, setSelectedWarehouse] = React.useState<number | null>(null)

  const { data: locationsData } = useLocations({ warehouseId: selectedWarehouse, page: 1, limit: 100 })
  const { data: warehousesData } = useWarehouses({ page: 1, limit: 50 })
  const putAwayMutation = usePutAwayPurchaseOrder()

  React.useEffect(() => {
    if (order && warehousesData?.data) {
      setSelectedWarehouse(order.warehouseId)
    }
  }, [order, warehousesData])

  React.useEffect(() => {
    if (order?.items) {
      setFormData({
        notes: `Put-away for ${order.refCode}`,
        items: order.items.map((item) => ({
          purchaseOrderItemId: item.id,
          quantity: item.quantityReceived - (item.quantityPutAway || 0),
          fromLocationId: "",
          toLocationId: "",
          batchId: item.batchNumber ? 1 : undefined, // Simplified
        }))
      })
    }
  }, [order])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!order) return

    const cleanedItems = formData.items
      .filter((item: any) => item.quantity > 0 && item.fromLocationId && item.toLocationId)
      .map((item: any) => ({
        purchaseOrderItemId: item.purchaseOrderItemId,
        quantity: item.quantity,
        fromLocationId: parseInt(item.fromLocationId, 10),
        toLocationId: parseInt(item.toLocationId, 10),
        batchId: item.batchId,
      }))

    putAwayMutation.mutate(
      { id: order.id.toString(), data: { notes: formData.notes, items: cleanedItems } },
      {
        onSuccess: () => {
          onSuccess()
          onClose()
        },
      }
    )
  }

  if (!order) return null

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Package className="h-5 w-5" />
            Put-Away - {order.refCode}
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-20">Quantity</TableHead>
                  <TableHead className="w-48">From Location</TableHead>
                  <TableHead className="w-48">To Location</TableHead>
                  <TableHead>Batch</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {order.items?.map((item: any, idx: number) => {
                  const maxQty = item.quantityReceived - (item.quantityPutAway || 0)
                  return (
                    <TableRow key={item.id}>
                      <TableCell className="font-medium">Product #{item.productId}</TableCell>
                      <TableCell>
                        <Badge variant={maxQty > 0 ? "outline" : "success"}>
                          {maxQty > 0 ? "Pending Put-Away" : "Completed"}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <Input
                          type="number"
                          min={0}
                          max={maxQty}
                          value={formData.items[idx]?.quantity || 0}
                          onChange={(e) => {
                            const newItems = [...formData.items]
                            newItems[idx] = { ...newItems[idx], quantity: parseInt(e.target.value, 10) || 0 }
                            setFormData({ ...formData, items: newItems })
                          }}
                          disabled={maxQty <= 0}
                        />
                      </TableCell>
                      <TableCell>
                        <LocationSelect
                          locations={locationsData?.data || []}
                          value={formData.items[idx]?.fromLocationId || ""}
                          onChange={(value) => {
                            const newItems = [...formData.items]
                            newItems[idx] = { ...newItems[idx], fromLocationId: value }
                            setFormData({ ...formData, items: newItems })
                          }}
                          disabled={maxQty <= 0}
                          placeholder="Receiving..."
                        />
                      </TableCell>
                      <TableCell>
                        <LocationSelect
                          locations={locationsData?.data || []}
                          value={formData.items[idx]?.toLocationId || ""}
                          onChange={(value) => {
                            const newItems = [...formData.items]
                            newItems[idx] = { ...newItems[idx], toLocationId: value }
                            setFormData({ ...formData, items: newItems })
                          }}
                          disabled={maxQty <= 0}
                          placeholder="Storage..."
                        />
                      </TableCell>
                      <TableCell>
                        {item.batchNumber || "-"}
                      </TableCell>
                    </TableRow>
                  )
                })}
              </TableBody>
            </Table>
          </div>

          <div className="grid grid-cols-1 gap-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="notes" className="text-right">Notes</Label>
              <Input
                id="notes"
                value={formData.notes}
                onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
                className="col-span-3"
                placeholder="Optional notes for this put-away..."
              />
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={onClose}>Cancel</Button>
            <Button
              type="submit"
              disabled={putAwayMutation.isPending || formData.items.filter((item: any) => item.quantity > 0).length === 0}
            >
              {putAwayMutation.isPending ? "Processing..." : "Complete Put-Away"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}

function LocationSelect({ locations, value, onChange, disabled, placeholder }: any) {
  if (disabled) {
    return <Input value={value} disabled placeholder="-" />
  }

  return (
    <select
      value={value}
      onChange={(e) => onChange(e.target.value)}
      disabled={disabled}
      className="flex h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
    >
      <option value="">{placeholder}</option>
      {locations.map((location: any) => (
        <option key={location.id} value={location.id}>
          {location.locationCode} - {location.zone || "Zone"}/{location.aisle || "Aisle"}
        </option>
      ))}
    </select>
  )
}