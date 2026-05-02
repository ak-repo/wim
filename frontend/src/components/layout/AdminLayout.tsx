import * as React from "react"
import { useLocation, useNavigate } from "react-router-dom"
import {
  LayoutDashboard,
  Users,
  UserRound,
  Package,
  Warehouse,
  MapPin,
  Menu,
  X,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Shield,
  Tag,
} from "lucide-react"
import { cn } from "@/utils"
import { useLogout } from "@/features/auth/hooks"

interface SidebarItem {
  name: string
  href: string
  icon: React.ComponentType<{ className?: string }>
}

interface SidebarGroup {
  name: string
  icon: React.ComponentType<{ className?: string }>
  items: SidebarItem[]
  defaultPath: string
}

const navigation: (SidebarItem | SidebarGroup)[] = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  {
    name: "User Master",
    icon: Shield,
    items: [
      { name: "Users", href: "/masters/users", icon: Users },
      { name: "User Roles", href: "/masters/user-roles", icon: Tag },
    ],
    defaultPath: "/masters/users",
  },
  {
    name: "Product Master",
    icon: Package,
    items: [
      { name: "Products", href: "/masters/products", icon: Package },
      { name: "Product Categories", href: "/masters/product-categories", icon: Tag },
    ],
    defaultPath: "/masters/products",
  },
  {
    name: "Customer Master",
    icon: UserRound,
    items: [
      { name: "Customers", href: "/masters/customers", icon: UserRound },
      { name: "Customer Types", href: "/masters/customer-types", icon: Tag },
    ],
    defaultPath: "/masters/customers",
  },
  {
    name: "Warehouse Master",
    icon: Warehouse,
    items: [
      { name: "Warehouses", href: "/masters/warehouses", icon: Warehouse },
      { name: "Locations", href: "/masters/locations", icon: MapPin },
    ],
    defaultPath: "/masters/warehouses",
  },
]

const Sidebar: React.FC<{
  open: boolean
  setOpen: (open: boolean) => void
  collapsed: boolean
  setCollapsed: (collapsed: boolean) => void
}> = ({ open, setOpen, collapsed, setCollapsed }) => {
  const location = useLocation()
  const navigate = useNavigate()
  const logout = useLogout()
  const [expandedGroups, setExpandedGroups] = React.useState<Set<string>>(new Set())

  const isGroupActive = (group: SidebarGroup) => {
    return group.items.some(
      (item) => location.pathname === item.href || location.pathname.startsWith(`${item.href}/`)
    )
  }

  const isItemActive = (item: SidebarItem) => {
    return location.pathname === item.href
  }

  const toggleGroup = (groupName: string, group: SidebarGroup) => {
    setExpandedGroups((prev) => {
      const newSet = new Set(prev)
      if (newSet.has(groupName)) {
        newSet.delete(groupName)
      } else {
        newSet.add(groupName)
        navigate(group.defaultPath)
      }
      return newSet
    })
  }

  React.useEffect(() => {
    const activeGroup = navigation.find((item) => {
      if ("items" in item) {
        return isGroupActive(item)
      }
      return false
    })
    if (activeGroup && "name" in activeGroup) {
      setExpandedGroups((prev) => new Set([...prev, activeGroup.name]))
    }
  }, [location.pathname])

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
                if ("items" in item) {
                  const isExpanded = expandedGroups.has(item.name)
                  const isActive = isGroupActive(item)
                  const Icon = item.icon
                  return (
                    <li key={item.name}>
                      <button
                        onClick={() => toggleGroup(item.name, item)}
                        className={cn(
                          "flex w-full items-center gap-3 rounded-md border border-transparent px-3 py-2 text-sm font-medium transition-colors",
                          isActive
                            ? "border-primary/40 bg-primary/20 text-primary"
                            : "text-muted-foreground hover:border-border hover:bg-muted/60 hover:text-foreground",
                          collapsed && "justify-center px-2"
                        )}
                        title={collapsed ? item.name : undefined}
                      >
                        <Icon className="h-5 w-5 flex-shrink-0" />
                        <span
                          className={cn(
                            "flex-1 text-left overflow-hidden whitespace-nowrap transition-all",
                            collapsed && "w-0 opacity-0",
                            !collapsed && "w-auto opacity-100"
                          )}
                        >
                          {item.name}
                        </span>
                        {!collapsed && (
                          <ChevronRight
                            className={cn(
                              "h-4 w-4 flex-shrink-0 transition-transform",
                              isExpanded && "rotate-90"
                            )}
                          />
                        )}
                      </button>
                      {!collapsed && isExpanded && (
                        <ul className="mt-1 ml-4 space-y-1">
                          {item.items.map((subItem) => {
                            const isSubItemActive = isItemActive(subItem)
                            const SubIcon = subItem.icon
                            return (
                              <li key={subItem.name}>
                                <a
                                  href={subItem.href}
                                  className={cn(
                                    "flex items-center gap-3 rounded-md border border-transparent px-3 py-2 text-sm font-medium transition-colors",
                                    isSubItemActive
                                      ? "border-primary/40 bg-primary/20 text-primary"
                                      : "text-muted-foreground hover:border-border hover:bg-muted/60 hover:text-foreground"
                                  )}
                                >
                                  <SubIcon className="h-4 w-4 flex-shrink-0" />
                                  <span className="overflow-hidden whitespace-nowrap">
                                    {subItem.name}
                                  </span>
                                </a>
                              </li>
                            )
                          })}
                        </ul>
                      )}
                    </li>
                  )
                } else {
                  const isActive = isItemActive(item)
                  const Icon = item.icon
                  return (
                    <li key={item.name}>
                      <a
                        href={item.href}
                        className={cn(
                          "flex items-center gap-3 rounded-md border border-transparent px-3 py-2 text-sm font-medium transition-colors",
                          isActive
                            ? "border-primary/40 bg-primary/20 text-primary"
                            : "text-muted-foreground hover:border-border hover:bg-muted/60 hover:text-foreground",
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
                }
              })}
            </ul>
          </nav>

          {/* Logout button */}
          <div className="border-t border-border p-3">
            <button
              onClick={logout}
              className={cn(
                "flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground",
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
    <header className="sticky top-0 z-30 flex h-14 items-center gap-4 border-b border-border bg-card px-4 shadow-none">
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

        <main className="p-6">
          <div className="mx-auto max-w">{children}</div>
        </main>
      </div>
    </div>
  )
}
