import {
  Package,
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
import { useAuthStore } from "@/stores/authStore"
import { StatCard } from "@/components/ui/StatCard"
import { MovementItem } from "@/components/ui/MovementItem"

export default function DashboardPage() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  console.log("[Dashboard] Rendering - isAuthenticated:", isAuthenticated)
  
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
      description: "Active items",
      icon: Package,
      loading: productsLoading,
    },
    {
      title: "Total Stock",
      value: totalStock.toLocaleString(),
      description: "Units available",
      icon: Database,
      loading: inventoryLoading,
      trend: "up" as const,
    },
    {
      title: "Warehouses",
      value: warehousesData?.data?.length ? warehousesData.data.length.toString() : "0",
      description: "Facilities",
      icon: Warehouse,
      loading: warehousesLoading,
    },
    {
      title: "Locations",
      value: locationsData?.data?.length ? locationsData.data.length.toString() : "0",
      description: "Storage areas",
      icon: MapPin,
      loading: locationsLoading,
    },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-[20px] font-medium tracking-tight text-ink">Dashboard</h1>
        <p className="text-[12px] text-ink-3 mt-0.5">
          Overview of your warehouse inventory and operations.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-[10px]">
        {stats.map((stat) => (
          <StatCard key={stat.title} {...stat} />
        ))}
      </div>

      <div className="bg-white border-[0.5px] border-border-default rounded-[10px] flex flex-col overflow-hidden">
        <div className="p-[14px_16px] border-b-[0.5px] border-border-default flex items-center justify-between">
          <div>
            <h2 className="text-[13px] font-medium text-ink">Recent Stock Movements</h2>
            <p className="text-[10px] text-ink-3 mt-0.5">Latest inventory transactions</p>
          </div>
          <div className="flex h-[30px] w-[30px] items-center justify-center rounded-[7px] bg-blue-bg text-blue">
            <History className="h-4 w-4" />
          </div>
        </div>
        <div className="p-[14px_16px]">
          {movementsLoading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((i) => (
                <div key={i} className="h-12 w-full bg-surface-2 rounded-[7px] animate-pulse" />
              ))}
            </div>
          ) : movementsData?.data && movementsData.data.length > 0 ? (
            <div className="flex flex-col">
              {movementsData.data.map((m) => (
                <MovementItem key={m.id} movement={m} />
              ))}
            </div>
          ) : (
            <div className="flex flex-col items-center justify-center py-10 text-center">
              <div className="flex h-[36px] w-[36px] items-center justify-center rounded-[10px] bg-surface-2 mb-3">
                <Boxes className="h-5 w-5 text-ink-3" />
              </div>
              <p className="text-[12px] text-ink-3">No recent activity found</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}