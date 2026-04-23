import * as React from "react"
import { useLocation } from "react-router-dom"
import {
  LayoutDashboard,
  Users,
  Package,
  Layers,
  Warehouse,
  MapPin,
  Database,
  Menu,
  X,
  LogOut,
  ChevronLeft,
  ShoppingCart,
  Bell,
  Search,
} from "lucide-react"
import { cn } from "@/utils"
import { useLogout } from "@/features/auth/hooks"
import { Input } from "@/components/ui/Input"
import { Button } from "@/components/ui/Button"
import { ScrollArea } from "@/components/ui/scroll-area"

interface SidebarItem {
  name: string
  href: string
  icon: React.ComponentType<{ className?: string }>
}

const navigation: SidebarItem[] = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  { name: "User Master", href: "/users", icon: Users },
  { name: "Product Master", href: "/products", icon: Layers },
  { name: "Warehouses", href: "/warehouses", icon: Warehouse },
  { name: "Locations", href: "/locations", icon: MapPin },
  { name: "Inventory", href: "/inventory", icon: Database },
  { name: "Sales Orders", href: "/sales-orders", icon: ShoppingCart },
]

const Sidebar: React.FC<{
  open: boolean
  setOpen: (open: boolean) => void
  collapsed: boolean
  setCollapsed: (collapsed: boolean) => void
}> = ({ open, setOpen, collapsed, setCollapsed }) => {
  const location = useLocation()
  const logout = useLogout()

  return (
    <>
      {open && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden backdrop-blur-sm"
          onClick={() => setOpen(false)}
        />
      )}

      <aside
        className={cn(
          "fixed left-0 top-0 z-50 h-screen bg-background border-r border-border transition-all duration-300",
          collapsed && "w-16",
          !collapsed && "w-64",
          open && "translate-x-0",
          !open && "-translate-x-full lg:translate-x-0"
        )}
      >
        <div className="flex h-full flex-col">
          <div className="flex h-16 items-center justify-between border-b border-border px-4">
            <div
              className={cn(
                "flex items-center gap-3 font-bold text-lg overflow-hidden whitespace-nowrap transition-all",
                collapsed && "w-0 opacity-0",
                !collapsed && "w-auto opacity-100"
              )}
            >
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                <Package className="h-4 w-4 text-primary-foreground" />
              </div>
              <span className="text-foreground">WIM</span>
            </div>

            <button
              onClick={() => setCollapsed(!collapsed)}
              className={cn(
                "hidden lg:flex items-center justify-center h-7 w-7 rounded-lg hover:bg-muted text-muted-foreground transition-colors",
                collapsed && "rotate-180"
              )}
            >
              <ChevronLeft className="h-4 w-4" />
            </button>

            <button
              onClick={() => setOpen(false)}
              className="lg:hidden flex items-center justify-center h-7 w-7 rounded-lg hover:bg-muted text-muted-foreground"
            >
              <X className="h-4 w-4" />
            </button>
          </div>

          <ScrollArea className="flex-1 p-3">
            <ul className="space-y-1">
              {navigation.map((item) => {
                const isActive = location.pathname === item.href
                const Icon = item.icon
                return (
                  <li key={item.name}>
                    <a
                      href={item.href}
                      className={cn(
                        "flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all duration-200",
                        isActive
                          ? "bg-primary text-primary-foreground shadow-sm"
                          : "text-muted-foreground hover:bg-muted hover:text-foreground",
                        collapsed && "justify-center px-2"
                      )}
                      title={collapsed ? item.name : undefined}
                    >
                      <Icon className="h-5 w-5 flex-shrink-0" />
                      <span
                        className={cn(
                          "overflow-hidden whitespace-nowrap transition-all",
                          collapsed && "w-0 opacity-0",
                          !collapsed && "w-auto opacity-100"
                        )}
                      >
                        {item.name}
                      </span>
                    </a>
                  </li>
                )
              })}
            </ul>
          </ScrollArea>

          <div className="border-t border-border p-3">
            <button
              onClick={logout}
              className={cn(
                "flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition-all duration-200 w-full",
                collapsed && "justify-center px-2"
              )}
              title={collapsed ? "Logout" : undefined}
            >
              <LogOut className="h-5 w-5 flex-shrink-0" />
              <span
                className={cn(
                  "overflow-hidden whitespace-nowrap transition-all",
                  collapsed && "w-0 opacity-0",
                  !collapsed && "w-auto opacity-100"
                )}
              >
                Logout
              </span>
            </button>
          </div>
        </div>
      </aside>
    </>
  )
}

const Header: React.FC<{
  onMenuClick: () => void
}> = ({ onMenuClick }) => {
  return (
    <header className="sticky top-0 z-30 flex h-16 items-center gap-4 border-b border-border bg-background/80 backdrop-blur px-6">
      <button
        onClick={onMenuClick}
        className="lg:hidden flex items-center justify-center h-9 w-9 rounded-lg hover:bg-muted text-muted-foreground transition-colors"
      >
        <Menu className="h-5 w-5" />
      </button>

      <div className="flex items-center gap-4 flex-1">
        <div className="relative w-full max-w-sm hidden sm:block">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Search..."
            className="pl-9 bg-muted/50 border-transparent focus:bg-background"
          />
        </div>
      </div>

      <div className="flex items-center gap-2">
        <Button variant="ghost" size="icon" className="relative">
          <Bell className="h-5 w-5" />
          <span className="absolute top-2 right-2 h-2 w-2 rounded-full bg-primary" />
        </Button>
        <div className="h-8 w-px bg-border mx-2" />
        <div className="flex items-center gap-3">
          <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center">
            <span className="text-xs font-semibold text-primary">A</span>
          </div>
          <div className="hidden md:block">
            <p className="text-sm font-medium text-foreground">Admin</p>
            <p className="text-xs text-muted-foreground">Administrator</p>
          </div>
        </div>
      </div>
    </header>
  )
}

export const AdminLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [sidebarOpen, setSidebarOpen] = React.useState(false)
  const [sidebarCollapsed, setSidebarCollapsed] = React.useState(false)

  return (
    <div className="min-h-screen bg-background">
      <Sidebar
        open={sidebarOpen}
        setOpen={setSidebarOpen}
        collapsed={sidebarCollapsed}
        setCollapsed={setSidebarCollapsed}
      />

      <div
        className={cn(
          "transition-all duration-300",
          sidebarCollapsed && "lg:pl-16",
          !sidebarCollapsed && "lg:pl-64"
        )}
      >
        <Header onMenuClick={() => setSidebarOpen(true)} />

        <main className="p-6 lg:p-8">
          <div className="mx-auto max-w-7xl space-y-6">{children}</div>
        </main>
      </div>
    </div>
  )
}
