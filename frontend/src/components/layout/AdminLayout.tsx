import * as React from "react"
import { useLocation, Link } from "react-router-dom"
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
  ShoppingCart,
} from "lucide-react"
import { cn } from "@/utils"
import { useLogout } from "@/features/auth/hooks"
import { ScrollArea } from "@/components/ui/scroll-area"
import { useAuthStore } from "@/stores/authStore"

interface SidebarItem {
  name: string
  href: string
  icon: React.ComponentType<{ className?: string; strokeWidth?: number }>
  roles?: string[]
}

const navigation: SidebarItem[] = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard, roles: ["admin", "super-admin"] },
  { name: "User Master", href: "/users", icon: Users, roles: ["super-admin"] },
  { name: "Product Master", href: "/products", icon: Layers, roles: ["admin", "super-admin"] },
  { name: "Warehouses", href: "/warehouses", icon: Warehouse, roles: ["admin", "super-admin"] },
  { name: "Locations", href: "/locations", icon: MapPin, roles: ["admin", "super-admin"] },
  { name: "Inventory", href: "/inventory", icon: Database, roles: ["admin", "super-admin"] },
  { name: "Sales Orders", href: "/sales-orders", icon: ShoppingCart, roles: ["admin", "super-admin"] },
]

const Sidebar: React.FC<{
  open: boolean
  setOpen: (open: boolean) => void
}> = ({ open, setOpen }) => {
  const location = useLocation()
  const logout = useLogout()
  const role = useAuthStore((state) => state.user?.role)

  const visibleNavigation = navigation.filter((item) => {
    if (!item.roles || item.roles.length === 0) {
      return true
    }
    if (!role) {
      return true
    }
    return item.roles.includes(role)
  })

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
          "fixed left-0 top-0 z-50 h-screen bg-ink flex flex-col transition-transform duration-300 w-[200px]",
          open ? "translate-x-0" : "-translate-x-full lg:translate-x-0"
        )}
      >
        {/* Brand */}
        <div className="flex items-center gap-2 h-[52px] px-4 shrink-0">
          <div className="flex h-5 w-5 shrink-0 items-center justify-center rounded-[6px] bg-accent-green">
            <Package className="h-3 w-3 text-white" />
          </div>
          <div className="flex flex-col justify-center">
            <span className="text-white font-medium text-[13px] leading-tight">WIM</span>
            <span className="text-white/50 text-[9px] uppercase tracking-wider leading-tight">Warehouse</span>
          </div>
          <button
            onClick={() => setOpen(false)}
            className="lg:hidden ml-auto flex items-center justify-center h-7 w-7 rounded-lg text-white/50"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        <ScrollArea className="flex-1 px-3 py-2">
          <div className="mb-2 px-2">
            <span className="text-[9px] text-white/30 uppercase tracking-widest font-medium">Menu</span>
          </div>
          <ul className="space-y-0.5">
            {visibleNavigation.map((item) => {
              const isActive = location.pathname === item.href
              const Icon = item.icon
              return (
                <li key={item.name}>
                  <Link
                    to={item.href}
                    className={cn(
                      "flex items-center gap-2 rounded-[6px] px-2 py-[7px] text-[13px] font-medium transition-colors relative",
                      isActive
                        ? "bg-white/10 text-white"
                        : "text-white/50 hover:bg-white/5 hover:text-white"
                    )}
                  >
                    <Icon className="h-4 w-4 shrink-0" strokeWidth={2.5} />
                    <span className="truncate">{item.name}</span>
                  </Link>
                </li>
              )
            })}
          </ul>
        </ScrollArea>

        {/* Footer */}
        <div className="border-t-[0.5px] border-white/10 p-3 shrink-0">
          <div className="flex items-center gap-2 mb-3">
            <div className="h-[26px] w-[26px] shrink-0 rounded-[6px] bg-accent-green flex items-center justify-center">
              <span className="text-[11px] font-medium text-white">AK</span>
            </div>
            <div className="flex flex-col overflow-hidden">
              <span className="text-[12px] font-medium text-white/90 truncate">
                {useAuthStore.getState().user?.username || "Admin User"}
              </span>
              <span className="text-[10px] text-white/50 truncate">
                {role || "Administrator"}
              </span>
            </div>
          </div>
          <button
            onClick={logout}
            className="flex items-center gap-2 rounded-[6px] px-2 py-[7px] text-[13px] font-medium text-white/50 hover:bg-white/5 hover:text-white transition-colors w-full"
          >
            <LogOut className="h-4 w-4 shrink-0" strokeWidth={2.5} />
            <span>Logout</span>
          </button>
        </div>
      </aside>
    </>
  )
}

const Header: React.FC<{
  onMenuClick: () => void
}> = ({ onMenuClick }) => {
  return (
    <header className="sticky top-0 z-30 flex h-[52px] items-center gap-4 border-b-[0.5px] border-border-default bg-white px-4 lg:px-6 shrink-0">
      <button
        onClick={onMenuClick}
        className="lg:hidden flex items-center justify-center h-8 w-8 rounded-[7px] text-ink-3 hover:bg-surface-2 transition-colors"
      >
        <Menu className="h-4 w-4" />
      </button>

      <div className="flex items-center gap-4 flex-1">
        <span className="text-[14px] font-medium text-ink">Overview</span>
      </div>
    </header>
  )
}

export const AdminLayout: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [sidebarOpen, setSidebarOpen] = React.useState(false)

  return (
    <div className="min-h-screen bg-surface-2 flex flex-col font-medium text-ink">
      <Sidebar
        open={sidebarOpen}
        setOpen={setSidebarOpen}
      />
      <div className="flex-1 flex flex-col lg:ml-[200px] transition-all duration-300">
        <Header onMenuClick={() => setSidebarOpen(true)} />
        <main className="flex-1 p-4 lg:p-6 overflow-auto">
          <div className="w-full">
            {children}
          </div>
        </main>
      </div>
    </div>
  )
}