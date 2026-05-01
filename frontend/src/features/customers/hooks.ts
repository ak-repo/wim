import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { customerService } from "@/features/customers/services"
import type {
  CustomerParams,
  UpdateCustomerRequest,
} from "@/features/customers/types"

export const useCustomers = (params: CustomerParams) => {
  return useQuery({
    queryKey: ["customers", params],
    queryFn: () => customerService.getCustomers(params),
  })
}

export const useCustomer = (id: number) => {
  return useQuery({
    queryKey: ["customers", id],
    queryFn: () => customerService.getCustomer(id),
    enabled: !!id,
  })
}

export const useCreateCustomer = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: customerService.createCustomer,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["customers"] })
    },
  })
}

export const useUpdateCustomer = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateCustomerRequest }) =>
      customerService.updateCustomer(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["customers"] })
    },
  })
}

export const useDeleteCustomer = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: customerService.deleteCustomer,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["customers"] })
    },
  })
}
