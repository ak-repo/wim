import * as React from "react"
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
import { WarehouseFormDialog } from "@/features/warehouses/components/WarehouseFormDialog"
import { WarehouseDeleteDialog } from "@/features/warehouses/components/WarehouseDeleteDialog"
import type { Warehouse } from "@/features/warehouses/types"
import { Plus, Search, Pencil, Trash2, Loader2, Warehouse as WarehouseIcon } from "lucide-react"
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
          <h2 className="text-2xl font-bold tracking-tight">Warehouses</h2>
          <p className="text-muted-foreground">Manage warehouse locations.</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="h-4 w-4 mr-2" />
          Add Warehouse
        </Button>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <div className="flex items-center gap-4">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search warehouses..."
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
                  <TableHead>Name</TableHead>
                  <TableHead>Location</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center">
                      <Loader2 className="h-6 w-6 animate-spin mx-auto text-muted-foreground" />
                    </TableCell>
                  </TableRow>
                ) : filteredWarehouses.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center">
                      <div className="flex flex-col items-center justify-center text-muted-foreground">
                        <WarehouseIcon className="h-8 w-8 mb-2" />
                        <p>No warehouses found.</p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredWarehouses.map((warehouse) => (
                    <TableRow key={warehouse.id}>
                      <TableCell className="font-medium">{warehouse.code}</TableCell>
                      <TableCell>
                        <div>
                          <p className="font-medium">{warehouse.name}</p>
                          <p className="text-xs text-muted-foreground">{warehouse.country}</p>
                        </div>
                      </TableCell>
                      <TableCell>
                        {warehouse.city ? (
                          <div className="text-sm">
                            {warehouse.city}
                            {warehouse.state && `, ${warehouse.state}`}
                          </div>
                        ) : (
                          <span className="text-muted-foreground">-</span>
                        )}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={warehouse.isActive ? "success" : "secondary"}
                        >
                          {warehouse.isActive ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell>{formatDate(warehouse.createdAt)}</TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleEdit(warehouse)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(warehouse)}
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
