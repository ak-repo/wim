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
import { useCreateUserRole, useUpdateUserRole } from "@/features/user_roles/hooks"
import type { UserRole, CreateUserRoleRequest, UpdateUserRoleRequest } from "@/features/user_roles/types"

interface UserRoleFormDialogProps {
  userRole: UserRole | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const UserRoleFormDialog: React.FC<UserRoleFormDialogProps> = ({
  userRole,
  open,
  onOpenChange,
}) => {
  const createUserRole = useCreateUserRole()
  const updateUserRole = useUpdateUserRole()
  const isEditing = !!userRole

  const [formData, setFormData] = React.useState({
    name: userRole?.name || "",
    description: userRole?.description || "",
    isActive: userRole?.isActive ?? true,
  })

  React.useEffect(() => {
    setFormData({
      name: userRole?.name || "",
      description: userRole?.description || "",
      isActive: userRole?.isActive ?? true,
    })
  }, [userRole, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value,
    }))
  }

  const handleTextAreaChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && userRole) {
        const updateData: UpdateUserRoleRequest = {
          name: formData.name,
          description: formData.description || undefined,
          isActive: formData.isActive,
        }
        await updateUserRole.mutateAsync({ id: userRole.id, data: updateData })
      } else {
        const createData: CreateUserRoleRequest = {
          name: formData.name,
          description: formData.description || undefined,
          isActive: formData.isActive,
        }
        await createUserRole.mutateAsync(createData)
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
          <DialogTitle>{isEditing ? "Edit User Role" : "Create User Role"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update user role details."
              : "Add a new user role to the system."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Name *</label>
            <Input
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Admin"
              required
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Description</label>
            <textarea
              name="description"
              value={formData.description}
              onChange={handleTextAreaChange}
              placeholder="Description for user role"
              rows={3}
              className="flex min-h-[80px] w-full rounded-[6px] border border-input bg-[#111827] px-3 py-2 text-sm text-foreground transition-colors placeholder:text-muted-foreground focus-visible:border-primary focus-visible:outline-none focus-visible:ring-0 disabled:cursor-not-allowed disabled:opacity-50"
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
            <Button type="submit" disabled={createUserRole.isPending || updateUserRole.isPending}>
              {createUserRole.isPending || updateUserRole.isPending
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