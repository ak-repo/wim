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
import { useDeleteCustomer } from "@/features/customers/hooks"
import type { Customer } from "@/features/customers/types"
import { AlertTriangle } from "lucide-react"

interface CustomerDeleteDialogProps {
  customer: Customer | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onDeleted?: () => void
}

export const CustomerDeleteDialog: React.FC<CustomerDeleteDialogProps> = ({
  customer,
  open,
  onOpenChange,
  onDeleted,
}) => {
  const deleteCustomer = useDeleteCustomer()

  const handleDelete = async () => {
    if (!customer) return
    try {
      await deleteCustomer.mutateAsync(customer.id)
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
            Delete Customer
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this customer? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        {customer && (
          <div className="my-4 rounded-md bg-muted p-3">
            <p className="text-sm font-medium">{customer.name}</p>
            <p className="text-xs text-muted-foreground">{customer.email}</p>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={handleDelete} disabled={deleteCustomer.isPending}>
            {deleteCustomer.isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
