import * as React from "react"
import { Button } from "@/components/ui/Button"
import { Badge } from "@/components/ui/Badge"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/Dialog"
import {
  useAllocateSalesOrder,
  useDeallocateSalesOrder,
  useShipSalesOrder,
  useCancelSalesOrder,
} from "@/features/salesOrders/hooks"
import type { SalesOrder } from "@/features/salesOrders/types"
import { Loader2, PackageCheck, Truck, XCircle, Undo } from "lucide-react"

interface SalesOrderActionsDialogProps {
  action: "allocate" | "deallocate" | "ship" | "cancel"
  salesOrder: SalesOrder | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function SalesOrderActionsDialog({
  action,
  salesOrder,
  open,
  onOpenChange,
}: SalesOrderActionsDialogProps) {
  const [notes, setNotes] = React.useState("")

  const allocateMutation = useAllocateSalesOrder()
  const deallocateMutation = useDeallocateSalesOrder()
  const shipMutation = useShipSalesOrder()
  const cancelMutation = useCancelSalesOrder()

  const handleSubmit = async () => {
    if (!salesOrder) return

    try {
      switch (action) {
        case "allocate":
          await allocateMutation.mutateAsync(salesOrder.id)
          break
        case "deallocate":
          await deallocateMutation.mutateAsync(salesOrder.id)
          break
        case "ship":
          await shipMutation.mutateAsync({ id: salesOrder.id, data: { notes } })
          break
        case "cancel":
          await cancelMutation.mutateAsync(salesOrder.id)
          break
      }
      onOpenChange(false)
    } catch (error) {
      console.error(`Failed to ${action} sales order:`, error)
    }
  }

  const isPending =
    allocateMutation.isPending ||
    deallocateMutation.isPending ||
    shipMutation.isPending ||
    cancelMutation.isPending

  const getActionDetails = () => {
    switch (action) {
      case "allocate":
        return {
          title: "Allocate Stock",
          description: "Allocate inventory for this sales order. This will reserve stock in the warehouse.",
          icon: PackageCheck,
          buttonLabel: "Allocate",
        }
      case "deallocate":
        return {
          title: "Deallocate Stock",
          description: "Remove stock allocation from this sales order. Reserved stock will be released.",
          icon: Undo,
          buttonLabel: "Deallocate",
        }
      case "ship":
        return {
          title: "Ship Order",
          description: "Ship this sales order. Stock will be deducted from inventory.",
          icon: Truck,
          buttonLabel: "Ship",
        }
      case "cancel":
        return {
          title: "Cancel Order",
          description: "Cancel this sales order. If allocated, stock will be released first.",
          icon: XCircle,
          buttonLabel: "Cancel",
        }
      default:
        return {
          title: "",
          description: "",
          icon: PackageCheck,
          buttonLabel: "",
        }
    }
  }

  const details = getActionDetails()
  const Icon = details.icon

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "PENDING":
        return <Badge variant="secondary">Pending</Badge>
      case "PROCESSING":
        return <Badge variant="outline">Processing</Badge>
      case "SHIPPED":
        return <Badge variant="success">Shipped</Badge>
      case "CANCELLED":
        return <Badge variant="destructive">Cancelled</Badge>
      default:
        return <Badge>{status}</Badge>
    }
  }

  const getAllocationStatusBadge = (status: string) => {
    switch (status) {
      case "UNALLOCATED":
        return <Badge variant="secondary">Unallocated</Badge>
      case "PARTIALLY_ALLOCATED":
        return <Badge variant="outline">Partially Allocated</Badge>
      case "FULLY_ALLOCATED":
        return <Badge variant="success">Fully Allocated</Badge>
      default:
        return <Badge>{status}</Badge>
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Icon className="h-5 w-5" />
            {details.title}
          </DialogTitle>
          <DialogDescription>{details.description}</DialogDescription>
        </DialogHeader>

        {salesOrder && (
          <div className="space-y-4">
            <div className="p-4 border rounded-md">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-muted-foreground">Order Ref</p>
                  <p className="font-medium">{salesOrder.refCode}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Customer ID</p>
                  <p className="font-medium">{salesOrder.customerId}</p>
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Order Status</p>
                  {getStatusBadge(salesOrder.status)}
                </div>
                <div>
                  <p className="text-xs text-muted-foreground">Allocation</p>
                  {getAllocationStatusBadge(salesOrder.allocationStatus)}
                </div>
              </div>
            </div>

            {salesOrder.items && salesOrder.items.length > 0 && (
              <div className="border rounded-md">
                <table className="w-full">
                  <thead>
                    <tr className="border-b">
                      <th className="text-left p-2 text-xs">Product</th>
                      <th className="text-right p-2 text-xs">Qty Ordered</th>
                      <th className="text-right p-2 text-xs">Qty Reserved</th>
                      <th className="text-right p-2 text-xs">Qty Shipped</th>
                    </tr>
                  </thead>
                  <tbody>
                    {salesOrder.items.map((item) => (
                      <tr key={item.id} className="border-b">
                        <td className="p-2">Product {item.productId}</td>
                        <td className="p-2 text-right">{item.quantityOrdered}</td>
                        <td className="p-2 text-right">{item.quantityReserved}</td>
                        <td className="p-2 text-right">{item.quantityShipped}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            {action === "ship" && (
              <div>
                <label className="text-sm font-medium">Shipping Notes</label>
                <textarea
                  className="w-full h-20 px-3 py-2 rounded-md border border-input bg-background mt-1"
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  placeholder="Add shipping notes (optional)"
                />
              </div>
            )}
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} disabled={isPending}>
            {isPending && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
            {details.buttonLabel}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}