import * as React from "react"
import { useLocations } from "@/features/locations/hooks"
import { useWarehouses } from "@/features/warehouses/hooks"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Badge } from "@/components/ui/Badge"
import { Skeleton } from "@/components/ui/skeleton"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { LocationFormDialog } from "@/features/locations/components/LocationFormDialog"
import { LocationDeleteDialog } from "@/features/locations/components/LocationDeleteDialog"
import type { Location } from "@/features/locations/types"
import { Plus, Search, Pencil, Trash2, MapPin, Layers } from "lucide-react"

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
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Locations</h1>
          <p className="text-muted-foreground mt-1">Manage storage locations within warehouses.</p>
        </div>
        <Button onClick={handleCreate} size="lg">
          <Plus className="h-4 w-4 mr-2" />
          Add Location
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Storage Locations</CardTitle>
              <CardDescription>
                Search and manage warehouse storage locations
              </CardDescription>
            </div>
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              <Layers className="h-4 w-4 text-primary" />
            </div>
          </div>
          <div className="flex items-center gap-4 mt-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by code, zone, or type..."
                className="pl-9 bg-muted/50 border-transparent focus:bg-background"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
          </div>
        </CardHeader>
        <CardContent>
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
                Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-10" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-16" /></TableCell>
                    <TableCell>
                      <div className="flex justify-end gap-2">
                        <Skeleton className="h-8 w-8" />
                        <Skeleton className="h-8 w-8" />
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              ) : filteredLocations.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="h-40">
                    <div className="flex flex-col items-center justify-center text-center">
                      <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                        <MapPin className="h-6 w-6 text-muted-foreground" />
                      </div>
                      <p className="text-sm font-medium text-foreground">No locations found</p>
                      <p className="text-xs text-muted-foreground mt-1">
                        {search ? "Try adjusting your search" : "Add a location to get started"}
                      </p>
                      {!search && (
                        <Button variant="outline" size="sm" className="mt-3" onClick={handleCreate}>
                          <Plus className="h-3 w-3 mr-1" />
                          Add Location
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                filteredLocations.map((location) => (
                  <TableRow key={location.id}>
                    <TableCell>
                      <span className="font-mono text-xs font-medium">{location.locationCode}</span>
                    </TableCell>
                    <TableCell className="font-medium text-foreground">
                      {getWarehouseName(location.warehouseId)}
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">{location.zone}</TableCell>
                    <TableCell>
                      <Badge variant="outline" className="text-xs">{location.locationType}</Badge>
                    </TableCell>
                    <TableCell>
                      {location.isPickFace ? (
                        <Badge variant="success" className="text-xs">Yes</Badge>
                      ) : (
                        <span className="text-muted-foreground text-sm">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <Badge variant={location.isActive ? "success" : "secondary"}>
                        {location.isActive ? "Active" : "Inactive"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleEdit(location)}
                          className="h-8 w-8"
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(location)}
                          className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
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
