import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { salesOrderService } from "@/features/salesOrders/services"
import type {
  CreateSalesOrderRequest,
  SalesOrderParams,
} from "@/features/salesOrders/types"

export const useSalesOrders = (params: SalesOrderParams) => {
  return useQuery({
    queryKey: ["sales-orders", params],
    queryFn: () => salesOrderService.getSalesOrders(params),
  })
}

export const useCreateSalesOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateSalesOrderRequest) => salesOrderService.createSalesOrder(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sales-orders"] })
    },
  })
}

export const useUpdateSalesOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<CreateSalesOrderRequest> }) =>
      salesOrderService.updateSalesOrder(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sales-orders"] })
    },
  })
}

export const useCancelSalesOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => salesOrderService.cancelSalesOrder(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sales-orders"] })
    },
  })
}

export const useAllocateSalesOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => salesOrderService.allocateSalesOrder(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sales-orders"] })
    },
  })
}

export const useDeallocateSalesOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => salesOrderService.deallocateSalesOrder(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sales-orders"] })
    },
  })
}

export const useShipSalesOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: { shippedDate?: string; notes?: string } }) =>
      salesOrderService.shipSalesOrder(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["sales-orders"] })
    },
  })
}