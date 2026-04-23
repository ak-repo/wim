import { apiService } from "@/lib/api"
import type {
  UserRole,
  CreateUserRoleRequest,
  UpdateUserRoleRequest,
  UserRoleParams,
} from "@/features/userRoles/types"
import type { PaginatedResponse } from "@/types"

export const userRoleService = {
  getRoles: async (params: UserRoleParams): Promise<PaginatedResponse<UserRole>> => {
    const response = await apiService.get<PaginatedResponse<UserRole>>(
      "/admin/user-roles",
      params as unknown as Record<string, unknown>
    )
    return response.data
  },

  getRole: async (id: number): Promise<UserRole> => {
    const response = await apiService.get<UserRole>(`/admin/user-roles/${id}`)
    return response.data
  },

  createRole: async (data: CreateUserRoleRequest): Promise<UserRole> => {
    const response = await apiService.post<UserRole>("/admin/user-roles", data)
    return response.data
  },

  updateRole: async (id: number, data: UpdateUserRoleRequest): Promise<UserRole> => {
    const response = await apiService.patch<UserRole>(`/admin/user-roles/${id}`, data)
    return response.data
  },

  deleteRole: async (id: number): Promise<void> => {
    await apiService.delete(`/admin/user-roles/${id}`)
  },
}
