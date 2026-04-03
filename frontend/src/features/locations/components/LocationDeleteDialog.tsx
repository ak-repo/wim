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
import { useDeleteLocation } from "@/features/locations/hooks"
import type { Location } from "@/features/locations/types"
import { AlertTriangle } from "lucide-react"

interface LocationDeleteDialogProps {
  location: Location | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const LocationDeleteDialog: React.FC<LocationDeleteDialogProps> = ({
  location,
  open,
  onOpenChange,
}) => {
  const deleteLocation = useDeleteLocation()

  const handleDelete = async () => {
    if (!location) return
    try {
      await deleteLocation.mutateAsync(location.id)
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
            Delete Location
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to delete this location? This action cannot be undone.
          </DialogDescription>
        </DialogHeader>

        {location && (
          <div className="my-4 p-3 bg-muted rounded-md">
            <p className="text-sm font-medium">{location.locationCode}</p>
            <p className="text-xs text-muted-foreground">Zone: {location.zone}</p>
          </div>
        )}

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            variant="destructive"
            onClick={handleDelete}
            disabled={deleteLocation.isPending}
          >
            {deleteLocation.isPending ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
