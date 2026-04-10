import { apiService } from "@/lib/api"
import type {
  SalesOrder,
  CreateSalesOrderRequest,
  SalesOrderParams,
} from "@/features/salesOrders/types"
import type { PaginatedResponse } from "@/types"

export const salesOrderService = {
  getSalesOrders: async (params: SalesOrderParams): Promise<PaginatedResponse<SalesOrder>> => {
    const response = await apiService.get<PaginatedResponse<SalesOrder>>("/admin/sales-orders", params as unknown as Record<string, unknown>)
    return response.data
  },

  getSalesOrder: async (id: number): Promise<SalesOrder> => {
    const response = await apiService.get<SalesOrder>(`/admin/sales-orders/${id}`)
    return response.data
  },

  getSalesOrderByRefCode: async (refCode: string): Promise<SalesOrder> => {
    const response = await apiService.get<SalesOrder>(`/admin/sales-orders/ref?refCode=${refCode}`)
    return response.data
  },

  createSalesOrder: async (data: CreateSalesOrderRequest): Promise<SalesOrder> => {
    const response = await apiService.post<SalesOrder>("/admin/sales-orders", data)
    return response.data
  },

  updateSalesOrder: async (id: number, data: Partial<CreateSalesOrderRequest>): Promise<SalesOrder> => {
    const response = await apiService.put<SalesOrder>(`/admin/sales-orders/${id}`, data)
    return response.data
  },

  cancelSalesOrder: async (id: number): Promise<void> => {
    await apiService.patch(`/admin/sales-orders/${id}/cancel`, {})
  },

  allocateSalesOrder: async (id: number, data?: { strategy?: string; notes?: string }): Promise<void> => {
    await apiService.patch(`/admin/sales-orders/${id}/allocate`, data || {})
  },

  deallocateSalesOrder: async (id: number): Promise<void> => {
    await apiService.patch(`/admin/sales-orders/${id}/deallocate`, {})
  },

  shipSalesOrder: async (id: number, data: { shippedDate?: string; notes?: string }): Promise<void> => {
    await apiService.patch(`/admin/sales-orders/${id}/ship`, data)
  },
}