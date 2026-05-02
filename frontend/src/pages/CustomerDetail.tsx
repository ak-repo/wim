import * as React from "react"
import { useNavigate, useParams } from "react-router-dom"
import { useCustomer } from "@/features/customers/hooks"
import { CustomerFormDialog } from "@/features/customers/components/CustomerFormDialog"
import { CustomerDeleteDialog } from "@/features/customers/components/CustomerDeleteDialog"
import type { Customer } from "@/features/customers/types"
import { Button } from "@/components/ui/Button"
import { Badge } from "@/components/ui/Badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/Card"
import { AlertMessage } from "@/components/ui/Alert"
import { Loader2, Pencil, Trash2, ArrowLeft, User, Mail, Phone, MapPin, Hash } from "lucide-react"
import { formatDateTime } from "@/utils"

export default function CustomerDetailPage() {
  const navigate = useNavigate()
  const params = useParams()
  const customerId = Number(params.id)
  const [customerSnapshot, setCustomerSnapshot] = React.useState<Customer | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = React.useState(false)

  const { data: customer, isLoading, error } = useCustomer(customerId)

  React.useEffect(() => {
    if (customer) {
      setCustomerSnapshot(customer)
    }
  }, [customer])

  const displayedCustomer = customer || customerSnapshot

  const handleEdit = () => {
    if (displayedCustomer) {
      setCustomerSnapshot(displayedCustomer)
      setIsFormOpen(true)
    }
  }

  const handleDelete = () => {
    if (displayedCustomer) {
      setCustomerSnapshot(displayedCustomer)
      setIsDeleteOpen(true)
    }
  }

  if (!customerId || Number.isNaN(customerId)) {
    return <AlertMessage variant="destructive" message="Invalid customer ID." />
  }

  if (isLoading && !displayedCustomer) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    )
  }

  if (error && !displayedCustomer) {
    return (
      <div className="space-y-4">
        <Button variant="outline" onClick={() => navigate("/masters/customers")}>
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Customers
        </Button>
        <AlertMessage
          variant="destructive"
          message={error instanceof Error ? error.message : "Failed to load customer."}
        />
      </div>
    )
  }

  if (!displayedCustomer) {
    return null
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div className="space-y-2">
          <Button variant="outline" onClick={() => navigate("/masters/customers")}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Customers
          </Button>
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Customer Details</h2>
            <p className="text-muted-foreground">Read-only overview of the customer profile.</p>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Badge variant={displayedCustomer.isActive ? "success" : "secondary"}>
            {displayedCustomer.isActive ? "Active" : "Inactive"}
          </Badge>
          <Button variant="outline" onClick={handleEdit}>
            <Pencil className="mr-2 h-4 w-4" />
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            <Trash2 className="mr-2 h-4 w-4" />
            Delete
          </Button>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-primary/10 p-2 text-primary">
                <User className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Name</p>
                <p className="text-lg font-semibold">{displayedCustomer.name}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-primary/10 p-2 text-primary">
                <Mail className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Email</p>
                <p className="text-lg font-semibold">{displayedCustomer.email}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-primary/10 p-2 text-primary">
                <Hash className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Ref Code</p>
                <p className="text-lg font-semibold">{displayedCustomer.refCode}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="pt-6">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-primary/10 p-2 text-primary">
                <Phone className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Contact</p>
                <p className="text-lg font-semibold">{displayedCustomer.contact || "-"}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Profile</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Address</p>
              <p className="text-sm">{displayedCustomer.address || "No address provided."}</p>
            </div>

            <div className="space-y-1">
              <p className="text-sm text-muted-foreground">Status</p>
              <Badge variant={displayedCustomer.isActive ? "success" : "secondary"}>
                {displayedCustomer.isActive ? "Active" : "Inactive"}
              </Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>System Info</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-start gap-3">
              <MapPin className="mt-0.5 h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">Created</p>
                <p className="text-xs text-muted-foreground">{formatDateTime(displayedCustomer.createdAt)}</p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <MapPin className="mt-0.5 h-4 w-4 text-muted-foreground" />
              <div>
                <p className="text-sm font-medium">Updated</p>
                <p className="text-xs text-muted-foreground">{formatDateTime(displayedCustomer.updatedAt)}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <CustomerFormDialog
        customer={displayedCustomer}
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
      />

      <CustomerDeleteDialog
        customer={displayedCustomer}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
        onDeleted={() => navigate("/masters/customers")}
      />
    </div>
  )
}
