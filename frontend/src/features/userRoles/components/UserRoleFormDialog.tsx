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
import {
  useCreateUserRole,
  useUpdateUserRole,
} from "@/features/userRoles/hooks"
import type {
  UserRole,
  CreateUserRoleRequest,
  UpdateUserRoleRequest,
} from "@/features/userRoles/types"

interface UserRoleFormDialogProps {
  role: UserRole | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const UserRoleFormDialog: React.FC<UserRoleFormDialogProps> = ({
  role,
  open,
  onOpenChange,
}) => {
  const createRole = useCreateUserRole()
  const updateRole = useUpdateUserRole()
  const isEditing = !!role

  const [formData, setFormData] = React.useState({
    name: role?.name || "",
    isActive: role?.isActive ?? true,
  })

  React.useEffect(() => {
    setFormData({
      name: role?.name || "",
      isActive: role?.isActive ?? true,
    })
  }, [role, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.checked }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && role) {
        const updateData: UpdateUserRoleRequest = {
          name: formData.name,
          isActive: formData.isActive,
        }
        await updateRole.mutateAsync({ id: role.id, data: updateData })
      } else {
        const createData: CreateUserRoleRequest = {
          name: formData.name,
          isActive: formData.isActive,
        }
        await createRole.mutateAsync(createData)
      }
      onOpenChange(false)
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEditing ? "Edit Role" : "Create Role"}</DialogTitle>
          <DialogDescription>
            {isEditing ? "Update role details." : "Add a new role."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Name</label>
            <Input
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Role name"
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              name="isActive"
              id="roleActive"
              checked={formData.isActive}
              onChange={handleCheckboxChange}
              className="h-4 w-4 rounded border-input"
            />
            <label htmlFor="roleActive" className="text-sm">Active</label>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={createRole.isPending || updateRole.isPending}>
              {createRole.isPending || updateRole.isPending
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
