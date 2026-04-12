import * as React from "react"
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
import { WarehouseFormDialog } from "@/features/warehouses/components/WarehouseFormDialog"
import { WarehouseDeleteDialog } from "@/features/warehouses/components/WarehouseDeleteDialog"
import type { Warehouse } from "@/features/warehouses/types"
import { Plus, Search, Pencil, Trash2, Warehouse as WarehouseIcon, Building2 } from "lucide-react"
import { formatDate } from "@/utils"

export default function WarehousesPage() {
  const [search, setSearch] = React.useState("")
  const [selectedWarehouse, setSelectedWarehouse] = React.useState<Warehouse | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = React.useState(false)

  const { data: warehousesData, isLoading } = useWarehouses({
    page: 1,
    limit: 10,
  })

  const handleEdit = (warehouse: Warehouse) => {
    setSelectedWarehouse(warehouse)
    setIsFormOpen(true)
  }

  const handleDelete = (warehouse: Warehouse) => {
    setSelectedWarehouse(warehouse)
    setIsDeleteOpen(true)
  }

  const handleCreate = () => {
    setSelectedWarehouse(null)
    setIsFormOpen(true)
  }

  const filteredWarehouses =
    warehousesData?.data?.filter(
      (warehouse) =>
        warehouse.name.toLowerCase().includes(search.toLowerCase()) ||
        warehouse.code.toLowerCase().includes(search.toLowerCase()) ||
        (warehouse.city?.toLowerCase() || "").includes(search.toLowerCase())
    ) || []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Warehouses</h1>
          <p className="text-muted-foreground mt-1">Manage warehouse locations and facilities.</p>
        </div>
        <Button onClick={handleCreate} size="lg">
          <Plus className="h-4 w-4 mr-2" />
          Add Warehouse
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle>Warehouse List</CardTitle>
              <CardDescription>
                Search and manage your warehouse facilities
              </CardDescription>
            </div>
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
              <Building2 className="h-4 w-4 text-primary" />
            </div>
          </div>
          <div className="flex items-center gap-4 mt-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search by name, code, or city..."
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
                <TableHead>Name</TableHead>
                <TableHead>Location</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Created</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-20" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-48" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-32" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-16" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                    <TableCell>
                      <div className="flex justify-end gap-2">
                        <Skeleton className="h-8 w-8" />
                        <Skeleton className="h-8 w-8" />
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              ) : filteredWarehouses.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-40">
                    <div className="flex flex-col items-center justify-center text-center">
                      <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                        <WarehouseIcon className="h-6 w-6 text-muted-foreground" />
                      </div>
                      <p className="text-sm font-medium text-foreground">No warehouses found</p>
                      <p className="text-xs text-muted-foreground mt-1">
                        {search ? "Try adjusting your search" : "Add a warehouse to get started"}
                      </p>
                      {!search && (
                        <Button variant="outline" size="sm" className="mt-3" onClick={handleCreate}>
                          <Plus className="h-3 w-3 mr-1" />
                          Add Warehouse
                        </Button>
                      )}
                    </div>
                  </TableCell>
                </TableRow>
              ) : (
                filteredWarehouses.map((warehouse) => (
                  <TableRow key={warehouse.id}>
                    <TableCell>
                      <span className="font-mono text-xs font-medium">{warehouse.code}</span>
                    </TableCell>
                    <TableCell>
                      <div>
                        <p className="font-medium text-foreground">{warehouse.name}</p>
                        <p className="text-xs text-muted-foreground">{warehouse.country}</p>
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {warehouse.city ? (
                        <span>
                          {warehouse.city}
                          {warehouse.state && `, ${warehouse.state}`}
                        </span>
                      ) : (
                        <span>-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <Badge variant={warehouse.isActive ? "success" : "secondary"}>
                        {warehouse.isActive ? "Active" : "Inactive"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {formatDate(warehouse.createdAt)}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex items-center justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleEdit(warehouse)}
                          className="h-8 w-8"
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => handleDelete(warehouse)}
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

      <WarehouseFormDialog
        warehouse={selectedWarehouse}
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
      />

      <WarehouseDeleteDialog
        warehouse={selectedWarehouse}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </div>
  )
}
