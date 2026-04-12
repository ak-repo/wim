import * as React from "react"
import { useSalesOrders } from "@/features/salesOrders/hooks"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Badge } from "@/components/ui/Badge"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { SalesOrderFormDialog } from "@/features/salesOrders/components/SalesOrderFormDialog"
import { SalesOrderActionsDialog } from "@/features/salesOrders/components/SalesOrderActionsDialog"
import type { SalesOrder } from "@/features/salesOrders/types"
import { Plus, Search, Pencil, Package, Truck, PackageCheck, XCircle, ShoppingCart } from "lucide-react"

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
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Sales Orders</h1>
          <p className="text-muted-foreground mt-1">Manage customer sales orders and fulfillment.</p>
        </div>
        <Button onClick={handleCreate} size="lg">
          <Plus className="h-4 w-4 mr-2" />
          New Order
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Order Management</CardTitle>
              <CardDescription>
                Search and manage sales orders
              </CardDescription>
            </div>
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              <ShoppingCart className="h-4 w-4 text-primary" />
            </div>
          </div>
          <div className="flex items-center gap-4 mt-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by ref code or status..."
                className="pl-9 bg-muted/50 border-transparent focus:bg-background"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
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
                Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-28" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-8 w-24 ml-auto" /></TableCell>
                  </TableRow>
                ))
              ) : filteredOrders.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="h-40">
                    <div className="flex flex-col items-center justify-center text-center">
                      <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                        <Package className="h-6 w-6 text-muted-foreground" />
                      </div>
                      <p className="text-sm font-medium text-foreground">No sales orders found</p>
                      <p className="text-xs text-muted-foreground mt-1">
                        {search ? "Try adjusting your search" : "Create an order to get started"}
                      </p>
                      {!search && (
                        <Button variant="outline" size="sm" className="mt-3" onClick={handleCreate}>
                          <Plus className="h-3 w-3 mr-1" />
                          New Order
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                filteredOrders.map((order) => (
                  <TableRow key={order.id}>
                    <TableCell>
                      <span className="font-mono text-xs font-medium">{order.refCode}</span>
                    </TableCell>
                    <TableCell className="font-medium text-foreground">
                      Customer #{order.customerId}
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      Warehouse #{order.warehouseId}
                    </TableCell>
                    <TableCell>{getStatusBadge(order.status)}</TableCell>
                    <TableCell>{getAllocationStatusBadge(order.allocationStatus)}</TableCell>
                    <TableCell className="text-muted-foreground text-sm">
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
                              className="h-8 w-8"
                            >
                              <PackageCheck className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => handleAction(order, "cancel")}
                              title="Cancel"
                              className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
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
                              className="h-8 w-8"
                            >
                              <Truck className="h-4 w-4" />
                            </Button>
                            {order.allocationStatus === "FULLY_ALLOCATED" && order.status === "PROCESSING" && (
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleAction(order, "deallocate")}
                                title="Deallocate"
                                className="h-8 w-8"
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
                            className="h-8 w-8"
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