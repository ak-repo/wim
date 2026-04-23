import { apiService } from "@/lib/api"
import type {
  ProductCategory,
  CreateProductCategoryRequest,
  UpdateProductCategoryRequest,
  ProductCategoryParams,
} from "@/features/productCategories/types"
import type { PaginatedResponse } from "@/types"

export const productCategoryService = {
  getCategories: async (
    params: ProductCategoryParams
  ): Promise<PaginatedResponse<ProductCategory>> => {
    const response = await apiService.get<PaginatedResponse<ProductCategory>>(
      "/admin/product-categories",
      params as unknown as Record<string, unknown>
    )
    return response.data
  },

  getCategory: async (id: number): Promise<ProductCategory> => {
    const response = await apiService.get<ProductCategory>(
      `/admin/product-categories/${id}`
    )
    return response.data
  },

  createCategory: async (
    data: CreateProductCategoryRequest
  ): Promise<ProductCategory> => {
    const response = await apiService.post<ProductCategory>(
      "/admin/product-categories",
      data
    )
    return response.data
  },

  updateCategory: async (
    id: number,
    data: UpdateProductCategoryRequest
  ): Promise<ProductCategory> => {
    const response = await apiService.patch<ProductCategory>(
      `/admin/product-categories/${id}`,
      data
    )
    return response.data
  },

  deleteCategory: async (id: number): Promise<void> => {
    await apiService.delete(`/admin/product-categories/${id}`)
  },
}
