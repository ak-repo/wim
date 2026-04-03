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
import { useDeleteWarehouse } from "@/features/warehouses/hooks"
import type { Warehouse } from "@/features/warehouses/types"
import { AlertTriangle } from "lucide-react"

interface WarehouseDeleteDialogProps {
  warehouse: Warehouse | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const WarehouseDeleteDialog: React.FC<WarehouseDeleteDialogProps> = ({
  warehouse,
  open,
  onOpenChange,
}) => {
  const deleteWarehouse = useDeleteWarehouse()

  const handleDelete = async () => {
    if (!warehouse) return
    try {
      await deleteWarehouse.mutateAsync(warehouse.id)
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
            Delete Warehouse
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this warehouse? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        {warehouse && (
          <div className="my-4 p-3 bg-muted rounded-md">
            <p className="text-sm font-medium">{warehouse.name}</p>
            <p className="text-xs text-muted-foreground">Code: {warehouse.code}</p>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteWarehouse.isPending}
          >
            {deleteWarehouse.isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
