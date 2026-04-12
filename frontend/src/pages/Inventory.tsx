import * as React from "react"
import { cn } from "@/utils"
import { useInventoryList, useStockMovements } from "@/features/inventory/hooks"
import { InventoryAdjustDialog } from "@/features/inventory/components/InventoryAdjustDialog"
import { Button } from "@/components/ui/Button"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { Badge } from "@/components/ui/Badge"
import { Skeleton } from "@/components/ui/skeleton"
import { useProducts } from "@/features/products/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { useLocations } from "@/features/locations/hooks"
import { Database, History, Plus, ArrowUpRight, ArrowDownRight } from "lucide-react"

export default function InventoryPage() {
  const [isAdjustOpen, setIsAdjustOpen] = React.useState(false)
  const [activeTab, setActiveTab] = React.useState<"current" | "movements">("current")
  
  const { data: inventoryData, isLoading: invLoading } = useInventoryList({ page: 1, limit: 50 })
  const { data: movementsData, isLoading: movLoading } = useStockMovements({ page: 1, limit: 50 })
  
  const { data: productsData } = useProducts({ page: 1, limit: 100 })
  const { data: warehousesData } = useWarehouses({ page: 1, limit: 100 })
  const { data: locationsData } = useLocations({ page: 1, limit: 100 })

  const getProductName = (id: number) => productsData?.data.find(p => p.id === id)?.name || `ID: ${id}`
  const getWarehouseName = (id: string) => warehousesData?.data.find(w => w.id === id)?.name || `ID: ${id}`
  const getLocationCode = (id: string) => locationsData?.data.find(l => l.id === id)?.locationCode || `ID: ${id}`

  const getWarehouseNameById = (id: number) => getWarehouseName(String(id))
  const getLocationCodeById = (id: number) => getLocationCode(String(id))

  const inventoryItems = inventoryData?.data || []
  const movements = movementsData?.data || []

  const tabs = [
    { id: "current", label: "Current Stock", icon: Database },
    { id: "movements", label: "Stock Movements", icon: History },
  ] as const

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Inventory</h1>
          <p className="text-muted-foreground mt-1">
            Manage your stock levels and track movements.
          </p>
        </div>
        <Button onClick={() => setIsAdjustOpen(true)} size="lg">
          <Plus className="h-4 w-4 mr-2" />
          Adjust Inventory
        </Button>
      </div>

      <div className="flex items-center gap-2 rounded-lg bg-muted/50 p-1 w-fit">
        {tabs.map((tab) => {
          const Icon = tab.icon
          const isActive = activeTab === tab.id
          return (
            <button
              key={tab.id}
              className={cn(
                "flex items-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-all duration-200",
                isActive
                  ? "bg-background text-foreground shadow-sm"
                  : "text-muted-foreground hover:text-foreground"
              )}
              onClick={() => setActiveTab(tab.id)}
            >
              <Icon className="h-4 w-4" />
              {tab.label}
            </button>
          )
        })}
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>
                {activeTab === "current" ? "Current Stock Levels" : "Recent Stock Movements"}
              </CardTitle>
              <CardDescription>
                {activeTab === "current"
                  ? "Real-time inventory across all locations"
                  : "Latest inventory transactions and adjustments"}
              </CardDescription>
            </div>
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              {activeTab === "current" ? (
                <Database className="h-4 w-4 text-primary" />
              ) : (
                <History className="h-4 w-4 text-primary" />
              )}
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {activeTab === "current" ? (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product</TableHead>
                  <TableHead>Warehouse</TableHead>
                  <TableHead>Location</TableHead>
                  <TableHead>Quantity</TableHead>
                  <TableHead>Reserved</TableHead>
                  <TableHead>Available</TableHead>
                  <TableHead className="text-right">Last Updated</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {invLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-12" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-12" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-12" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-24 ml-auto" /></TableCell>
                    </TableRow>
                  ))
                ) : inventoryItems.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-32">
                      <div className="flex flex-col items-center justify-center text-center">
                        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                          <Database className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <p className="text-sm font-medium text-foreground">No inventory records</p>
                        <p className="text-xs text-muted-foreground mt-1">
                          Adjust inventory to get started
                        </p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  inventoryItems.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell className="font-medium">{getProductName(item.productId)}</TableCell>
                      <TableCell>{getWarehouseNameById(item.warehouseId)}</TableCell>
                      <TableCell>
                        <Badge variant="outline">{getLocationCodeById(item.locationId)}</Badge>
                      </TableCell>
                      <TableCell className="font-semibold">{item.quantity.toLocaleString()}</TableCell>
                      <TableCell className="text-muted-foreground">{item.reservedQty.toLocaleString()}</TableCell>
                      <TableCell>
                        <Badge variant={item.availableQty > 0 ? "success" : "secondary"}>
                          {item.availableQty.toLocaleString()}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right text-sm text-muted-foreground">
                        {new Date(item.updatedAt).toLocaleDateString()}
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Date</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Product</TableHead>
                  <TableHead>Warehouse</TableHead>
                  <TableHead>Quantity</TableHead>
                  <TableHead>Notes</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {movLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                      <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-40" /></TableCell>
                    </TableRow>
                  ))
                ) : movements.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-32">
                      <div className="flex flex-col items-center justify-center text-center">
                        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                          <History className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <p className="text-sm font-medium text-foreground">No movements recorded</p>
                        <p className="text-xs text-muted-foreground mt-1">
                          Stock movements will appear here
                        </p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  movements.map((m) => {
                    const isPositive = m.quantity > 0
                    return (
                      <TableRow key={m.id}>
                        <TableCell className="text-sm text-muted-foreground whitespace-nowrap">
                          {new Date(m.createdAt).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline" className="text-xs">{m.movementType}</Badge>
                        </TableCell>
                        <TableCell className="font-medium">{getProductName(m.productId)}</TableCell>
                        <TableCell>{getWarehouseNameById(m.warehouseId)}</TableCell>
                        <TableCell>
                          <div className={`flex items-center gap-1 ${isPositive ? "text-emerald-500" : "text-red-500"}`}>
                            {isPositive ? (
                              <ArrowUpRight className="h-3 w-3" />
                            ) : (
                              <ArrowDownRight className="h-3 w-3" />
                            )}
                            <span className="font-semibold">
                              {isPositive ? "+" : ""}{m.quantity.toLocaleString()}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground truncate max-w-xs">
                          {m.notes || "-"}
                        </TableCell>
                      </TableRow>
                    )
                  })
                )}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <InventoryAdjustDialog 
        open={isAdjustOpen} 
        onOpenChange={setIsAdjustOpen} 
      />
    </div>
  )
}
