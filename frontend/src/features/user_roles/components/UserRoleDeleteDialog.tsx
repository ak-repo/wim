import * as React from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/Dialog"
import { Button } from "@/components/ui/Button"
import { useDeleteUserRole } from "@/features/user_roles/hooks"
import type { UserRole } from "@/features/user_roles/types"
import { AlertTriangle } from "lucide-react"

interface UserRoleDeleteDialogProps {
  userRole: UserRole | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onDeleted?: () => void
}

export const UserRoleDeleteDialog: React.FC<UserRoleDeleteDialogProps> = ({
  userRole,
  open,
  onOpenChange,
  onDeleted,
}) => {
  const deleteUserRole = useDeleteUserRole()

  const handleDelete = async () => {
    if (!userRole) return
    try {
      await deleteUserRole.mutateAsync(userRole.id)
      onOpenChange(false)
      onDeleted?.()
    } catch {
      // Error handled by mutation
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-destructive" />
            Delete User Role
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this user role? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        {userRole && (
          <div className="my-4 rounded-md bg-muted p-3">
            <p className="text-sm font-medium">{userRole.name}</p>
            <p className="text-xs text-muted-foreground">{userRole.description || "No description"}</p>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete} disabled={deleteUserRole.isPending}>
            {deleteUserRole.isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}