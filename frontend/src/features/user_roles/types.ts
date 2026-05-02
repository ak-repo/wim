export interface UserRole {
  id: number
  name: string
  description?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateUserRoleRequest {
  name: string
  description?: string
  isActive?: boolean
}

export interface UpdateUserRoleRequest {
  name?: string
  description?: string
  isActive?: boolean
}

export interface UserRoleParams {
  active?: boolean
  page: number
  limit: number
}