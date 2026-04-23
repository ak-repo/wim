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
import { useDeleteUserRole } from "@/features/userRoles/hooks"
import type { UserRole } from "@/features/userRoles/types"
import { AlertTriangle } from "lucide-react"

interface UserRoleDeleteDialogProps {
  role: UserRole | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const UserRoleDeleteDialog: React.FC<UserRoleDeleteDialogProps> = ({
  role,
  open,
  onOpenChange,
}) => {
  const deleteRole = useDeleteUserRole()

  const handleDelete = async () => {
    if (!role) return
    try {
      await deleteRole.mutateAsync(role.id)
      onOpenChange(false)
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
            Delete Role
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this role? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        {role && (
          <div className="my-4 p-3 bg-muted rounded-md">
            <p className="text-sm font-medium">{role.name}</p>
            <p className="text-xs text-muted-foreground">{role.refCode}</p>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteRole.isPending}
          >
            {deleteRole.isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
