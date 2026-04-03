import type { UUID } from "@/types"

export interface Location {
  id: UUID
  warehouseId: UUID
  zone: string
  aisle?: string
  rack?: string
  bin?: string
  locationCode: string
  locationType: string
  isPickFace: boolean
  maxWeight?: number
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateLocationRequest {
  warehouseId: UUID
  zone: string
  aisle?: string
  rack?: string
  bin?: string
  locationCode: string
  locationType: string
  isPickFace: boolean
  maxWeight?: number
}

export interface UpdateLocationRequest {
  zone?: string
  aisle?: string
  rack?: string
  bin?: string
  locationCode?: string
  locationType?: string
  isPickFace?: boolean
  maxWeight?: number
  isActive?: boolean
}

export interface LocationParams {
  active?: boolean
  warehouseId?: UUID
  zone?: string
  page: number
  limit: number
}
