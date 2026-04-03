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
import { Select } from "@/components/ui/Select"
import { useCreateUser, useUpdateUser } from "@/features/auth/hooks"
import type { User, UserRequest } from "@/features/auth/types"

interface UserFormDialogProps {
  user: User | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const UserFormDialog: React.FC<UserFormDialogProps> = ({
  user,
  open,
  onOpenChange,
}) => {
  const createUser = useCreateUser()
  const updateUser = useUpdateUser()
  const isEditing = !!user

  const [formData, setFormData] = React.useState({
    username: user?.username || "",
    email: user?.email || "",
    password: "",
    role: user?.role || "admin",
    contact: user?.contact || "",
    isActive: user?.isActive ?? true,
  })

  React.useEffect(() => {
    setFormData({
      username: user?.username || "",
      email: user?.email || "",
      password: "",
      role: user?.role || "admin",
      contact: user?.contact || "",
      isActive: user?.isActive ?? true,
    })
  }, [user, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.checked }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && user) {
        await updateUser.mutateAsync({
          id: user.id,
          data: {
            username: formData.username,
            email: formData.email,
            role: formData.role,
            contact: formData.contact || undefined,
            isActive: formData.isActive,
          },
        })
      } else {
        await createUser.mutateAsync({
          username: formData.username,
          email: formData.email,
          password: formData.password,
          role: formData.role,
          contact: formData.contact || undefined,
          isActive: formData.isActive,
        } as UserRequest)
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
          <DialogTitle>{isEditing ? "Edit User" : "Create User"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update user details and permissions."
              : "Add a new user to the system."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Username</label>
            <Input
              name="username"
              value={formData.username}
              onChange={handleChange}
              placeholder="johndoe"
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Email</label>
            <Input
              name="email"
              type="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="john@example.com"
              disabled={isEditing}
            />
          </div>

          {!isEditing && (
            <div className="space-y-2">
              <label className="text-sm font-medium">Password</label>
              <Input
                name="password"
                type="password"
                value={formData.password}
                onChange={handleChange}
                placeholder="••••••••"
              />
            </div>
          )}

          <div className="space-y-2">
            <label className="text-sm font-medium">Role</label>
            <Select
              name="role"
              value={formData.role}
              onChange={handleChange}
            >
              <option value="admin">Admin</option>
              <option value="manager">Manager</option>
              <option value="worker">Worker</option>
            </Select>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Contact</label>
            <Input
              name="contact"
              value={formData.contact}
              onChange={handleChange}
              placeholder="+1 234 567 890"
            />
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              name="isActive"
              id="isActive"
              checked={formData.isActive}
              onChange={handleCheckboxChange}
              className="h-4 w-4 rounded border-input"
            />
            <label htmlFor="isActive" className="text-sm">Active</label>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={createUser.isPending || updateUser.isPending}>
              {createUser.isPending || updateUser.isPending
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
