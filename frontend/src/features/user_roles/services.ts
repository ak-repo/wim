import { apiService } from "@/lib/api"
import type {
  UserRole,
  CreateUserRoleRequest,
  UpdateUserRoleRequest,
  UserRoleParams,
} from "@/features/user_roles/types"
import type { PaginatedResponse } from "@/types"

interface UserRoleListResponse {
  data: UserRole[]
  total_count: number
  total_page: number
  current_page: number
  limit: number
}

interface CreateUserRoleResponse {
  id: number
}

interface UpdateUserRoleResponse {
  message: string
}

interface DeleteUserRoleResponse {
  message: string
}

const mapUserRoleListResponse = (response: UserRoleListResponse): PaginatedResponse<UserRole> => {
  return {
    data: response.data,
    total: response.total_count,
    page: response.current_page,
    limit: response.limit,
    totalPages: response.total_page,
  }
}

export const userRoleService = {
  getUserRoles: async (params: UserRoleParams): Promise<PaginatedResponse<UserRole>> => {
    const response = await apiService.get<UserRoleListResponse>(
      "/admin/user-roles",
      params as unknown as Record<string, unknown>
    )
    return mapUserRoleListResponse(response.data)
  },

  getUserRole: async (id: number): Promise<UserRole> => {
    const response = await apiService.get<UserRole>(`/admin/user-roles/${id}`)
    return response.data
  },

  createUserRole: async (data: CreateUserRoleRequest): Promise<CreateUserRoleResponse> => {
    const response = await apiService.post<CreateUserRoleResponse>("/admin/user-roles", data)
    return response.data
  },

  updateUserRole: async (id: number, data: UpdateUserRoleRequest): Promise<UpdateUserRoleResponse> => {
    const response = await apiService.put<UpdateUserRoleResponse>(`/admin/user-roles/${id}`, data)
    return response.data
  },

  deleteUserRole: async (id: number): Promise<DeleteUserRoleResponse> => {
    const response = await apiService.delete<DeleteUserRoleResponse>(`/admin/user-roles/${id}`)
    return response.data
  },
}