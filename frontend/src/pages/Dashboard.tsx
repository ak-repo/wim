import * as React from "react"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { Button } from "@/components/ui/Button"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Package,
  Warehouse,
  MapPin,
  Boxes,
  Database,
  History,
  ArrowUpRight,
  ArrowDownRight,
  Zap,
  TrendingUp,
} from "lucide-react"
import { useProducts } from "@/features/products/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { useLocations } from "@/features/locations/hooks"
import { useInventoryList, useStockMovements } from "@/features/inventory/hooks"

function StatCard({
  title,
  value,
  description,
  icon: Icon,
  trend,
  loading,
}: {
  title: string
  value: string
  description: string
  icon: React.ComponentType<{ className?: string }>
  trend?: "up" | "down"
  loading?: boolean
}) {
  if (loading) {
    return (
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-4 w-4" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-16 mb-2" />
          <Skeleton className="h-3 w-32" />
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className="relative overflow-hidden group">
      <div className="absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
          <Icon className="h-4 w-4 text-primary" />
        </div>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold tracking-tight">{value}</div>
        <div className="flex items-center gap-2 mt-1">
          {trend && (
            trend === "up" ? (
              <ArrowUpRight className="h-3 w-3 text-emerald-500" />
            ) : (
              <ArrowDownRight className="h-3 w-3 text-red-500" />
            )
          )}
          <p className="text-xs text-muted-foreground">{description}</p>
        </div>
      </CardContent>
    </Card>
  )
}

function QuickActionCard({
  href,
  icon: Icon,
  title,
  description,
}: {
  href: string
  icon: React.ComponentType<{ className?: string }>
  title: string
  description: string
}) {
  return (
    <a
      href={href}
      className="group flex items-start gap-4 rounded-xl border border-border bg-background p-4 transition-all duration-200 hover:border-primary/50 hover:shadow-md hover:shadow-primary/5"
    >
      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 transition-colors group-hover:bg-primary/20">
        <Icon className="h-5 w-5 text-primary" />
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <p className="text-sm font-medium text-foreground">{title}</p>
          <ArrowUpRight className="h-3 w-3 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity" />
        </div>
        <p className="text-xs text-muted-foreground mt-0.5">{description}</p>
      </div>
    </a>
  )
}

function MovementItem({
  movement,
}: {
  movement: {
    id: number
    movementType: string
    quantity: number
    createdAt: string
  }
}) {
  const isPositive = movement.quantity > 0

  return (
    <div className="group flex items-center justify-between rounded-lg border border-border bg-background p-3 transition-all hover:border-primary/30 hover:bg-muted/50">
      <div className="flex items-center gap-3">
        <div
          className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-full transition-colors ${
            isPositive
              ? "bg-emerald-500/10 text-emerald-500"
              : "bg-red-500/10 text-red-500"
          }`}
        >
          {isPositive ? (
            <ArrowUpRight className="h-4 w-4" />
          ) : (
            <ArrowDownRight className="h-4 w-4" />
          )}
        </div>
        <div>
          <p className="text-sm font-medium text-foreground">{movement.movementType}</p>
          <p className="text-xs text-muted-foreground">
            {new Date(movement.createdAt).toLocaleString()}
          </p>
        </div>
      </div>
      <div
        className={`text-sm font-bold ${
          isPositive ? "text-emerald-500" : "text-red-500"
        }`}
      >
        {isPositive ? "+" : ""}
        {movement.quantity}
      </div>
    </div>
  )
}

export default function DashboardPage() {
  const { data: productsData, isLoading: productsLoading } = useProducts({ page: 1, limit: 1 })
  const { data: warehousesData, isLoading: warehousesLoading } = useWarehouses({ page: 1, limit: 1 })
  const { data: locationsData, isLoading: locationsLoading } = useLocations({ page: 1, limit: 1 })
  const { data: inventoryData, isLoading: inventoryLoading } = useInventoryList({ page: 1, limit: 1000 })
  const { data: movementsData, isLoading: movementsLoading } = useStockMovements({ page: 1, limit: 5 })

  const totalStock = inventoryData?.data?.reduce((acc, curr) => acc + curr.quantity, 0) || 0

  const stats = [
    {
      title: "Total Products",
      value: productsData?.data?.length ? productsData.data.length.toString() : "0",
      description: "Active products in catalog",
      icon: Package,
      loading: productsLoading,
    },
    {
      title: "Total Stock",
      value: totalStock.toLocaleString(),
      description: "Units across all locations",
      icon: Database,
      loading: inventoryLoading,
    },
    {
      title: "Warehouses",
      value: warehousesData?.data?.length ? warehousesData.data.length.toString() : "0",
      description: "Active warehouses",
      icon: Warehouse,
      loading: warehousesLoading,
    },
    {
      title: "Locations",
      value: locationsData?.data?.length ? locationsData.data.length.toString() : "0",
      description: "Storage locations",
      icon: MapPin,
      loading: locationsLoading,
    },
  ]

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Dashboard</h1>
          <p className="text-muted-foreground mt-1">
            Overview of your warehouse inventory and operations.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm">
            <TrendingUp className="h-4 w-4 mr-2" />
            Reports
          </Button>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <StatCard
            key={stat.title}
            {...stat}
            trend={stat.title === "Total Stock" ? "up" : undefined}
          />
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-1">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Quick Actions</CardTitle>
                <CardDescription>Common tasks and operations</CardDescription>
              </div>
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                <Zap className="h-4 w-4 text-primary" />
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-3">
            <QuickActionCard
              href="/inventory"
              icon={Database}
              title="Adjust Stock"
              description="Update inventory levels"
            />
            <QuickActionCard
              href="/products"
              icon={Package}
              title="Manage Products"
              description="Add or edit products"
            />
            <QuickActionCard
              href="/locations"
              icon={MapPin}
              title="Configure Locations"
              description="Set up storage areas"
            />
            <QuickActionCard
              href="/warehouses"
              icon={Warehouse}
              title="Manage Warehouses"
              description="Add or edit warehouses"
            />
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Recent Stock Movements</CardTitle>
                <CardDescription>Latest inventory transactions</CardDescription>
              </div>
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                <History className="h-4 w-4 text-primary" />
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {movementsLoading ? (
              <div className="space-y-3">
                {[1, 2, 3, 4, 5].map((i) => (
                  <Skeleton key={i} className="h-16 w-full" />
                ))}
              </div>
            ) : movementsData?.data && movementsData.data.length > 0 ? (
              <div className="space-y-3">
                {movementsData.data.map((m) => (
                  <MovementItem key={m.id} movement={m} />
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-4">
                  <Boxes className="h-6 w-6 text-muted-foreground" />
                </div>
                <p className="text-sm font-medium text-foreground">No recent activity</p>
                <p className="text-xs text-muted-foreground mt-1">
                  Stock movements will appear here
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
