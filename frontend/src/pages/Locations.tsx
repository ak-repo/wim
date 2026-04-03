import * as React from "react"
import { useLocations } from "@/features/locations/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Badge } from "@/components/ui/Badge"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Card, CardContent, CardHeader } from "@/components/ui/Card"
import { LocationFormDialog } from "@/features/locations/components/LocationFormDialog"
import { LocationDeleteDialog } from "@/features/locations/components/LocationDeleteDialog"
import type { Location } from "@/features/locations/types"
import { Plus, Search, Pencil, Trash2, Loader2, MapPin } from "lucide-react"

export default function LocationsPage() {
  const [search, setSearch] = React.useState("")
  const [selectedLocation, setSelectedLocation] = React.useState<Location | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = React.useState(false)

  const { data: locationsData, isLoading } = useLocations({
    page: 1,
    limit: 10,
  })

  const { data: warehousesData } = useWarehouses({ page: 1, limit: 100 })

  const getWarehouseName = (id: string) => {
    return warehousesData?.data?.find((w) => w.id === id)?.name || "Unknown"
  }

  const handleEdit = (location: Location) => {
    setSelectedLocation(location)
    setIsFormOpen(true)
  }

  const handleDelete = (location: Location) => {
    setSelectedLocation(location)
    setIsDeleteOpen(true)
  }

  const handleCreate = () => {
    setSelectedLocation(null)
    setIsFormOpen(true)
  }

  const filteredLocations =
    locationsData?.data?.filter(
      (location) =>
        location.locationCode.toLowerCase().includes(search.toLowerCase()) ||
        location.zone.toLowerCase().includes(search.toLowerCase()) ||
        location.locationType.toLowerCase().includes(search.toLowerCase())
    ) || []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Locations</h2>
          <p className="text-muted-foreground">Manage storage locations within warehouses.</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="h-4 w-4 mr-2" />
          Add Location
        </Button>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center gap-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search locations..."
                className="pl-9"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Code</TableHead>
                  <TableHead>Warehouse</TableHead>
                  <TableHead>Zone</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Pick Face</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-24 text-center">
                      <Loader2 className="h-6 w-6 animate-spin mx-auto text-muted-foreground" />
                    </TableCell>
                  </TableRow>
                ) : filteredLocations.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-24 text-center">
                      <div className="flex flex-col items-center justify-center text-muted-foreground">
                        <MapPin className="h-8 w-8 mb-2" />
                        <p>No locations found.</p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredLocations.map((location) => (
                    <TableRow key={location.id}>
                      <TableCell className="font-medium">
                        {location.locationCode}
                      </TableCell>
                      <TableCell>
                        {getWarehouseName(location.warehouseId)}
                      </TableCell>
                      <TableCell>{location.zone}</TableCell>
                      <TableCell>
                        <Badge variant="outline">{location.locationType}</Badge>
                      </TableCell>
                      <TableCell>
                        {location.isPickFace ? (
                          <Badge variant="success">Yes</Badge>
                        ) : (
                          <span className="text-muted-foreground">-</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={location.isActive ? "success" : "secondary"}
                        >
                          {location.isActive ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleEdit(location)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(location)}
                            className="text-destructive hover:text-destructive"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <LocationFormDialog
        location={selectedLocation}
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
      />

      <LocationDeleteDialog
        location={selectedLocation}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </div>
  )
}
