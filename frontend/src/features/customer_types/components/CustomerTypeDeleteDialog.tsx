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
import { useDeleteCustomerType } from "@/features/customer_types/hooks"
import type { CustomerType } from "@/features/customer_types/types"
import { AlertTriangle } from "lucide-react"

interface CustomerTypeDeleteDialogProps {
  customerType: CustomerType | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onDeleted?: () => void
}

export const CustomerTypeDeleteDialog: React.FC<CustomerTypeDeleteDialogProps> = ({
  customerType,
  open,
  onOpenChange,
  onDeleted,
}) => {
  const deleteCustomerType = useDeleteCustomerType()

  const handleDelete = async () => {
    if (!customerType) return
    try {
      await deleteCustomerType.mutateAsync(customerType.id)
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
            Delete Customer Type
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this customer type? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        {customerType && (
          <div className="my-4 rounded-md bg-muted p-3">
            <p className="text-sm font-medium">{customerType.name}</p>
            <p className="text-xs text-muted-foreground">{customerType.description || "No description"}</p>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete} disabled={deleteCustomerType.isPending}>
            {deleteCustomerType.isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}