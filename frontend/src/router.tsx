import { createBrowserRouter, Navigate, Outlet } from "react-router-dom"
import { QueryClient } from "@tanstack/react-query"
import { AdminLayout } from "@/components/layout/AdminLayout"
import { useAuthStore } from "@/stores/authStore"

// Pages
import LoginPage from "@/pages/Login"
import DashboardPage from "@/pages/Dashboard"
import UsersPage from "@/pages/Users"
import ProductsPage from "@/pages/Products"
import WarehousesPage from "@/pages/Warehouses"
import LocationsPage from "@/pages/Locations"
import InventoryPage from "@/pages/Inventory"
import SalesOrdersPage from "@/pages/SalesOrders"

// Protected Route wrapper - with hydration check
function ProtectedRoute() {
  const accessToken = localStorage.getItem("accessToken")
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  
  console.log("[Router] ProtectedRoute check - isAuthenticated:", isAuthenticated, "hasToken:", !!accessToken)
  
  // Must have token in localStorage (even before zustand rehydrates)
  if (!accessToken) {
    console.log("[Router] No token - redirecting to login")
    return <Navigate to="/login" replace />
  }
  return <Outlet />
}

// Public Route wrapper (redirects to dashboard if authenticated)
function PublicRoute() {
  const accessToken = localStorage.getItem("accessToken")
  
  if (!accessToken) {
    return <Outlet />
  }
  return <Navigate to="/" replace />
}

function RoleRoute({ allowedRoles }: { allowedRoles: string[] }) {
  const accessToken = localStorage.getItem("accessToken")
  const role = useAuthStore((state) => state.user?.role)

  if (!accessToken) {
    return <Navigate to="/login" replace />
  }

  if (role && !allowedRoles.includes(role) && role !== "super-admin") {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}

// Admin Layout wrapper
function AdminLayoutWrapper() {
  return (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  )
}

// Create query client
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

export const router = createBrowserRouter([
  {
    element: <PublicRoute />,
    children: [
      {
        path: "/login",
        element: <LoginPage />,
      },
    ],
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <AdminLayoutWrapper />,
        children: [
          {
            element: <RoleRoute allowedRoles={["admin"]} />,
            children: [
              {
                path: "/",
                element: <DashboardPage />,
              },
              {
                path: "/products",
                element: <ProductsPage />,
              },
              {
                path: "/warehouses",
                element: <WarehousesPage />,
              },
              {
                path: "/locations",
                element: <LocationsPage />,
              },
              {
                path: "/inventory",
                element: <InventoryPage />,
              },
              {
                path: "/sales-orders",
                element: <SalesOrdersPage />,
              },
            ],
          },
          {
            element: <RoleRoute allowedRoles={["super-admin"]} />,
            children: [
              {
                path: "/users",
                element: <UsersPage />,
              },
            ],
          },
        ],
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/" replace />,
  },
])
