export interface UserRole {
  id: number
  name: string
  refCode: string
  isActive: boolean
}

export interface CreateUserRoleRequest {
  name: string
  isActive?: boolean
}

export interface UpdateUserRoleRequest {
  name?: string
  isActive?: boolean
}

export interface UserRoleParams {
  active?: boolean
  page: number
  limit: number
}
