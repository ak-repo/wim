import * as React from "react"
import { useUsers } from "@/features/auth/hooks"
import { useUserRoles } from "@/features/userRoles/hooks"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Badge } from "@/components/ui/Badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Tabs } from "@/components/ui/Tabs"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/Card"
import { UserFormDialog } from "@/features/users/components/UserFormDialog"
import { UserDeleteDialog } from "@/features/users/components/UserDeleteDialog"
import { UserRoleFormDialog } from "@/features/userRoles/components/UserRoleFormDialog"
import { UserRoleDeleteDialog } from "@/features/userRoles/components/UserRoleDeleteDialog"
import type { User } from "@/features/auth/types"
import type { UserRole } from "@/features/userRoles/types"
import { Plus, Search, Pencil, Trash2, Users, Shield, BadgeCheck } from "lucide-react"
import { formatDate } from "@/utils"

export default function UsersPage() {
  const [search, setSearch] = React.useState("")
  const [selectedUser, setSelectedUser] = React.useState<User | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = React.useState(false)
  const [activeTab, setActiveTab] = React.useState("users")
  const [selectedRole, setSelectedRole] = React.useState<UserRole | null>(null)
  const [isRoleFormOpen, setIsRoleFormOpen] = React.useState(false)
  const [isRoleDeleteOpen, setIsRoleDeleteOpen] = React.useState(false)

  const { data: usersData, isLoading } = useUsers({
    page: 1,
    limit: 10,
  })

  const { data: rolesData, isLoading: rolesLoading } = useUserRoles({
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

  const handleEditRole = (role: UserRole) => {
    setSelectedRole(role)
    setIsRoleFormOpen(true)
  }

  const handleDeleteRole = (role: UserRole) => {
    setSelectedRole(role)
    setIsRoleDeleteOpen(true)
  }

  const handleCreateRole = () => {
    setSelectedRole(null)
    setIsRoleFormOpen(true)
  }

  const filteredUsers =
    usersData?.data?.filter(
      (user) =>
        user.username.toLowerCase().includes(search.toLowerCase()) ||
        user.email.toLowerCase().includes(search.toLowerCase())
    ) || []

  const filteredRoles =
    rolesData?.data?.filter((role) =>
      role.name.toLowerCase().includes(search.toLowerCase())
    ) || []

  const getRoleBadgeVariant = (role: string) => {
    switch (role) {
      case "super_admin":
        return "destructive"
      case "admin":
        return "default"
      case "manager":
        return "secondary"
      default:
        return "outline"
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">User Master</h1>
          <p className="text-muted-foreground mt-1">Manage users and roles in one place.</p>
        </div>
        <div className="flex items-center gap-3">
          <Tabs
            value={activeTab}
            onChange={setActiveTab}
            tabs={[
              { id: "users", label: "Users", icon: Users },
              { id: "roles", label: "Roles", icon: BadgeCheck },
            ]}
          />
          <Button
            onClick={activeTab === "users" ? handleCreate : handleCreateRole}
            size="lg"
          >
            <Plus className="h-4 w-4 mr-2" />
            {activeTab === "users" ? "Add User" : "Add Role"}
          </Button>
        </div>
      </div>

      {activeTab === "users" ? (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>User Management</CardTitle>
                <CardDescription>
                  Search and manage user accounts
                </CardDescription>
              </div>
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                <Shield className="h-4 w-4 text-primary" />
              </div>
            </div>
            <div className="flex items-center gap-4 mt-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search by name or email..."
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
                  <TableHead>User</TableHead>
                  <TableHead>Role</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell><Skeleton className="h-4 w-40" /></TableCell>
                      <TableCell><Skeleton className="h-6 w-20" /></TableCell>
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
                ) : filteredUsers.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} className="h-40">
                      <div className="flex flex-col items-center justify-center text-center">
                        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                          <Users className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <p className="text-sm font-medium text-foreground">No users found</p>
                        <p className="text-xs text-muted-foreground mt-1">
                          {search ? "Try adjusting your search" : "Add a user to get started"}
                        </p>
                        {!search && (
                          <Button variant="outline" size="sm" className="mt-3" onClick={handleCreate}>
                            <Plus className="h-3 w-3 mr-1" />
                            Add User
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredUsers.map((user) => (
                    <TableRow key={user.id}>
                      <TableCell>
                        <div className="flex items-center gap-3">
                          <div className="h-8 w-8 rounded-full bg-primary/10 flex items-center justify-center">
                            <span className="text-xs font-semibold text-primary uppercase">
                              {user.username.charAt(0)}
                            </span>
                          </div>
                          <div>
                            <p className="font-medium text-foreground">{user.username}</p>
                            <p className="text-xs text-muted-foreground">{user.email}</p>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant={getRoleBadgeVariant(user.role)}>
                          {user.role.replace("_", " ")}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <Badge variant={user.isActive ? "success" : "secondary"}>
                          {user.isActive ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-muted-foreground text-sm">
                        {formatDate(user.created_at)}
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleEdit(user)}
                            className="h-8 w-8"
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(user)}
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
      ) : (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>User Roles</CardTitle>
                <CardDescription>
                  Manage role definitions for access control
                </CardDescription>
              </div>
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                <BadgeCheck className="h-4 w-4 text-primary" />
              </div>
            </div>
            <div className="flex items-center gap-4 mt-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search by role name..."
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
                  <TableHead>Name</TableHead>
                  <TableHead>Ref Code</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rolesLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell><Skeleton className="h-4 w-40" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                      <TableCell><Skeleton className="h-6 w-16" /></TableCell>
                      <TableCell>
                        <div className="flex justify-end gap-2">
                          <Skeleton className="h-8 w-8" />
                          <Skeleton className="h-8 w-8" />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                ) : filteredRoles.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} className="h-40">
                      <div className="flex flex-col items-center justify-center text-center">
                        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                          <BadgeCheck className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <p className="text-sm font-medium text-foreground">No roles found</p>
                        <p className="text-xs text-muted-foreground mt-1">
                          {search ? "Try adjusting your search" : "Add a role to get started"}
                        </p>
                        {!search && (
                          <Button variant="outline" size="sm" className="mt-3" onClick={handleCreateRole}>
                            <Plus className="h-3 w-3 mr-1" />
                            Add Role
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredRoles.map((role) => (
                    <TableRow key={role.id}>
                      <TableCell>
                        <p className="font-medium text-foreground">{role.name}</p>
                      </TableCell>
                      <TableCell>
                        <span className="font-mono text-xs font-medium text-muted-foreground">
                          {role.refCode}
                        </span>
                      </TableCell>
                      <TableCell>
                        <Badge variant={role.isActive ? "success" : "secondary"}>
                          {role.isActive ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleEditRole(role)}
                            className="h-8 w-8"
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDeleteRole(role)}
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
      )}

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

      <UserRoleFormDialog
        role={selectedRole}
        open={isRoleFormOpen}
        onOpenChange={setIsRoleFormOpen}
      />

      <UserRoleDeleteDialog
        role={selectedRole}
        open={isRoleDeleteOpen}
        onOpenChange={setIsRoleDeleteOpen}
      />
    </div>
  )
}
