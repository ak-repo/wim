import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { inventoryService } from "@/features/inventory/services"
import type {
  InventoryParams,
  StockMovementParams,
  AdjustInventoryRequest,
} from "@/features/inventory/types"

export const useInventoryList = (params: InventoryParams) => {
  return useQuery({
    queryKey: ["inventory", params],
    queryFn: () => inventoryService.getInventoryList(params),
  })
}

export const useInventory = (id: number) => {
  return useQuery({
    queryKey: ["inventory", id],
    queryFn: () => inventoryService.getInventory(id),
    enabled: !!id,
  })
}

export const useAdjustInventory = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: AdjustInventoryRequest) => inventoryService.adjustInventory(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["inventory"] })
      queryClient.invalidateQueries({ queryKey: ["stock-movements"] })
    },
  })
}

export const useStockMovements = (params: StockMovementParams) => {
  return useQuery({
    queryKey: ["stock-movements", params],
    queryFn: () => inventoryService.getStockMovements(params),
  })
}
