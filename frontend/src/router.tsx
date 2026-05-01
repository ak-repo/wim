import { createBrowserRouter, Navigate, Outlet } from "react-router-dom"
import { QueryClient } from "@tanstack/react-query"
import { AdminLayout } from "@/components/layout/AdminLayout"
import { useAuthStore } from "@/stores/authStore"

// Pages
import LoginPage from "@/pages/Login"
import DashboardPage from "@/pages/Dashboard"
import UsersPage from "@/pages/Users"
import CustomersPage from "@/pages/Customers"
import CustomerDetailPage from "@/pages/CustomerDetail"
import ProductsPage from "@/pages/Products"
import WarehousesPage from "@/pages/Warehouses"
import LocationsPage from "@/pages/Locations"

// Protected Route wrapper
function ProtectedRoute() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />
}

// Public Route wrapper (redirects to dashboard if authenticated)
function PublicRoute() {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  return !isAuthenticated ? <Outlet /> : <Navigate to="/" replace />
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
            path: "/",
            element: <DashboardPage />,
          },
          {
            path: "/users",
            element: <UsersPage />,
          },
          {
            path: "/customers",
            element: <CustomersPage />,
          },
          {
            path: "/customers/:id",
            element: <CustomerDetailPage />,
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
        ],
      },
    ],
  },
  {
    path: "*",
    element: <Navigate to="/" replace />,
  },
])
