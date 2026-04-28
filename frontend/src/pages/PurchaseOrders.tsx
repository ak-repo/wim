import * as React from "react"
import { Plus, Search, Clock, Package, AlertCircle, XCircle, ArrowRightToLine } from "lucide-react"

import { usePurchaseOrders } from "@/features/purchaseOrders/hooks"
import { Button } from "@/components/ui/Button"
import { Badge } from "@/components/ui/Badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Input } from "@/components/ui/Input"
import { PutAwayDialog } from "@/features/purchaseOrders/components/PutAwayDialog"

export default function PurchaseOrdersPage() {
  const [search, setSearch] = React.useState("")
  const [selectedOrder, setSelectedOrder] = React.useState<any>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [showPutAway, setShowPutAway] = React.useState(false)

  const { data: ordersData, isLoading } = usePurchaseOrders({
    page: 1,
    limit: 10,
  })

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "PENDING":
        return <Clock className="h-4 w-4 text-muted-foreground" />
      case "PARTIAL_RECEIVED":
        return <AlertCircle className="h-4 w-4 text-yellow-500" />
      case "RECEIVED":
        return <Package className="h-4 w-4 text-blue-500" />
      case "CANCELLED":
        return <XCircle className="h-4 w-4 text-red-500" />
      default:
        return <Package className="h-4 w-4 text-muted-foreground" />
    }
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "PENDING":
        return <Badge variant="secondary">Pending</Badge>
      case "PARTIAL_RECEIVED":
        return <Badge variant="outline">Partial Received</Badge>
      case "RECEIVED":
        return <Badge variant="success">Received</Badge>
      case "CANCELLED":
        return <Badge variant="destructive">Cancelled</Badge>
      default:
        return <Badge>{status}</Badge>
    }
  }

  const filteredOrders = ordersData?.data?.filter(
    (order) =>
      order.refCode.toLowerCase().includes(search.toLowerCase()) ||
      order.status.toLowerCase().includes(search.toLowerCase())
  ) || []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">
            Purchase Orders
          </h1>
          <p className="text-muted-foreground mt-1">
            Manage supplier purchase orders and receiving.
          </p>
        </div>
        <Button onClick={() => setIsFormOpen(true)} size="lg">
          <Plus className="h-4 w-4 mr-2" />
          New Order
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between mb-4">
            <CardTitle className="text-lg">All Purchase Orders</CardTitle>
            <div className="relative w-64">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by ID or status..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-9"
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-2">
              <Skeleton className="h-12 w-full" />
              <Skeleton className="h-12 w-full" />
              <Skeleton className="h-12 w-full" />
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Supplier</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Items</TableHead>
                  <TableHead>Expected Date</TableHead>
                  <TableHead>Warehouse</TableHead>
                  <TableHead>Notes</TableHead>
                  <TableHead className="w-12">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredOrders.map((order) => (
                  <TableRow key={order.id}>
                    <TableCell className="font-medium">
                      <div className="flex items-center gap-2">
                        {getStatusIcon(order.status)}
                        {order.refCode}
                      </div>
                    </TableCell>
                    <TableCell>Supplier #{order.supplierId}</TableCell>
                    <TableCell>{getStatusBadge(order.status)}</TableCell>
                    <TableCell>
                      {order.items?.length || 0} item
                      {order.items && order.items.length !== 1 ? "s" : ""}
                    </TableCell>
                    <TableCell>
                      {order.expectedDate
                        ? new Date(order.expectedDate).toLocaleDateString()
                        : "-"}
                    </TableCell>
                    <TableCell>Warehouse #{order.warehouseId}</TableCell>
                    <TableCell className="max-w-xs truncate">
                      {order.notes || "-"}
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setSelectedOrder(order)}
                      >
                        View
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {selectedOrder && (
        <>
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Order Details</CardTitle>
                {selectedOrder.status === "RECEIVED" && (
                  <Button onClick={() => setShowPutAway(true)} size="sm">
                    <ArrowRightToLine className="h-4 w-4 mr-2" />
                    Put Away
                  </Button>
                )}
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-2 text-sm">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <span className="font-medium">ID:</span> {selectedOrder.refCode}
                  </div>
                  <div>
                    <span className="font-medium">Status:</span> {getStatusBadge(selectedOrder.status)}
                  </div>
                  <div>
                    <span className="font-medium">Supplier:</span> #{selectedOrder.supplierId}
                  </div>
                  <div>
                    <span className="font-medium">Warehouse:</span> #{selectedOrder.warehouseId}
                  </div>
                  <div>
                    <span className="font-medium">Expected Date:</span> {selectedOrder.expectedDate ? new Date(selectedOrder.expectedDate).toLocaleDateString() : "-"}
                  </div>
                  <div>
                    <span className="font-medium">Notes:</span> {selectedOrder.notes || "-"}
                  </div>
                </div>
                <div className="mt-4">
                  <h4 className="font-medium mb-2">Items:</h4>
                  {selectedOrder.items && selectedOrder.items.length > 0 ? (
                    <div className="space-y-1">
                      {selectedOrder.items.map((item: any) => (
                        <div key={item.id} className="flex justify-between py-1 border-b">
                          <span>Product #{item.productId}</span>
                          <span>{item.quantityReceived}/{item.quantityOrdered} received</span>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <p className="text-muted-foreground">No items</p>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          <PutAwayDialog
            order={selectedOrder}
            open={showPutAway}
            onClose={() => setShowPutAway(false)}
            onSuccess={() => {
              setSelectedOrder(null)
            }}
          />
        </>
      )}
    </div>
  )
}
