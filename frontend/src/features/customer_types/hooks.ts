import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { customerTypeService } from "@/features/customer_types/services"
import type {
  CustomerTypeParams,
  UpdateCustomerTypeRequest,
} from "@/features/customer_types/types"

export const useCustomerTypes = (params: CustomerTypeParams) => {
  return useQuery({
    queryKey: ["customer-types", params],
    queryFn: () => customerTypeService.getCustomerTypes(params),
  })
}

export const useCustomerType = (id: number) => {
  return useQuery({
    queryKey: ["customer-types", id],
    queryFn: () => customerTypeService.getCustomerType(id),
    enabled: !!id,
  })
}

export const useCreateCustomerType = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: customerTypeService.createCustomerType,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["customer-types"] })
    },
  })
}

export const useUpdateCustomerType = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateCustomerTypeRequest }) =>
      customerTypeService.updateCustomerType(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["customer-types"] })
    },
  })
}

export const useDeleteCustomerType = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: customerTypeService.deleteCustomerType,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["customer-types"] })
    },
  })
}