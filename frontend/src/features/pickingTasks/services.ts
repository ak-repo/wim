import { apiService } from "@/lib/api"
import type {
  PickingTaskResponse,
  PickingTasksData,
  PickingTaskParams,
  CreatePickingTaskRequest,
  AssignPickingTaskRequest,
  PickItemRequest,
  CompletePickingRequest,
} from "@/types/pickingTask"

export const pickingTaskService = {
  getPickingTasks: async (params: PickingTaskParams): Promise<PickingTasksData> => {
    const response = await apiService.get<PickingTasksData>("/admin/picking-tasks", params as unknown as Record<string, unknown>)
    return response.data
  },

  getPickingTask: async (id: string): Promise<PickingTaskResponse> => {
    const response = await apiService.get<PickingTaskResponse>(`/admin/picking-tasks/${id}`)
    return response.data
  },

  createPickingTask: async (data: CreatePickingTaskRequest): Promise<PickingTaskResponse> => {
    const response = await apiService.post<PickingTaskResponse>("/admin/picking-tasks", data)
    return response.data
  },

  assignPickingTask: async (id: string, data: AssignPickingTaskRequest): Promise<void> => {
    await apiService.post(`/admin/picking-tasks/${id}/assign`, data)
  },

  startPickingTask: async (id: string): Promise<void> => {
    await apiService.post(`/admin/picking-tasks/${id}/start`, {})
  },

  pickItem: async (id: string, data: PickItemRequest): Promise<void> => {
    await apiService.post(`/admin/picking-tasks/${id}/pick`, data)
  },

  completePickingTask: async (id: string, data: CompletePickingRequest): Promise<void> => {
    await apiService.post(`/admin/picking-tasks/${id}/complete`, data)
  },

  cancelPickingTask: async (id: string, notes?: string): Promise<void> => {
    await apiService.post(`/admin/picking-tasks/${id}/cancel`, { notes })
  },
}