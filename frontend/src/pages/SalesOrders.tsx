import * as React from "react"
import { useSalesOrders } from "@/features/salesOrders/hooks"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Badge } from "@/components/ui/Badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Card, CardContent, CardHeader } from "@/components/ui/Card"
import { SalesOrderFormDialog } from "@/features/salesOrders/components/SalesOrderFormDialog"
import { SalesOrderActionsDialog } from "@/features/salesOrders/components/SalesOrderActionsDialog"
import type { SalesOrder } from "@/features/salesOrders/types"
import { Plus, Search, Pencil, Loader2, Package, Truck, PackageCheck, XCircle } from "lucide-react"

export default function SalesOrdersPage() {
  const [search, setSearch] = React.useState("")
  const [selectedOrder, setSelectedOrder] = React.useState<SalesOrder | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [actionDialog, setActionDialog] = React.useState<{
    action: "allocate" | "deallocate" | "ship" | "cancel"
    open: boolean
  }>({ action: "allocate", open: false })

  const { data: ordersData, isLoading } = useSalesOrders({
    page: 1,
    limit: 10,
  })

  const handleEdit = (order: SalesOrder) => {
    setSelectedOrder(order)
    setIsFormOpen(true)
  }

  const handleCreate = () => {
    setSelectedOrder(null)
    setIsFormOpen(true)
  }

  const handleAction = (order: SalesOrder, action: "allocate" | "deallocate" | "ship" | "cancel") => {
    setSelectedOrder(order)
    setActionDialog({ action, open: true })
  }

  const filteredOrders =
    ordersData?.data?.filter(
      (order) =>
        order.refCode.toLowerCase().includes(search.toLowerCase()) ||
        order.status.toLowerCase().includes(search.toLowerCase())
    ) || []

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
        return <Badge variant="outline">Partially</Badge>
      case "FULLY_ALLOCATED":
        return <Badge variant="success">Allocated</Badge>
      default:
        return <Badge>{status}</Badge>
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Sales Orders</h2>
          <p className="text-muted-foreground">Manage customer sales orders.</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="h-4 w-4 mr-2" />
          New Order
        </Button>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center gap-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search orders..."
                className="pl-9"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Ref Code</TableHead>
                  <TableHead>Customer</TableHead>
                  <TableHead>Warehouse</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Allocation</TableHead>
                  <TableHead>Order Date</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-24 text-center">
                      <Loader2 className="h-6 w-6 animate-spin mx-auto text-muted-foreground" />
                    </TableCell>
                  </TableRow>
                ) : filteredOrders.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-24 text-center">
                      <div className="flex flex-col items-center justify-center text-muted-foreground">
                        <Package className="h-8 w-8 mb-2" />
                        <p>No sales orders found.</p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredOrders.map((order) => (
                    <TableRow key={order.id}>
                      <TableCell className="font-medium">{order.refCode}</TableCell>
                      <TableCell>Customer #{order.customerId}</TableCell>
                      <TableCell>Warehouse #{order.warehouseId}</TableCell>
                      <TableCell>{getStatusBadge(order.status)}</TableCell>
                      <TableCell>{getAllocationStatusBadge(order.allocationStatus)}</TableCell>
                      <TableCell>
                        {new Date(order.orderDate).toLocaleDateString()}
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-1">
                          {order.status === "PENDING" && (
                            <>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleAction(order, "allocate")}
                                title="Allocate"
                              >
                                <PackageCheck className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleAction(order, "cancel")}
                                title="Cancel"
                                className="text-destructive hover:text-destructive"
                              >
                                <XCircle className="h-4 w-4" />
                              </Button>
                            </>
                          )}
                          {(order.allocationStatus === "PARTIALLY_ALLOCATED" || order.allocationStatus === "FULLY_ALLOCATED") && order.status !== "SHIPPED" && (
                            <>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleAction(order, "ship")}
                                title="Ship"
                              >
                                <Truck className="h-4 w-4" />
                              </Button>
                              {order.allocationStatus === "FULLY_ALLOCATED" && order.status === "PROCESSING" && (
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  onClick={() => handleAction(order, "deallocate")}
                                  title="Deallocate"
                                >
                                  <PackageCheck className="h-4 w-4 rotate-180" />
                                </Button>
                              )}
                            </>
                          )}
                          {order.status === "PROCESSING" && (
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleEdit(order)}
                              title="Edit"
                            >
                              <Pencil className="h-4 w-4" />
                            </Button>
                          )}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <SalesOrderFormDialog
        salesOrder={selectedOrder}
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
      />

      <SalesOrderActionsDialog
        action={actionDialog.action}
        salesOrder={selectedOrder}
        open={actionDialog.open}
        onOpenChange={(open) => setActionDialog({ ...actionDialog, open })}
      />
    </div>
  )
}