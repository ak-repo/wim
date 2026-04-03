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
import { useCreateLocation, useUpdateLocation } from "@/features/locations/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import type { Location, CreateLocationRequest, UpdateLocationRequest } from "@/features/locations/types"

interface LocationFormDialogProps {
  location: Location | null
  open: boolean
  onOpenChange: (open: boolean) => void
}

export const LocationFormDialog: React.FC<LocationFormDialogProps> = ({
  location,
  open,
  onOpenChange,
}) => {
  const createLocation = useCreateLocation()
  const updateLocation = useUpdateLocation()
  const { data: warehousesData } = useWarehouses({ page: 1, limit: 100 })
  const isEditing = !!location

  const [formData, setFormData] = React.useState({
    warehouseId: location?.warehouseId || "",
    zone: location?.zone || "",
    aisle: location?.aisle || "",
    rack: location?.rack || "",
    bin: location?.bin || "",
    locationCode: location?.locationCode || "",
    locationType: location?.locationType || "storage",
    isPickFace: location?.isPickFace ?? false,
    maxWeight: location?.maxWeight?.toString() || "",
  })

  React.useEffect(() => {
    setFormData({
      warehouseId: location?.warehouseId || "",
      zone: location?.zone || "",
      aisle: location?.aisle || "",
      rack: location?.rack || "",
      bin: location?.bin || "",
      locationCode: location?.locationCode || "",
      locationType: location?.locationType || "storage",
      isPickFace: location?.isPickFace ?? false,
      maxWeight: location?.maxWeight?.toString() || "",
    })
  }, [location, open])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }))
  }

  const handleCheckboxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.checked }))
  }

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (isEditing && location) {
        const updateData: UpdateLocationRequest = {
          zone: formData.zone,
          aisle: formData.aisle || undefined,
          rack: formData.rack || undefined,
          bin: formData.bin || undefined,
          locationCode: formData.locationCode,
          locationType: formData.locationType,
          isPickFace: formData.isPickFace,
          maxWeight: formData.maxWeight ? parseFloat(formData.maxWeight) : undefined,
        }
        await updateLocation.mutateAsync({ id: location.id, data: updateData })
      } else {
        const createData: CreateLocationRequest = {
          warehouseId: formData.warehouseId,
          zone: formData.zone,
          aisle: formData.aisle || undefined,
          rack: formData.rack || undefined,
          bin: formData.bin || undefined,
          locationCode: formData.locationCode,
          locationType: formData.locationType,
          isPickFace: formData.isPickFace,
          maxWeight: formData.maxWeight ? parseFloat(formData.maxWeight) : undefined,
        }
        await createLocation.mutateAsync(createData)
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
          <DialogTitle>{isEditing ? "Edit Location" : "Create Location"}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? "Update storage location details."
              : "Add a new storage location."}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={onSubmit} className="space-y-4 mt-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Warehouse *</label>
            <Select
              name="warehouseId"
              value={formData.warehouseId}
              onChange={handleChange}
              disabled={isEditing}
            >
              <option value="">Select warehouse</option>
              {warehousesData?.data?.map((w) => (
                <option key={w.id} value={w.id}>
                  {w.name} ({w.code})
                </option>
              ))}
            </Select>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Location Code *</label>
              <Input
                name="locationCode"
                value={formData.locationCode}
                onChange={handleChange}
                placeholder="LOC-001"
                disabled={isEditing}
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Location Type *</label>
              <Select
                name="locationType"
                value={formData.locationType}
                onChange={handleChange}
              >
                <option value="storage">Storage</option>
                <option value="receiving">Receiving</option>
                <option value="shipping">Shipping</option>
                <option value="quarantine">Quarantine</option>
                <option value="staging">Staging</option>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Zone *</label>
            <Input
              name="zone"
              value={formData.zone}
              onChange={handleChange}
              placeholder="Zone A"
            />
          </div>

          <div className="grid grid-cols-3 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Aisle</label>
              <Input
                name="aisle"
                value={formData.aisle}
                onChange={handleChange}
                placeholder="A1"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Rack</label>
              <Input
                name="rack"
                value={formData.rack}
                onChange={handleChange}
                placeholder="R1"
              />
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Bin</label>
              <Input
                name="bin"
                value={formData.bin}
                onChange={handleChange}
                placeholder="B1"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Max Weight (kg)</label>
              <Input
                name="maxWeight"
                type="number"
                step="0.1"
                value={formData.maxWeight}
                onChange={handleChange}
                placeholder="0.0"
              />
            </div>

            <div className="flex items-center gap-2 pt-6">
              <input
                type="checkbox"
                name="isPickFace"
                id="isPickFace"
                checked={formData.isPickFace}
                onChange={handleCheckboxChange}
                className="h-4 w-4 rounded border-input"
              />
              <label htmlFor="isPickFace" className="text-sm">Pick Face</label>
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={createLocation.isPending || updateLocation.isPending}
            >
              {createLocation.isPending || updateLocation.isPending
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
