import { apiService } from "@/lib/api"
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  User,
  UserRequest,
  UserParams,
} from "@/features/auth/types"
import type { PaginatedResponse } from "@/types"

export const authService = {
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await apiService.post<AuthResponse>("/adminPublic/login", data)
    return response.data
  },

  me: async (): Promise<User> => {
    const response = await apiService.get<User>("/admin/me")
    return response.data
  },

  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await apiService.post<AuthResponse>("/adminPublic/register", data)
    return response.data
  },

  logout: (): void => {
    localStorage.removeItem("accessToken")
    localStorage.removeItem("refreshToken")
  },
}

export const userService = {
  getUsers: async (params: UserParams): Promise<PaginatedResponse<User>> => {
    const response = await apiService.get<PaginatedResponse<User>>("/admin/users", params as unknown as Record<string, unknown>)
    return response.data
  },

  createUser: async (data: UserRequest): Promise<User> => {
    const response = await apiService.post<User>("/admin/users", data)
    return response.data
  },

  updateUser: async (id: number, data: Partial<UserRequest>): Promise<User> => {
    const response = await apiService.put<User>(`/admin/users/${id}`, data)
    return response.data
  },

  deleteUser: async (id: number): Promise<void> => {
    await apiService.delete(`/admin/users/${id}`)
  },
}
