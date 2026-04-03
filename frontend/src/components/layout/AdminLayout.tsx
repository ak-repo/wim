import * as React from "react"
import { useLocation } from "react-router-dom"
import {
  LayoutDashboard,
  Users,
  Package,
  Warehouse,
  MapPin,
  Menu,
  X,
  LogOut,
  ChevronLeft,
} from "lucide-react"
import { cn } from "@/utils"
import { useLogout } from "@/features/auth/hooks"

interface SidebarItem {
  name: string
  href: string
  icon: React.ComponentType<{ className?: string }>
}

const navigation: SidebarItem[] = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  { name: "Users", href: "/users", icon: Users },
  { name: "Products", href: "/products", icon: Package },
  { name: "Warehouses", href: "/warehouses", icon: Warehouse },
  { name: "Locations", href: "/locations", icon: MapPin },
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
      {/* Mobile sidebar backdrop */}
      {open && (
        <div
          className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          onClick={() => setOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside
        className={cn(
          "fixed left-0 top-0 z-50 h-screen bg-card border-r border-border transition-all duration-300",
          collapsed && "w-16",
          !collapsed && "w-64",
          open && "translate-x-0",
          !open && "-translate-x-full lg:translate-x-0"
        )}
      >
        <div className="flex h-full flex-col">
          {/* Logo area */}
          <div className="flex h-14 items-center justify-between border-b border-border px-4">
            <div
              className={cn(
                "flex items-center gap-2 font-bold text-lg overflow-hidden whitespace-nowrap transition-all",
                collapsed && "w-0 opacity-0",
                !collapsed && "w-auto opacity-100"
              )}
            >
              <Package className="h-6 w-6 text-primary flex-shrink-0" />
              <span className="text-foreground">WIM Admin</span>
            </div>

            {/* Collapse toggle - desktop only */}
            <button
              onClick={() => setCollapsed(!collapsed)}
              className={cn(
                "hidden lg:flex items-center justify-center h-7 w-7 rounded-md hover:bg-muted text-muted-foreground transition-colors",
                collapsed && "rotate-180"
              )}
            >
              <ChevronLeft className="h-4 w-4" />
            </button>

            {/* Close button - mobile only */}
            <button
              onClick={() => setOpen(false)}
              className="lg:hidden flex items-center justify-center h-7 w-7 rounded-md hover:bg-muted text-muted-foreground"
            >
              <X className="h-4 w-4" />
            </button>
          </div>

          {/* Navigation */}
          <nav className="flex-1 overflow-y-auto p-3">
            <ul className="space-y-1">
              {navigation.map((item) => {
                const isActive = location.pathname === item.href
                const Icon = item.icon
                return (
                  <li key={item.name}>
                    <a
                      href={item.href}
                      className={cn(
                        "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors",
                        isActive
                          ? "bg-primary text-primary-foreground"
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
          </nav>

          {/* Logout button */}
          <div className="border-t border-border p-3">
            <button
              onClick={logout}
              className={cn(
                "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-muted hover:text-foreground transition-colors w-full",
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
    <header className="sticky top-0 z-30 flex h-14 items-center gap-4 border-b border-border bg-card px-4 shadow-sm">
      <button
        onClick={onMenuClick}
        className="lg:hidden flex items-center justify-center h-8 w-8 rounded-md hover:bg-muted text-muted-foreground"
      >
        <Menu className="h-5 w-5" />
      </button>

      <div className="flex-1">
        <h1 className="text-lg font-semibold text-foreground">Warehouse Management</h1>
      </div>

      <div className="flex items-center gap-4">
        <div className="text-sm text-muted-foreground">
          Admin Panel
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

        <main className="p-4 lg:p-6">
          <div className="mx-auto max-w-7xl">{children}</div>
        </main>
      </div>
    </div>
  )
}
