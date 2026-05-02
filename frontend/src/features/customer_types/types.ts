export interface CustomerType {
  id: number
  name: string
  description?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateCustomerTypeRequest {
  name: string
  description?: string
  isActive?: boolean
}

export interface UpdateCustomerTypeRequest {
  name?: string
  description?: string
  isActive?: boolean
}

export interface CustomerTypeParams {
  active?: boolean
  page: number
  limit: number
}