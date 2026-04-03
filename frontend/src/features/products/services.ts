import { apiService } from "@/lib/api"
import type {
  Product,
  CreateProductRequest,
  UpdateProductRequest,
  ProductParams,
} from "@/features/products/types"
import type { PaginatedResponse } from "@/types"

export const productService = {
  getProducts: async (params: ProductParams): Promise<PaginatedResponse<Product>> => {
    const response = await apiService.get<PaginatedResponse<Product>>("/admin/products", params as unknown as Record<string, unknown>)
    return response.data
  },

  getProduct: async (id: number): Promise<Product> => {
    const response = await apiService.get<Product>(`/admin/products/${id}`)
    return response.data
  },

  createProduct: async (data: CreateProductRequest): Promise<Product> => {
    const response = await apiService.post<Product>("/admin/products", data)
    return response.data
  },

  updateProduct: async (id: number, data: UpdateProductRequest): Promise<Product> => {
    const response = await apiService.patch<Product>(`/admin/products/${id}`, data)
    return response.data
  },

  deleteProduct: async (id: number): Promise<void> => {
    await apiService.delete(`/admin/products/${id}`)
  },
}
