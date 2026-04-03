import { apiService } from "@/lib/api"
import type {
  Warehouse,
  CreateWarehouseRequest,
  UpdateWarehouseRequest,
  WarehouseParams,
} from "@/features/warehouses/types"
import type { PaginatedResponse } from "@/types"

export const warehouseService = {
  getWarehouses: async (params: WarehouseParams): Promise<PaginatedResponse<Warehouse>> => {
    const response = await apiService.get<PaginatedResponse<Warehouse>>("/admin/warehouses", params as unknown as Record<string, unknown>)
    return response.data
  },

  getWarehouse: async (id: string): Promise<Warehouse> => {
    const response = await apiService.get<Warehouse>(`/admin/warehouses/${id}`)
    return response.data
  },

  createWarehouse: async (data: CreateWarehouseRequest): Promise<Warehouse> => {
    const response = await apiService.post<Warehouse>("/admin/warehouses", data)
    return response.data
  },

  updateWarehouse: async (id: string, data: UpdateWarehouseRequest): Promise<Warehouse> => {
    const response = await apiService.patch<Warehouse>(`/admin/warehouses/${id}`, data)
    return response.data
  },

  deleteWarehouse: async (id: string): Promise<void> => {
    await apiService.delete(`/admin/warehouses/${id}`)
  },
}
