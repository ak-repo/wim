import { apiService } from "@/lib/api"
import type {
  PurchaseOrderResponse,
  PurchaseOrderParams,
  CreatePurchaseOrderRequest,
  ReceivePurchaseOrderRequest,
  PutAwayPurchaseOrderRequest,
  PurchaseOrdersData,
} from "@/types/purchaseOrder"

export const purchaseOrderService = {
  getPurchaseOrders: async (params: PurchaseOrderParams): Promise<PurchaseOrdersData> => {
    const response = await apiService.get<PurchaseOrdersData>("/admin/purchase-orders", params as unknown as Record<string, unknown>)
    return response.data
  },

  getPurchaseOrder: async (id: string): Promise<PurchaseOrderResponse> => {
    const response = await apiService.get<PurchaseOrderResponse>(`/admin/purchase-orders/${id}`)
    return response.data
  },

  createPurchaseOrder: async (data: CreatePurchaseOrderRequest): Promise<PurchaseOrderResponse> => {
    const response = await apiService.post<PurchaseOrderResponse>("/admin/purchase-orders", data)
    return response.data
  },

  receivePurchaseOrder: async (id: string, data: ReceivePurchaseOrderRequest): Promise<void> => {
    await apiService.post(`/admin/purchase-orders/${id}/receive`, data)
  },

  putAwayPurchaseOrder: async (id: string, data: PutAwayPurchaseOrderRequest): Promise<void> => {
    await apiService.post(`/admin/purchase-orders/${id}/put-away`, data)
  },

  getPurchaseOrderByRefCode: async (refCode: string): Promise<PurchaseOrderResponse> => {
    const response = await apiService.get<PurchaseOrderResponse>("/admin/purchase-orders/ref", { refCode })
    return response.data
  },
}