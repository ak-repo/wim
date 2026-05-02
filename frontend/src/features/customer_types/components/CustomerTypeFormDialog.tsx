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
import { Textarea } from "@/components/ui/Textarea"
import { useCreateCustomerType, useUpdateCustomerType } from "@/features/customer_types/hooks"
import type { CustomerType, CreateCustomerTypeRequest, UpdateCustomerTypeRequest } from "@/features/customer_types/types"

interface CustomerTypeFormDialogProps {
  customerType: CustomerType | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const CustomerTypeFormDialog: React.FC<CustomerTypeFormDialogProps> = ({
  customerType,
  open,
  onOpenChange,
}) => {
  const createCustomerType = useCreateCustomerType()
  const updateCustomerType = useUpdateCustomerType()
  const isEditing = !!customerType

  const [formData, setFormData] = React.useState({
    name: customerType?.name || "",
    description: customerType?.description || "",
    isActive: customerType?.isActive ?? true,
  })

  React.useEffect(() => {
    setFormData({
      name: customerType?.name || "",
      description: customerType?.description || "",
      isActive: customerType?.isActive ?? true,
    })
  }, [customerType, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type, checked } = e.target as HTMLInputElement
    setFormData((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && customerType) {
        const updateData: UpdateCustomerTypeRequest = {
          name: formData.name,
          description: formData.description || undefined,
          isActive: formData.isActive,
        }
        await updateCustomerType.mutateAsync({ id: customerType.id, data: updateData })
      } else {
        const createData: CreateCustomerTypeRequest = {
          name: formData.name,
          description: formData.description || undefined,
          isActive: formData.isActive,
        }
        await createCustomerType.mutateAsync(createData)
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
          <DialogTitle>{isEditing ? "Edit Customer Type" : "Create Customer Type"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update customer type details."
              : "Add a new customer type to the system."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Name *</label>
            <Input
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Regular Customer"
              required
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Description</label>
            <Textarea
              name="description"
              value={formData.description}
              onChange={handleChange}
              placeholder="Description for customer type"
              rows={3}
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
            <Button type="submit" disabled={createCustomerType.isPending || updateCustomerType.isPending}>
              {createCustomerType.isPending || updateCustomerType.isPending
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
