// Purchase Order types
export interface PurchaseOrderItemRequest {
  productId: number
  batchNumber?: string
  quantityOrdered: number
  unitPrice?: number
}

export interface CreatePurchaseOrderRequest {
  supplierId: number
  warehouseId: number
  expectedDate?: string
  notes?: string
  items: PurchaseOrderItemRequest[]
}

export interface ReceivePurchaseOrderItemRequest {
  purchaseOrderItemId: number
  quantityReceived: number
  locationId?: number
  batchId?: number
}

export interface ReceivePurchaseOrderRequest {
  receivedDate?: string
  notes?: string
  items: ReceivePurchaseOrderItemRequest[]
}

export interface PutAwayPurchaseOrderItemRequest {
  purchaseOrderItemId: number
  quantity: number
  fromLocationId: number
  toLocationId: number
  batchId?: number
}

export interface PutAwayPurchaseOrderRequest {
  notes?: string
  items: PutAwayPurchaseOrderItemRequest[]
}

export interface PurchaseOrderItemResponse {
  id: number
  purchaseOrderId: number
  productId: number
  quantityOrdered: number
  quantityReceived: number
  batchNumber?: string
  unitPrice?: number
  createdAt: string
  updatedAt: string
}

export interface PurchaseOrderResponse {
  id: number
  refCode: string
  supplierId: number
  warehouseId: number
  status: string
  expectedDate?: string
  receivedDate?: string
  notes?: string
  createdBy?: number
  createdAt: string
  updatedAt: string
  items?: PurchaseOrderItemResponse[]
}

export interface PurchaseOrderParams {
  supplierId?: number
  warehouseId?: number
  status?: string
  page?: number
  limit?: number
}

export interface PurchaseOrdersData {
  data: PurchaseOrderResponse[]
  total_count: number
  total_page: number
  current_page: number
  limit: number
}
