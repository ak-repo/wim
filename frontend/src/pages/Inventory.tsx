import * as React from "react"
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
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import { Badge } from "@/components/ui/Badge"
import { useProducts } from "@/features/products/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { useLocations } from "@/features/locations/hooks"
import { Database, History, Plus, Loader2 } from "lucide-react"

export default function InventoryPage() {
  const [isAdjustOpen, setIsAdjustOpen] = React.useState(false)
  const [activeTab, setActiveTab] = React.useState<"current" | "movements">("current")
  
  const { data: inventoryData, isLoading: invLoading } = useInventoryList({ limit: 50 })
  const { data: movementsData, isLoading: movLoading } = useStockMovements({ limit: 50 })
  
  const { data: productsData } = useProducts({ limit: 100 })
  const { data: warehousesData } = useWarehouses({ limit: 100 })
  const { data: locationsData } = useLocations({ limit: 100 })

  const getProductName = (id: number) => productsData?.data.find(p => p.id === id)?.name || id
  const getWarehouseName = (id: number) => warehousesData?.data.find(w => w.id === id)?.name || id
  const getLocationCode = (id: number) => locationsData?.data.find(l => l.id === id)?.code || id

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Inventory</h2>
          <p className="text-muted-foreground">
            Manage your stock levels and track movements.
          </p>
        </div>
        <Button onClick={() => setIsAdjustOpen(true)}>
          <Plus className="h-4 w-4 mr-2" />
          Adjust Inventory
        </Button>
      </div>

      <div className="flex space-x-4 border-b">
        <button
          className={`pb-2 px-1 text-sm font-medium transition-colors hover:text-primary flex items-center gap-2 ${
            activeTab === "current" ? "border-b-2 border-primary text-primary" : "text-muted-foreground"
          }`}
          onClick={() => setActiveTab("current")}
        >
          <Database className="h-4 w-4" />
          Current Stock
        </button>
        <button
          className={`pb-2 px-1 text-sm font-medium transition-colors hover:text-primary flex items-center gap-2 ${
            activeTab === "movements" ? "border-b-2 border-primary text-primary" : "text-muted-foreground"
          }`}
          onClick={() => setActiveTab("movements")}
        >
          <History className="h-4 w-4" />
          Stock Movements
        </button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>
            {activeTab === "current" ? "Current Stock Levels" : "Recent Stock Movements"}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
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
                    <TableRow>
                      <TableCell colSpan={7} className="h-24 text-center">
                        <Loader2 className="h-6 w-6 animate-spin mx-auto text-muted-foreground" />
                      </TableCell>
                    </TableRow>
                  ) : inventoryData?.data.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="h-24 text-center text-muted-foreground">
                        No inventory records found.
                      </TableCell>
                    </TableRow>
                  ) : (
                    inventoryData?.data.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell className="font-medium">{getProductName(item.productId)}</TableCell>
                        <TableCell>{getWarehouseName(item.warehouseId)}</TableCell>
                        <TableCell>{getLocationCode(item.locationId)}</TableCell>
                        <TableCell className="font-semibold">{item.quantity}</TableCell>
                        <TableCell className="text-muted-foreground">{item.reservedQty}</TableCell>
                        <TableCell>
                          <Badge variant={item.availableQty > 0 ? "success" : "destructive"}>
                            {item.availableQty}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right text-sm text-muted-foreground">
                          {new Date(item.updatedAt).toLocaleString()}
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
                    <TableHead>Qty</TableHead>
                    <TableHead>Notes</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {movLoading ? (
                    <TableRow>
                      <TableCell colSpan={6} className="h-24 text-center">
                        <Loader2 className="h-6 w-6 animate-spin mx-auto text-muted-foreground" />
                      </TableCell>
                    </TableRow>
                  ) : movementsData?.data.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                        No movements recorded.
                      </TableCell>
                    </TableRow>
                  ) : (
                    movementsData?.data.map((m) => (
                      <TableRow key={m.id}>
                        <TableCell className="text-sm whitespace-nowrap">
                          {new Date(m.createdAt).toLocaleString()}
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline">{m.movementType}</Badge>
                        </TableCell>
                        <TableCell className="font-medium">{getProductName(m.productId)}</TableCell>
                        <TableCell>{getWarehouseName(m.warehouseId)}</TableCell>
                        <TableCell>
                          <span className={m.quantity > 0 ? "text-green-600 font-semibold" : "text-red-600 font-semibold"}>
                            {m.quantity > 0 ? `+${m.quantity}` : m.quantity}
                          </span>
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground truncate max-w-xs">
                          {m.notes || "-"}
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            )}
          </div>
        </CardContent>
      </Card>

      <InventoryAdjustDialog 
        open={isAdjustOpen} 
        onOpenChange={setIsAdjustOpen} 
      />
    </div>
  )
}
