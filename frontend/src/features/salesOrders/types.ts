// Sales Orders types

export interface SalesOrder {
  id: number
  refCode: string
  customerId: number
  warehouseId: number
  status: string
  allocationStatus: string
  orderDate: string
  requiredDate?: string
  shippedDate?: string
  shippingMethod?: string
  shippingAddress?: string
  billingAddress?: string
  notes?: string
  createdBy?: number
  createdAt: string
  updatedAt: string
  items?: SalesOrderItem[]
}

export interface SalesOrderItem {
  id: number
  salesOrderId: number
  productId: number
  quantityOrdered: number
  quantityShipped: number
  quantityReserved: number
  unitPrice?: number
  allocationStatus: string
  batchId?: number
  allocatedLocationId?: number
  createdAt: string
  updatedAt: string
}

export interface CreateSalesOrderRequest {
  customerId: number
  warehouseId: number
  requiredDate?: string
  shippingMethod?: string
  shippingAddress?: string
  billingAddress?: string
  notes?: string
  items: SalesOrderItemRequest[]
}

export interface SalesOrderItemRequest {
  productId: number
  quantityOrdered: number
  unitPrice?: number
}

export interface AllocateSalesOrderRequest {
  strategy?: string
  notes?: string
}

export interface ShipSalesOrderRequest {
  shippedDate?: string
  notes?: string
  items: ShipSalesOrderItemRequest[]
}

export interface ShipSalesOrderItemRequest {
  salesOrderItemId: number
  quantityShipped: number
  locationId?: number
  batchId?: number
}

export interface SalesOrderParams {
  customerId?: number
  warehouseId?: number
  status?: string
  allocationStatus?: string
  page: number
  limit: number
}