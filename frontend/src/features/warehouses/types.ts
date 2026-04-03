import type { UUID } from "@/types"

export interface Warehouse {
  id: UUID
  code: string
  name: string
  addressLine1?: string
  addressLine2?: string
  city?: string
  state?: string
  postalCode?: string
  country: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateWarehouseRequest {
  code: string
  name: string
  addressLine1?: string
  addressLine2?: string
  city?: string
  state?: string
  postalCode?: string
  country: string
}

export interface UpdateWarehouseRequest {
  name?: string
  addressLine1?: string
  addressLine2?: string
  city?: string
  state?: string
  postalCode?: string
  country?: string
  isActive?: boolean
}

export interface WarehouseParams {
  active?: boolean
  page: number
  limit: number
}
