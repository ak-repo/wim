import * as React from "react"
import { useUserRoles } from "@/features/user_roles/hooks"
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
import { UserRoleFormDialog } from "@/features/user_roles/components/UserRoleFormDialog"
import { UserRoleDeleteDialog } from "@/features/user_roles/components/UserRoleDeleteDialog"
import type { UserRole } from "@/features/user_roles/types"
import { Plus, Search, Pencil, Trash2, Loader2, ChevronDown } from "lucide-react"
import { formatDate } from "@/utils"

export default function UserRolesPage() {
  const [search, setSearch] = React.useState("")
  const [openActionMenu, setOpenActionMenu] = React.useState<string | null>(null)
  const [selectedRole, setSelectedRole] = React.useState<UserRole | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = React.useState(false)

  const { data: rolesData, isLoading } = useUserRoles({
    page: 1,
    limit: 10,
  })

  const handleEdit = (role: UserRole) => {
    setSelectedRole(role)
    setIsFormOpen(true)
  }

  const handleDelete = (role: UserRole) => {
    setSelectedRole(role)
    setIsDeleteOpen(true)
  }

  const handleCreate = () => {
    setSelectedRole(null)
    setIsFormOpen(true)
  }

  React.useEffect(() => {
    const closeMenu = () => setOpenActionMenu(null)
    window.addEventListener("click", closeMenu)
    return () => window.removeEventListener("click", closeMenu)
  }, [])

  const filteredRoles =
    rolesData?.data?.filter((role) => {
      const matchesSearch =
        role.name.toLowerCase().includes(search.toLowerCase()) ||
        (role.description && role.description.toLowerCase().includes(search.toLowerCase()))
      return matchesSearch
    }) || []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">User Roles Management</h2>
          <p className="text-sm text-muted-foreground">Manage user roles and permissions.</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="h-4 w-4 mr-2" />
          Add Role
        </Button>
      </div>

      <Card>
        <CardHeader className="pb-4">
          <div className="relative">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search roles..."
              className="pl-9"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Description</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Created Date</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-24 text-center">
                    <Loader2 className="mx-auto h-5 w-5 animate-spin text-muted-foreground" />
                  </TableCell>
                </TableRow>
              ) : filteredRoles.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-24 text-center text-muted-foreground">
                    No roles found.
                  </TableCell>
                </TableRow>
              ) : (
                filteredRoles.map((role) => (
                  <TableRow key={role.id}>
                    <TableCell className="font-medium">{role.name}</TableCell>
                    <TableCell className="text-muted-foreground">
                      {role.description || "-"}
                    </TableCell>
                    <TableCell>
                      <Badge variant={role.isActive ? "success" : "destructive"}>
                        {role.isActive ? "active" : "inactive"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground">{formatDate(role.createdAt)}</TableCell>
                    <TableCell className="text-right">
                      <div className="relative inline-block text-left" onClick={(e) => e.stopPropagation()}>
                        <Button
                          variant="outline"
                          className="h-8 px-2"
                          onClick={() =>
                            setOpenActionMenu((prev) => (prev === String(role.id) ? null : String(role.id)))
                          }
                        >
                          Actions
                          <ChevronDown className="h-3 w-3" />
                        </Button>
                        {openActionMenu === String(role.id) && (
                          <div className="absolute right-0 z-20 mt-1 min-w-32 overflow-hidden rounded-md border border-border bg-card shadow-lg">
                            <button
                              className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-foreground hover:bg-muted/60"
                              onClick={() => {
                                setOpenActionMenu(null)
                                handleEdit(role)
                              }}
                            >
                              <Pencil className="h-3.5 w-3.5" />
                              Edit
                            </button>
                            <button
                              className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-destructive hover:bg-muted/60"
                              onClick={() => {
                                setOpenActionMenu(null)
                                handleDelete(role)
                              }}
                            >
                              <Trash2 className="h-3.5 w-3.5" />
                              Delete
                            </button>
                          </div>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <UserRoleFormDialog
        userRole={selectedRole}
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
      />

      <UserRoleDeleteDialog
        userRole={selectedRole}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </div>
  )
}
