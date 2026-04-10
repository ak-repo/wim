export interface Inventory {
  id: number
  productId: number
  warehouseId: number
  locationId: number
  batchId?: number
  quantity: number
  reservedQty: number
  availableQty: number
  version: number
  createdAt: string
  updatedAt: string
}

export interface StockMovement {
  id: number
  movementType: string
  productId: number
  warehouseId: number
  locationIdFrom?: number
  locationIdTo?: number
  batchId?: number
  quantity: number
  referenceType?: string
  referenceId?: number
  performedBy?: number
  notes?: string
  createdAt: string
}

export interface AdjustInventoryRequest {
  productId: number
  warehouseId: number
  locationId: number
  batchId?: number
  quantity: number
  reason: string
  notes?: string
}

export interface InventoryParams {
  productId?: number
  warehouseId?: number
  locationId?: number
  batchId?: number
  page?: number
  limit?: number
}

export interface StockMovementParams {
  movementType?: string
  productId?: number
  warehouseId?: number
  locationId?: number
  batchId?: number
  referenceType?: string
  referenceId?: number
  page?: number
  limit?: number
}
