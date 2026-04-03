import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { warehouseService } from "@/features/warehouses/services"
import type {
  UpdateWarehouseRequest,
  WarehouseParams,
} from "@/features/warehouses/types"

export const useWarehouses = (params: WarehouseParams) => {
  return useQuery({
    queryKey: ["warehouses", params],
    queryFn: () => warehouseService.getWarehouses(params),
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
