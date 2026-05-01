import * as React from "react"
import { useNavigate } from "react-router-dom"
import { useCustomers } from "@/features/customers/hooks"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Select } from "@/components/ui/Select"
import { Badge } from "@/components/ui/Badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Card, CardContent, CardHeader } from "@/components/ui/Card"
import { CustomerFormDialog } from "@/features/customers/components/CustomerFormDialog"
import { CustomerDeleteDialog } from "@/features/customers/components/CustomerDeleteDialog"
import type { Customer } from "@/features/customers/types"
import { Plus, Search, Pencil, Trash2, Loader2, Eye, Users } from "lucide-react"
import { formatDate } from "@/utils"

type CustomerFilter = "all" | "active" | "inactive"

export default function CustomersPage() {
  const navigate = useNavigate()
  const [search, setSearch] = React.useState("")
  const [filter, setFilter] = React.useState<CustomerFilter>("all")
  const [page, setPage] = React.useState(1)
  const [selectedCustomer, setSelectedCustomer] = React.useState<Customer | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = React.useState(false)

  const activeParam = filter === "all" ? undefined : filter === "active"

  const { data: customersData, isLoading } = useCustomers({
    page,
    limit: 10,
    active: activeParam,
  })

  React.useEffect(() => {
    setPage(1)
  }, [filter])

  const handleEdit = (customer: Customer) => {
    setSelectedCustomer(customer)
    setIsFormOpen(true)
  }

  const handleDelete = (customer: Customer) => {
    setSelectedCustomer(customer)
    setIsDeleteOpen(true)
  }

  const handleCreate = () => {
    setSelectedCustomer(null)
    setIsFormOpen(true)
  }

  const filteredCustomers =
    customersData?.data?.filter(
      (customer) =>
        customer.name.toLowerCase().includes(search.toLowerCase()) ||
        customer.email.toLowerCase().includes(search.toLowerCase()) ||
        customer.refCode.toLowerCase().includes(search.toLowerCase()) ||
        (customer.contact?.toLowerCase() || "").includes(search.toLowerCase()) ||
        (customer.address?.toLowerCase() || "").includes(search.toLowerCase())
    ) || []

  const totalPages = customersData?.totalPages || 1
  const totalCustomers = customersData?.total || 0

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Customers</h2>
          <p className="text-muted-foreground">Manage customer accounts and profile details.</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Add Customer
        </Button>
      </div>

      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-primary/10 p-2 text-primary">
                <Users className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Customers</p>
                <p className="text-2xl font-bold">{totalCustomers}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div>
              <p className="text-sm text-muted-foreground">Current Filter</p>
              <p className="text-2xl font-bold capitalize">{filter}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div>
              <p className="text-sm text-muted-foreground">Page</p>
              <p className="text-2xl font-bold">
                {page} / {totalPages}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex flex-col gap-3 lg:flex-row lg:items-center">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search customers..."
                className="pl-9"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>

            <Select
              value={filter}
              onChange={(e) => setFilter(e.target.value as CustomerFilter)}
              className="lg:w-40"
            >
              <option value="all">All</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
            </Select>
          </div>
        </CardHeader>

        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Ref Code</TableHead>
                  <TableHead>Customer</TableHead>
                  <TableHead>Contact</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center">
                      <Loader2 className="mx-auto h-6 w-6 animate-spin text-muted-foreground" />
                    </TableCell>
                  </TableRow>
                ) : filteredCustomers.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center">
                      <div className="flex flex-col items-center justify-center text-muted-foreground">
                        <Users className="mb-2 h-8 w-8" />
                        <p>No customers found.</p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredCustomers.map((customer) => (
                    <TableRow key={customer.id}>
                      <TableCell className="font-medium">{customer.refCode}</TableCell>
                      <TableCell>
                        <div>
                          <p className="font-medium">{customer.name}</p>
                          <p className="text-sm text-muted-foreground">{customer.email}</p>
                          {customer.address && (
                            <p className="text-xs text-muted-foreground">{customer.address}</p>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="text-sm">
                          {customer.contact || <span className="text-muted-foreground">-</span>}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant={customer.isActive ? "success" : "secondary"}>
                          {customer.isActive ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell>{formatDate(customer.createdAt)}</TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => navigate(`/customers/${customer.id}`)}
                            title="View customer"
                          >
                            <Eye className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleEdit(customer)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(customer)}
                            className="text-destructive hover:text-destructive"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          <div className="mt-4 flex items-center justify-between gap-3">
            <p className="text-sm text-muted-foreground">
              Page {page} of {totalPages}
            </p>
            <div className="flex items-center gap-2">
              <Button variant="outline" onClick={() => setPage((current) => Math.max(1, current - 1))} disabled={page <= 1 || isLoading}>
                Previous
              </Button>
              <Button
                variant="outline"
                onClick={() => setPage((current) => Math.min(totalPages, current + 1))}
                disabled={page >= totalPages || isLoading}
              >
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <CustomerFormDialog customer={selectedCustomer} open={isFormOpen} onOpenChange={setIsFormOpen} />

      <CustomerDeleteDialog
        customer={selectedCustomer}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </div>
  )
}
