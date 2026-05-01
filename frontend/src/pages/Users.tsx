import * as React from "react"
import { useUsers } from "@/features/auth/hooks"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Select } from "@/components/ui/Select"
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
import { UserFormDialog } from "@/features/users/components/UserFormDialog"
import { UserDeleteDialog } from "@/features/users/components/UserDeleteDialog"
import type { User } from "@/features/auth/types"
import { Plus, Search, Pencil, Trash2, Loader2, ChevronDown } from "lucide-react"
import { formatDate } from "@/utils"

export default function UsersPage() {
  const [search, setSearch] = React.useState("")
  const [roleFilter, setRoleFilter] = React.useState("all")
  const [openActionMenu, setOpenActionMenu] = React.useState<string | null>(null)
  const [selectedUser, setSelectedUser] = React.useState<User | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = React.useState(false)

  const { data: usersData, isLoading } = useUsers({
    page: 1,
    limit: 10,
  })

  const handleEdit = (user: User) => {
    setSelectedUser(user)
    setIsFormOpen(true)
  }

  const handleDelete = (user: User) => {
    setSelectedUser(user)
    setIsDeleteOpen(true)
  }

  const handleCreate = () => {
    setSelectedUser(null)
    setIsFormOpen(true)
  }

  React.useEffect(() => {
    const closeMenu = () => setOpenActionMenu(null)
    window.addEventListener("click", closeMenu)
    return () => window.removeEventListener("click", closeMenu)
  }, [])

  const filteredUsers =
    usersData?.data?.filter((user) => {
      const matchesSearch =
        user.username.toLowerCase().includes(search.toLowerCase()) ||
        user.email.toLowerCase().includes(search.toLowerCase())

      const matchesRole = roleFilter === "all" || user.role === roleFilter

      return matchesSearch && matchesRole
    }) || []

  const getRoleBadgeVariant = (role: string) => {
    switch (role) {
      case "super_admin":
        return "warning"
      case "admin":
        return "default"
      case "manager":
      case "worker":
        return "secondary"
      default:
        return "secondary"
    }
  }

  const getRoleLabel = (role: string) => {
    return role === "admin" || role === "super_admin" ? "admin" : "user"
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">Users Management</h2>
          <p className="text-sm text-muted-foreground">Manage accounts, roles, and account status.</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="h-4 w-4 mr-2" />
          Add User
        </Button>
      </div>

      <Card>
        <CardHeader className="pb-4">
          <div className="grid gap-3 md:grid-cols-[minmax(220px,1fr)_200px]">
            <div className="relative">
              <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search users..."
                className="pl-9"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
            <Select value={roleFilter} onChange={(e) => setRoleFilter(e.target.value)}>
              <option value="all">All Roles</option>
              <option value="super_admin">Super Admin</option>
              <option value="admin">Admin</option>
              <option value="worker">User</option>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Role</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Created Date</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-24 text-center">
                    <Loader2 className="mx-auto h-5 w-5 animate-spin text-muted-foreground" />
                  </TableCell>
                </TableRow>
              ) : filteredUsers.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-24 text-center text-muted-foreground">
                    No users found.
                  </TableCell>
                </TableRow>
              ) : (
                filteredUsers.map((user) => (
                  <TableRow key={user.id}>
                    <TableCell className="font-medium">{user.username}</TableCell>
                    <TableCell className="text-muted-foreground">{user.email}</TableCell>
                    <TableCell>
                      <Badge variant={getRoleBadgeVariant(user.role)}>
                        {getRoleLabel(user.role)}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant={user.isActive ? "success" : "destructive"}>
                        {user.isActive ? "active" : "inactive"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground">{formatDate(user.created_at)}</TableCell>
                    <TableCell className="text-right">
                      <div className="relative inline-block text-left" onClick={(e) => e.stopPropagation()}>
                        <Button
                          variant="outline"
                          className="h-8 px-2"
                          onClick={() =>
                            setOpenActionMenu((prev) => (prev === user.id ? null : user.id))
                          }
                        >
                          Actions
                          <ChevronDown className="h-3 w-3" />
                        </Button>
                        {openActionMenu === user.id && (
                          <div className="absolute right-0 z-20 mt-1 min-w-32 overflow-hidden rounded-md border border-border bg-card shadow-lg">
                            <button
                              className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-foreground hover:bg-muted/60"
                              onClick={() => {
                                setOpenActionMenu(null)
                                handleEdit(user)
                              }}
                            >
                              <Pencil className="h-3.5 w-3.5" />
                              Edit
                            </button>
                            <button
                              className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-destructive hover:bg-muted/60"
                              onClick={() => {
                                setOpenActionMenu(null)
                                handleDelete(user)
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

      <UserFormDialog
        user={selectedUser}
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
      />

      <UserDeleteDialog
        user={selectedUser}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />
    </div>
  )
}
