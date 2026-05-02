import { apiService } from "@/lib/api"
import type {
  CustomerType,
  CreateCustomerTypeRequest,
  UpdateCustomerTypeRequest,
  CustomerTypeParams,
} from "@/features/customer_types/types"
import type { PaginatedResponse } from "@/types"

interface CustomerTypeListResponse {
  data: CustomerType[]
  total_count: number
  total_page: number
  current_page: number
  limit: number
}

interface CreateCustomerTypeResponse {
  id: number
}

interface UpdateCustomerTypeResponse {
  message: string
}

interface DeleteCustomerTypeResponse {
  message: string
}

const mapCustomerTypeListResponse = (response: CustomerTypeListResponse): PaginatedResponse<CustomerType> => {
  return {
    data: response.data,
    total: response.total_count,
    page: response.current_page,
    limit: response.limit,
    totalPages: response.total_page,
  }
}

export const customerTypeService = {
  getCustomerTypes: async (params: CustomerTypeParams): Promise<PaginatedResponse<CustomerType>> => {
    const response = await apiService.get<CustomerTypeListResponse>(
      "/admin/customer-types",
      params as unknown as Record<string, unknown>
    )
    return mapCustomerTypeListResponse(response.data)
  },

  getCustomerType: async (id: number): Promise<CustomerType> => {
    const response = await apiService.get<CustomerType>(`/admin/customer-types/${id}`)
    return response.data
  },

  createCustomerType: async (data: CreateCustomerTypeRequest): Promise<CreateCustomerTypeResponse> => {
    const response = await apiService.post<CreateCustomerTypeResponse>("/admin/customer-types", data)
    return response.data
  },

  updateCustomerType: async (id: number, data: UpdateCustomerTypeRequest): Promise<UpdateCustomerTypeResponse> => {
    const response = await apiService.put<UpdateCustomerTypeResponse>(`/admin/customer-types/${id}`, data)
    return response.data
  },

  deleteCustomerType: async (id: number): Promise<DeleteCustomerTypeResponse> => {
    const response = await apiService.delete<DeleteCustomerTypeResponse>(`/admin/customer-types/${id}`)
    return response.data
  },
}