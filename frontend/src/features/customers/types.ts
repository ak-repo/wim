export interface Customer {
  id: number
  refCode: string
  name: string
  email: string
  contact?: string
  address?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateCustomerRequest {
  name: string
  email: string
  contact?: string
  address?: string
  isActive?: boolean
}

export interface UpdateCustomerRequest {
  name?: string
  email?: string
  contact?: string
  address?: string
  isActive?: boolean
}

export interface CustomerParams {
  active?: boolean
  page: number
  limit: number
}
