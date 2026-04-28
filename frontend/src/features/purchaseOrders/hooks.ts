import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { purchaseOrderService } from "@/features/purchaseOrders/services"
import { useAuthStore } from "@/stores/authStore"
import type {
  CreatePurchaseOrderRequest,
  ReceivePurchaseOrderRequest,
  PutAwayPurchaseOrderRequest,
  PurchaseOrderParams,
} from "@/types/purchaseOrder"

export const usePurchaseOrders = (params: PurchaseOrderParams) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const hasToken = !!localStorage.getItem("accessToken")
  
  return useQuery({
    queryKey: ["purchaseOrders", params],
    queryFn: () => purchaseOrderService.getPurchaseOrders(params),
    enabled: isAuthenticated && hasToken,
  })
}

export const usePurchaseOrder = (id: string) => {
  return useQuery({
    queryKey: ["purchaseOrders", id],
    queryFn: () => purchaseOrderService.getPurchaseOrder(id),
    enabled: !!id,
  })
}

export const useCreatePurchaseOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: purchaseOrderService.createPurchaseOrder,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["purchaseOrders"] })
    },
  })
}

export const useReceivePurchaseOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ReceivePurchaseOrderRequest }) =>
      purchaseOrderService.receivePurchaseOrder(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["purchaseOrders"] })
      queryClient.invalidateQueries({ queryKey: ["purchaseOrders", variables.id] })
    },
  })
}

export const usePutAwayPurchaseOrder = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: PutAwayPurchaseOrderRequest }) =>
      purchaseOrderService.putAwayPurchaseOrder(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["purchaseOrders"] })
      queryClient.invalidateQueries({ queryKey: ["purchaseOrders", variables.id] })
      queryClient.invalidateQueries({ queryKey: ["inventory"] })
    },
  })
}
