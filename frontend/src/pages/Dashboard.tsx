import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import {
  Package,
  Users,
  Warehouse,
  MapPin,
  Boxes,
} from "lucide-react"

const stats = [
  {
    title: "Total Products",
    value: "0",
    description: "Active products in catalog",
    icon: Package,
    trend: "+0%",
  },
  {
    title: "Total Users",
    value: "0",
    description: "Registered users",
    icon: Users,
    trend: "+0%",
  },
  {
    title: "Warehouses",
    value: "0",
    description: "Active warehouses",
    icon: Warehouse,
    trend: "+0%",
  },
  {
    title: "Locations",
    value: "0",
    description: "Storage locations",
    icon: MapPin,
    trend: "+0%",
  },
]

export default function DashboardPage() {
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

      {/* Quick Actions */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>Quick Actions</CardTitle>
            <CardDescription>Common tasks and operations</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            <a
              href="/products"
              className="flex items-center gap-2 rounded-md border p-3 hover:bg-muted transition-colors"
            >
              <Package className="h-4 w-4" />
              <span>Add New Product</span>
            </a>
            <a
              href="/warehouses"
              className="flex items-center gap-2 rounded-md border p-3 hover:bg-muted transition-colors"
            >
              <Warehouse className="h-4 w-4" />
              <span>Manage Warehouses</span>
            </a>
            <a
              href="/locations"
              className="flex items-center gap-2 rounded-md border p-3 hover:bg-muted transition-colors"
            >
              <MapPin className="h-4 w-4" />
              <span>View Locations</span>
            </a>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Status</CardTitle>
            <CardDescription>Current system health</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="h-2 w-2 rounded-full bg-emerald-500" />
                <span className="text-sm">API Server</span>
              </div>
              <span className="text-xs text-muted-foreground">Online</span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="h-2 w-2 rounded-full bg-emerald-500" />
                <span className="text-sm">Database</span>
              </div>
              <span className="text-xs text-muted-foreground">Connected</span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <div className="h-2 w-2 rounded-full bg-emerald-500" />
                <span className="text-sm">Cache</span>
              </div>
              <span className="text-xs text-muted-foreground">Active</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
            <CardDescription>Latest system events</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col items-center justify-center py-8 text-center">
              <Boxes className="h-8 w-8 text-muted-foreground mb-2" />
              <p className="text-sm text-muted-foreground">No recent activity</p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
