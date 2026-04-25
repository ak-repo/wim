import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { warehouseService } from "@/features/warehouses/services"
import { useAuthStore } from "@/stores/authStore"
import type {
  UpdateWarehouseRequest,
  WarehouseParams,
} from "@/features/warehouses/types"

export const useWarehouses = (params: WarehouseParams) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const hasToken = !!localStorage.getItem("accessToken")
  console.log("[useWarehouses] enabled check - isAuthenticated:", isAuthenticated, "hasToken:", hasToken)
  
  return useQuery({
    queryKey: ["warehouses", params],
    queryFn: () => warehouseService.getWarehouses(params),
    enabled: isAuthenticated && hasToken,
  })
}

export const useWarehouse = (id: string) => {
  return useQuery({
    queryKey: ["warehouses", id],
    queryFn: () => warehouseService.getWarehouse(id),
    enabled: !!id,
  })
}

export const useCreateWarehouse = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: warehouseService.createWarehouse,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["warehouses"] })
    },
  })
}

export const useUpdateWarehouse = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateWarehouseRequest }) =>
      warehouseService.updateWarehouse(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["warehouses"] })
    },
  })
}

export const useDeleteWarehouse = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: warehouseService.deleteWarehouse,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["warehouses"] })
    },
  })
}
