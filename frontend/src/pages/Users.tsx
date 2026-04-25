import * as React from "react"
import { useUsers } from "@/features/auth/hooks"
import { useUserRoles } from "@/features/userRoles/hooks"
import { UserFormDialog } from "@/features/users/components/UserFormDialog"
import { UserDeleteDialog } from "@/features/users/components/UserDeleteDialog"
import { UserRoleFormDialog } from "@/features/userRoles/components/UserRoleFormDialog"
import { UserRoleDeleteDialog } from "@/features/userRoles/components/UserRoleDeleteDialog"
import type { User } from "@/features/auth/types"
import type { UserRole } from "@/features/userRoles/types"
import { Plus, Search, Pencil, Trash2, Users, BadgeCheck } from "lucide-react"

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

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-[20px] font-medium tracking-tight text-ink">User Master</h1>
          <p className="text-[12px] text-ink-3 mt-0.5">Manage users and roles in one place.</p>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex bg-white rounded-[7px] border-[0.5px] border-border-default p-1">
            <button
              onClick={() => setActiveTab("users")}
              className={`px-3 py-1.5 text-[12px] font-medium rounded-[5px] transition-colors flex items-center gap-1.5 ${
                activeTab === "users" ? "bg-surface-2 text-ink" : "text-ink-3 hover:text-ink-2"
              }`}
            >
              <Users className="h-3.5 w-3.5" />
              Users
            </button>
            <button
              onClick={() => setActiveTab("roles")}
              className={`px-3 py-1.5 text-[12px] font-medium rounded-[5px] transition-colors flex items-center gap-1.5 ${
                activeTab === "roles" ? "bg-surface-2 text-ink" : "text-ink-3 hover:text-ink-2"
              }`}
            >
              <BadgeCheck className="h-3.5 w-3.5" />
              Roles
            </button>
          </div>
          <button
            onClick={activeTab === "users" ? handleCreate : handleCreateRole}
            className="flex items-center gap-2 bg-ink text-white rounded-[7px] px-3 py-1.5 text-[12px] font-medium hover:bg-ink-2 transition-colors"
          >
            <Plus className="h-3.5 w-3.5" />
            {activeTab === "users" ? "Add User" : "Add Role"}
          </button>
        </div>
      </div>

      <div className="bg-white border-[0.5px] border-border-default rounded-[10px] overflow-hidden flex flex-col">
        <div className="p-[14px_16px] border-b-[0.5px] border-border-default flex items-center justify-between gap-4">
          <div className="relative flex-1 max-w-[240px]">
            <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-ink-3" />
            <input
              type="text"
              placeholder={activeTab === "users" ? "Search users..." : "Search roles..."}
              className="h-[30px] w-full bg-surface-2 border-[0.5px] border-border-default rounded-[7px] pl-8 pr-3 text-[12px] text-ink placeholder:text-ink-3 focus:outline-none focus:border-border-2"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
        </div>
        
        <div className="overflow-x-auto">
          <table className="w-full text-left text-[13px]">
            <thead>
              <tr className="border-b-[0.5px] border-border-default">
                <th className="font-medium text-ink-3 px-4 py-3 whitespace-nowrap">
                  {activeTab === "users" ? "User" : "Name"}
                </th>
                <th className="font-medium text-ink-3 px-4 py-3 whitespace-nowrap">
                  {activeTab === "users" ? "Role" : "Ref Code"}
                </th>
                <th className="font-medium text-ink-3 px-4 py-3 whitespace-nowrap">Status</th>
                {activeTab === "users" && (
                  <th className="font-medium text-ink-3 px-4 py-3 whitespace-nowrap">Created</th>
                )}
                <th className="font-medium text-ink-3 px-4 py-3 text-right whitespace-nowrap">Actions</th>
              </tr>
            </thead>
            <tbody>
              {activeTab === "users" ? (
                isLoading ? (
                  <tr>
                    <td colSpan={5} className="px-4 py-8 text-center text-[12px] text-ink-3">Loading users...</td>
                  </tr>
                ) : filteredUsers.length === 0 ? (
                  <tr>
                    <td colSpan={5} className="px-4 py-8 text-center">
                      <div className="flex flex-col items-center justify-center">
                        <div className="flex h-[36px] w-[36px] items-center justify-center rounded-[10px] bg-surface-2 mb-3">
                          <Users className="h-5 w-5 text-ink-3" />
                        </div>
                        <p className="text-[12px] text-ink-3">No users found</p>
                      </div>
                    </td>
                  </tr>
                ) : (
                  filteredUsers.map((user) => (
                    <tr key={user.id} className="border-b-[0.5px] border-border-default last:border-0 hover:bg-surface-2/50 transition-colors">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-3">
                          <div className="h-[32px] w-[32px] rounded-full bg-accent-bg flex shrink-0 items-center justify-center">
                            <span className="text-[12px] font-medium text-accent-green uppercase">
                              {user.username.charAt(0)}
                            </span>
                          </div>
                          <div className="flex flex-col">
                            <span className="font-medium text-ink">{user.username}</span>
                            <span className="text-[11px] text-ink-3">{user.email}</span>
                          </div>
                        </div>
                      </td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium ${
                          user.role === 'super_admin' ? 'bg-accent-bg text-accent-green' : 'bg-surface-2 text-ink-2'
                        }`}>
                          {user.role.replace("_", " ")}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium ${
                          user.isActive ? 'bg-accent-bg text-accent-green' : 'bg-surface-2 text-ink-3'
                        }`}>
                          {user.isActive ? "Active" : "Inactive"}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-[12px] text-ink-3">
                        {new Date(user.created_at).toLocaleDateString()}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-3">
                          <button
                            onClick={() => handleEdit(user)}
                            className="text-ink-3 hover:text-ink transition-colors"
                          >
                            <Pencil className="h-[15px] w-[15px]" />
                          </button>
                          <button
                            onClick={() => handleDelete(user)}
                            className="text-ink-3 hover:text-coral transition-colors"
                          >
                            <Trash2 className="h-[15px] w-[15px]" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )
              ) : (
                rolesLoading ? (
                  <tr>
                    <td colSpan={4} className="px-4 py-8 text-center text-[12px] text-ink-3">Loading roles...</td>
                  </tr>
                ) : filteredRoles.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="px-4 py-8 text-center">
                      <div className="flex flex-col items-center justify-center">
                        <div className="flex h-[36px] w-[36px] items-center justify-center rounded-[10px] bg-surface-2 mb-3">
                          <BadgeCheck className="h-5 w-5 text-ink-3" />
                        </div>
                        <p className="text-[12px] text-ink-3">No roles found</p>
                      </div>
                    </td>
                  </tr>
                ) : (
                  filteredRoles.map((role) => (
                    <tr key={role.id} className="border-b-[0.5px] border-border-default last:border-0 hover:bg-surface-2/50 transition-colors">
                      <td className="px-4 py-3 font-medium text-ink">{role.name}</td>
                      <td className="px-4 py-3">
                        <span className="font-mono text-[11px] text-ink-3 bg-surface-2 px-1.5 py-0.5 rounded">
                          {role.refCode}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium ${
                          role.isActive ? 'bg-accent-bg text-accent-green' : 'bg-surface-2 text-ink-3'
                        }`}>
                          {role.isActive ? "Active" : "Inactive"}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-3">
                          <button
                            onClick={() => handleEditRole(role)}
                            className="text-ink-3 hover:text-ink transition-colors"
                          >
                            <Pencil className="h-[15px] w-[15px]" />
                          </button>
                          <button
                            onClick={() => handleDeleteRole(role)}
                            className="text-ink-3 hover:text-coral transition-colors"
                          >
                            <Trash2 className="h-[15px] w-[15px]" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )
              )}
            </tbody>
          </table>
        </div>
      </div>

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
