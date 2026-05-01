import { useMemo } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Badge } from "@/components/ui/Badge"
import { useUsers } from "@/features/auth/hooks"
import { useProducts } from "@/features/products/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { useLocations } from "@/features/locations/hooks"
import {
  Package,
  Users,
  Warehouse,
  MapPin,
  BarChart3,
} from "lucide-react"

export default function DashboardPage() {
  const usersQuery = useUsers({ page: 1, limit: 1 })
  const productsQuery = useProducts({ page: 1, limit: 1 })
  const warehousesQuery = useWarehouses({ page: 1, limit: 1 })
  const locationsQuery = useLocations({ page: 1, limit: 1 })

  const stats = useMemo(
    () => [
      {
        title: "Total Products",
        value: productsQuery.data?.total,
        description: "Catalog items",
        icon: Package,
      },
      {
        title: "Total Users",
        value: usersQuery.data?.total,
        description: "Account records",
        icon: Users,
      },
      {
        title: "Warehouses",
        value: warehousesQuery.data?.total,
        description: "Connected sites",
        icon: Warehouse,
      },
      {
        title: "Locations",
        value: locationsQuery.data?.total,
        description: "Storage zones",
        icon: MapPin,
      },
    ],
    [productsQuery.data?.total, usersQuery.data?.total, warehousesQuery.data?.total, locationsQuery.data?.total]
  )

  const queryStates = [usersQuery, productsQuery, warehousesQuery, locationsQuery]
  const hasAnyMetricData = stats.some((stat) => typeof stat.value === "number")
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-semibold tracking-tight">Dashboard</h2>
        <p className="text-sm text-muted-foreground">System overview from available backend data.</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon
          const isLoading = queryStates.some((query) => query.isLoading)
          return (
            <Card key={stat.title} className="h-full">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-1">
                <CardTitle className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                  {stat.title}
                </CardTitle>
                <Icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="font-mono text-3xl font-semibold leading-none">
                  {typeof stat.value === "number" ? stat.value : isLoading ? "..." : "--"}
                </div>
                <p className="mt-2 text-xs text-muted-foreground">{stat.description}</p>
              </CardContent>
            </Card>
          )
        })}
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        <Card className="h-full lg:col-span-2">
          <CardHeader>
            <CardTitle className="text-sm font-semibold">Inventory Throughput (7d)</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex min-h-[220px] flex-col items-center justify-center rounded-md border border-border bg-[#111827] p-4 text-center">
              <BarChart3 className="mb-2 h-5 w-5 text-muted-foreground" />
              <p className="text-sm text-foreground">No throughput data available.</p>
              <p className="mt-2 text-xs text-muted-foreground">Connect a backend analytics endpoint to populate this panel.</p>
            </div>
          </CardContent>
        </Card>

        <Card className="h-full lg:col-span-1">
          <CardHeader>
            <CardTitle className="text-sm font-semibold">System Summary</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            {stats
              .filter((stat) => typeof stat.value === "number")
              .map((stat) => (
                <div key={stat.title} className="rounded-md border border-border bg-[#111827] p-3">
                  <div className="flex items-center justify-between text-xs">
                    <span className="text-muted-foreground">{stat.title}</span>
                    <Badge variant={stat.value && stat.value > 0 ? "success" : "warning"}>
                      {stat.value ?? 0}
                    </Badge>
                  </div>
                </div>
              ))}
            {!hasAnyMetricData && !queryStates.some((query) => query.isLoading) && (
              <div className="rounded-md border border-border bg-[#111827] p-3 text-xs text-muted-foreground">
                No dashboard metrics available from backend.
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-semibold">Recent Activity</CardTitle>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Event</TableHead>
                <TableHead>Actor</TableHead>
                <TableHead>Scope</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Time</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell colSpan={5} className="py-8 text-center">
                  <p className="text-sm text-foreground">No recent activity available.</p>
                  <p className="mt-2 text-xs text-muted-foreground">No backend activity feed is configured for this dashboard.</p>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  )
}
