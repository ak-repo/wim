import * as React from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/Dialog"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { useCreateCustomer, useUpdateCustomer } from "@/features/customers/hooks"
import type { Customer, CreateCustomerRequest, UpdateCustomerRequest } from "@/features/customers/types"

interface CustomerFormDialogProps {
  customer: Customer | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const CustomerFormDialog: React.FC<CustomerFormDialogProps> = ({
  customer,
  open,
  onOpenChange,
}) => {
  const createCustomer = useCreateCustomer()
  const updateCustomer = useUpdateCustomer()
  const isEditing = !!customer

  const [formData, setFormData] = React.useState({
    name: customer?.name || "",
    email: customer?.email || "",
    contact: customer?.contact || "",
    address: customer?.address || "",
    isActive: customer?.isActive ?? true,
  })

  React.useEffect(() => {
    setFormData({
      name: customer?.name || "",
      email: customer?.email || "",
      contact: customer?.contact || "",
      address: customer?.address || "",
      isActive: customer?.isActive ?? true,
    })
  }, [customer, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && customer) {
        const updateData: UpdateCustomerRequest = {
          name: formData.name,
          email: formData.email,
          contact: formData.contact || undefined,
          address: formData.address || undefined,
          isActive: formData.isActive,
        }
        await updateCustomer.mutateAsync({ id: customer.id, data: updateData })
      } else {
        const createData: CreateCustomerRequest = {
          name: formData.name,
          email: formData.email,
          contact: formData.contact || undefined,
          address: formData.address || undefined,
          isActive: formData.isActive,
        }
        await createCustomer.mutateAsync(createData)
      }
      onOpenChange(false)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>{isEditing ? "Edit Customer" : "Create Customer"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update customer profile details."
              : "Register a new customer in the system."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Name *</label>
            <Input
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Acme Trading Co."
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Email *</label>
            <Input
              name="email"
              type="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="contact@acme.com"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Contact</label>
            <Input
              name="contact"
              value={formData.contact}
              onChange={handleChange}
              placeholder="+1 555 0100"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Address</label>
            <Input
              name="address"
              value={formData.address}
              onChange={handleChange}
              placeholder="123 Warehouse Way"
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              name="isActive"
              id="isActive"
              checked={formData.isActive}
              onChange={handleChange}
              className="h-4 w-4 rounded border-input"
            />
            <label htmlFor="isActive" className="text-sm">
              Active
            </label>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={createCustomer.isPending || updateCustomer.isPending}>
              {createCustomer.isPending || updateCustomer.isPending
                ? "Saving..."
                : isEditing
                ? "Update"
                : "Create"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  )
}
