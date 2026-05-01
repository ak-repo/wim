import { apiService } from "@/lib/api"
import type {
  Customer,
  CreateCustomerRequest,
  UpdateCustomerRequest,
  CustomerParams,
} from "@/features/customers/types"
import type { PaginatedResponse } from "@/types"

interface CustomerListResponse {
  data: Customer[]
  total_count: number
  total_page: number
  current_page: number
  limit: number
}

interface CreateCustomerResponse {
  id: number
}

const mapCustomerListResponse = (response: CustomerListResponse): PaginatedResponse<Customer> => {
  return {
    data: response.data,
    total: response.total_count,
    page: response.current_page,
    limit: response.limit,
    totalPages: response.total_page,
  }
}

export const customerService = {
  getCustomers: async (params: CustomerParams): Promise<PaginatedResponse<Customer>> => {
    const response = await apiService.get<CustomerListResponse>(
      "/admin/customers",
      params as unknown as Record<string, unknown>
    )
    return mapCustomerListResponse(response.data)
  },

  getCustomer: async (id: number): Promise<Customer> => {
    const response = await apiService.get<Customer>(`/admin/customers/${id}`)
    return response.data
  },

  getCustomerByEmail: async (email: string): Promise<Customer> => {
    const response = await apiService.get<Customer>(`/admin/customers/email/${encodeURIComponent(email)}`)
    return response.data
  },

  createCustomer: async (data: CreateCustomerRequest): Promise<CreateCustomerResponse> => {
    const response = await apiService.post<CreateCustomerResponse>("/admin/customers", data)
    return response.data
  },

  updateCustomer: async (id: number, data: UpdateCustomerRequest): Promise<{ message: string }> => {
    const response = await apiService.patch<{ message: string }>(`/admin/customers/${id}`, data)
    return response.data
  },

  deleteCustomer: async (id: number): Promise<{ message: string }> => {
    const response = await apiService.delete<{ message: string }>(`/admin/customers/${id}`)
    return response.data
  },
}
