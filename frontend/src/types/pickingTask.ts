// Picking Task types
export interface PickingTaskItemRequest {
  pickingTaskItemId: number
  quantity: number
  locationId: number
  batchId?: number
}

export interface CreatePickingTaskRequest {
  salesOrderId: number
  priority?: string
  notes?: string
}

export interface AssignPickingTaskRequest {
  assignedTo: number
  notes?: string
}

export interface PickItemRequest {
  pickingTaskItemId: number
  quantity: number
  locationId: number
  batchId?: number
}

export interface CompletePickingRequest {
  items: PickItemRequest[]
  notes?: string
}

export interface PickingTaskItemResponse {
  id: number
  pickingTaskId: number
  salesOrderItemId: number
  productId: number
  productName?: string
  locationId?: number
  locationCode?: string
  batchId?: number
  quantityRequired: number
  quantityPicked: number
  pickedAt?: string
  status: string
  createdAt: string
  updatedAt: string
}

export interface PickingTaskResponse {
  id: number
  refCode: string
  salesOrderId: number
  warehouseId: number
  status: string
  priority: string
  assignedTo?: number
  assignedUser?: string
  startedAt?: string
  completedAt?: string
  notes?: string
  createdBy?: number
  createdAt: string
  updatedAt: string
  items?: PickingTaskItemResponse[]
}

export interface PickingTaskParams {
  warehouseId?: number
  status?: string
  priority?: string
  assignedTo?: number
  page?: number
  limit?: number
}

export interface PickingTasksData {
  data: PickingTaskResponse[]
  total_count: number
  total_page: number
  current_page: number
  limit: number
}