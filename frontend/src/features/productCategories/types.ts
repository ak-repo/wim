export interface ProductCategory {
  id: number
  name: string
  refCode: string
  isActive: boolean
}

export interface CreateProductCategoryRequest {
  name: string
  isActive?: boolean
}

export interface UpdateProductCategoryRequest {
  name?: string
  isActive?: boolean
}

export interface ProductCategoryParams {
  active?: boolean
  page: number
  limit: number
}
