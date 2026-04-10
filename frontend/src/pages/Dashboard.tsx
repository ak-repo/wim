import * as React from "react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import {
  Package,
  Users,
  Warehouse,
  MapPin,
  Boxes,
  Database,
  History,
} from "lucide-react"
import { useProducts } from "@/features/products/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { useLocations } from "@/features/locations/hooks"
import { useInventoryList, useStockMovements } from "@/features/inventory/hooks"

export default function DashboardPage() {
  const { data: productsData } = useProducts({ limit: 1 })
  const { data: warehousesData } = useWarehouses({ limit: 1 })
  const { data: locationsData } = useLocations({ limit: 1 })
  const { data: inventoryData } = useInventoryList({ limit: 1000 })
  const { data: movementsData } = useStockMovements({ limit: 5 })

  const totalStock = inventoryData?.data?.reduce((acc, curr) => acc + curr.quantity, 0) || 0

  const stats = [
    {
      title: "Total Products",
      value: productsData?.total_count || "0",
      description: "Active products in catalog",
      icon: Package,
    },
    {
      title: "Total Stock",
      value: totalStock.toString(),
      description: "Units across all locations",
      icon: Database,
    },
    {
      title: "Warehouses",
      value: warehousesData?.total_count || "0",
      description: "Active warehouses",
      icon: Warehouse,
    },
    {
      title: "Locations",
      value: locationsData?.total_count || "0",
      description: "Storage locations",
      icon: MapPin,
    },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
        <p className="text-muted-foreground">
          Welcome to your warehouse management dashboard.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon
          return (
            <Card key={stat.title}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">{stat.title}</CardTitle>
                <Icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
                <p className="text-xs text-muted-foreground">
                  {stat.description}
                </p>
              </CardContent>
            </Card>
          )
        })}
      </div>

      {/* Quick Actions & Recent Activity */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>Common tasks and operations</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <a
              href="/inventory"
              className="flex items-center gap-2 rounded-md border p-3 hover:bg-muted transition-colors text-sm font-medium"
            >
              <Database className="h-4 w-4 text-primary" />
              <span>Adjust Inventory Stock</span>
            </a>
            <a
              href="/products"
              className="flex items-center gap-2 rounded-md border p-3 hover:bg-muted transition-colors text-sm font-medium"
            >
              <Package className="h-4 w-4 text-primary" />
              <span>Manage Product Catalog</span>
            </a>
            <a
              href="/locations"
              className="flex items-center gap-2 rounded-md border p-3 hover:bg-muted transition-colors text-sm font-medium"
            >
              <MapPin className="h-4 w-4 text-primary" />
              <span>Configure Storage Locations</span>
            </a>
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Recent Stock Movements</CardTitle>
            <CardDescription>Latest inventory transactions</CardDescription>
          </CardHeader>
          <CardContent>
            {movementsData?.data && movementsData.data.length > 0 ? (
              <div className="space-y-4">
                {movementsData.data.map((m) => (
                  <div key={m.id} className="flex items-center justify-between border-b pb-2 last:border-0">
                    <div className="flex items-center gap-3">
                      <div className={`p-2 rounded-full ${m.quantity > 0 ? 'bg-emerald-100 text-emerald-600' : 'bg-red-100 text-red-600'}`}>
                        <History className="h-4 w-4" />
                      </div>
                      <div>
                        <p className="text-sm font-medium">{m.movementType}</p>
                        <p className="text-xs text-muted-foreground">{new Date(m.createdAt).toLocaleString()}</p>
                      </div>
                    </div>
                    <div className={`text-sm font-bold ${m.quantity > 0 ? 'text-emerald-600' : 'text-red-600'}`}>
                      {m.quantity > 0 ? `+${m.quantity}` : m.quantity}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-8 text-center">
                <Boxes className="h-8 w-8 text-muted-foreground mb-2" />
                <p className="text-sm text-muted-foreground">No recent activity</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
