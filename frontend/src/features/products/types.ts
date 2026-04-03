// Products types
// UUID type is available in parent types file
// UUID is used in related types

export interface Product {
  id: number
  sku: string
  name: string
  description?: string
  category?: string
  unitOfMeasure: string
  weight?: number
  length?: number
  width?: number
  height?: number
  barcode?: string
  isActive: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateProductRequest {
  sku: string
  name: string
  description?: string
  category?: string
  unitOfMeasure: string
  weight?: number
  length?: number
  width?: number
  height?: number
  barcode?: string
}

export interface UpdateProductRequest {
  name?: string
  description?: string
  category?: string
  unitOfMeasure?: string
  weight?: number
  length?: number
  width?: number
  height?: number
  barcode?: string
  isActive?: boolean
}

export interface ProductParams {
  active?: boolean
  category?: string
  page: number
  limit: number
}
