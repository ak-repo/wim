import { apiService } from "@/lib/api"
import type {
  ProductCategory,
  CreateProductCategoryRequest,
  UpdateProductCategoryRequest,
  ProductCategoryParams,
} from "@/features/product_categories/types"
import type { PaginatedResponse } from "@/types"

interface ProductCategoryListResponse {
  data: ProductCategory[]
  total_count: number
  total_page: number
  current_page: number
  limit: number
}

interface CreateProductCategoryResponse {
  id: number
}

interface UpdateProductCategoryResponse {
  message: string
}

interface DeleteProductCategoryResponse {
  message: string
}

const mapProductCategoryListResponse = (response: ProductCategoryListResponse): PaginatedResponse<ProductCategory> => {
  return {
    data: response.data,
    total: response.total_count,
    page: response.current_page,
    limit: response.limit,
    totalPages: response.total_page,
  }
}

export const productCategoryService = {
  getProductCategories: async (params: ProductCategoryParams): Promise<PaginatedResponse<ProductCategory>> => {
    const response = await apiService.get<ProductCategoryListResponse>(
      "/admin/product-categories",
      params as unknown as Record<string, unknown>
    )
    return mapProductCategoryListResponse(response.data)
  },

  getProductCategory: async (id: number): Promise<ProductCategory> => {
    const response = await apiService.get<ProductCategory>(`/admin/product-categories/${id}`)
    return response.data
  },

  createProductCategory: async (data: CreateProductCategoryRequest): Promise<CreateProductCategoryResponse> => {
    const response = await apiService.post<CreateProductCategoryResponse>("/admin/product-categories", data)
    return response.data
  },

  updateProductCategory: async (id: number, data: UpdateProductCategoryRequest): Promise<UpdateProductCategoryResponse> => {
    const response = await apiService.put<UpdateProductCategoryResponse>(`/admin/product-categories/${id}`, data)
    return response.data
  },

  deleteProductCategory: async (id: number): Promise<DeleteProductCategoryResponse> => {
    const response = await apiService.delete<DeleteProductCategoryResponse>(`/admin/product-categories/${id}`)
    return response.data
  },
}