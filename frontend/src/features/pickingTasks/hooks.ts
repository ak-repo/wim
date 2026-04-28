import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { pickingTaskService } from "@/features/pickingTasks/services"
import { useAuthStore } from "@/stores/authStore"
import type {
  CreatePickingTaskRequest,
  AssignPickingTaskRequest,
  PickItemRequest,
  CompletePickingRequest,
  PickingTaskParams,
} from "@/types/pickingTask"

export const usePickingTasks = (params: PickingTaskParams) => {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const hasToken = !!localStorage.getItem("accessToken")

  return useQuery({
    queryKey: ["pickingTasks", params],
    queryFn: () => pickingTaskService.getPickingTasks(params),
    enabled: isAuthenticated && hasToken,
  })
}

export const usePickingTask = (id: string) => {
  return useQuery({
    queryKey: ["pickingTasks", id],
    queryFn: () => pickingTaskService.getPickingTask(id),
    enabled: !!id,
  })
}

export const useCreatePickingTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: pickingTaskService.createPickingTask,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["pickingTasks"] })
    },
  })
}

export const useAssignPickingTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: AssignPickingTaskRequest }) =>
      pickingTaskService.assignPickingTask(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["pickingTasks"] })
      queryClient.invalidateQueries({ queryKey: ["pickingTasks", variables.id] })
    },
  })
}

export const useStartPickingTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: pickingTaskService.startPickingTask,
    onSuccess: (_data, id) => {
      queryClient.invalidateQueries({ queryKey: ["pickingTasks"] })
      queryClient.invalidateQueries({ queryKey: ["pickingTasks", id] })
    },
  })
}

export const usePickItem = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: PickItemRequest }) =>
      pickingTaskService.pickItem(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["pickingTasks"] })
      queryClient.invalidateQueries({ queryKey: ["pickingTasks", variables.id] })
      queryClient.invalidateQueries({ queryKey: ["inventory"] })
    },
  })
}

export const useCompletePickingTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: CompletePickingRequest }) =>
      pickingTaskService.completePickingTask(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["pickingTasks"] })
      queryClient.invalidateQueries({ queryKey: ["pickingTasks", variables.id] })
    },
  })
}

export const useCancelPickingTask = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, notes }: { id: string; notes?: string }) =>
      pickingTaskService.cancelPickingTask(id, notes),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["pickingTasks"] })
      queryClient.invalidateQueries({ queryKey: ["pickingTasks", variables.id] })
    },
  })
}