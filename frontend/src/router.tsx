import { createBrowserRouter, Navigate, Outlet } from "react-router-dom"
import { QueryClient } from "@tanstack/react-query"
import { AdminLayout } from "@/components/layout/AdminLayout"
import { useAuthStore } from "@/stores/authStore"

// Pages
import LoginPage from "@/pages/Login"
import DashboardPage from "@/pages/Dashboard"
import UsersPage from "@/pages/Users"
import UserRolesPage from "@/pages/UserRoles"
import CustomersPage from "@/pages/Customers"
import CustomerTypesPage from "@/pages/CustomerTypes"
import CustomerDetailPage from "@/pages/CustomerDetail"
import ProductsPage from "@/pages/Products"
import ProductCategoriesPage from "@/pages/ProductCategories"
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
            path: "/masters",
            element: <Navigate to="/masters/users" replace />,
          },
          {
            path: "/",
            element: <DashboardPage />,
          },
          {
            path: "/masters/users",
            element: <UsersPage />,
          },
          {
            path: "/masters/user-roles",
            element: <UserRolesPage />,
          },
          {
            path: "/masters/customers",
            element: <CustomersPage />,
          },
          {
            path: "/masters/customer-types",
            element: <CustomerTypesPage />,
          },
          {
            path: "/masters/customers/:id",
            element: <CustomerDetailPage />,
          },
          {
            path: "/masters/products",
            element: <ProductsPage />,
          },
          {
            path: "/masters/product-categories",
            element: <ProductCategoriesPage />,
          },
          {
            path: "/masters/warehouses",
            element: <WarehousesPage />,
          },
          {
            path: "/masters/locations",
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
