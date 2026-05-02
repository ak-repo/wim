export interface ProductCategory {
  id: number
  name: string
  description?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateProductCategoryRequest {
  name: string
  description?: string
  isActive?: boolean
}

export interface UpdateProductCategoryRequest {
  name?: string
  description?: string
  isActive?: boolean
}

export interface ProductCategoryParams {
  active?: boolean
  page: number
  limit: number
}