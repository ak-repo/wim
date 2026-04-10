import { apiService } from "@/lib/api"
import type {
  Inventory,
  AdjustInventoryRequest,
  InventoryParams,
  StockMovement,
  StockMovementParams,
} from "@/features/inventory/types"
import type { PaginatedResponse } from "@/types"

export const inventoryService = {
  getInventoryList: async (params: InventoryParams): Promise<PaginatedResponse<Inventory>> => {
    const response = await apiService.get<PaginatedResponse<Inventory>>("/admin/inventory", params as unknown as Record<string, unknown>)
    return response.data
  },

  getInventory: async (id: number): Promise<Inventory> => {
    const response = await apiService.get<Inventory>(`/admin/inventory/${id}`)
    return response.data
  },

  adjustInventory: async (data: AdjustInventoryRequest): Promise<{ message: string }> => {
    const response = await apiService.post<{ message: string }>("/admin/inventory/adjust", data)
    return response.data
  },

  getStockMovements: async (params: StockMovementParams): Promise<PaginatedResponse<StockMovement>> => {
    const response = await apiService.get<PaginatedResponse<StockMovement>>("/admin/stock-movements", params as unknown as Record<string, unknown>)
    return response.data
  },
}
