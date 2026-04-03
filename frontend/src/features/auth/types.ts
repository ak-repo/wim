import type { UUID } from "@/types"

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  role?: string
  contact?: string
}

export interface AuthResponse {
  accessToken: string
  refreshToken?: string
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface User {
  id: UUID
  username: string
  email: string
  role: string
  contact?: string
  isActive: boolean
  created_at: string
  updated_at: string
}

export interface UserRequest {
  id?: UUID
  username: string
  email: string
  password?: string
  role: string
  contact?: string
  isActive: boolean
}

export interface UserParams {
  active?: boolean
  page: number
  limit: number
}
