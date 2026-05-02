import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { userRoleService } from "@/features/user_roles/services"
import type {
  UserRoleParams,
  UpdateUserRoleRequest,
} from "@/features/user_roles/types"

export const useUserRoles = (params: UserRoleParams) => {
  return useQuery({
    queryKey: ["user-roles", params],
    queryFn: () => userRoleService.getUserRoles(params),
  })
}

export const useUserRole = (id: number) => {
  return useQuery({
    queryKey: ["user-roles", id],
    queryFn: () => userRoleService.getUserRole(id),
    enabled: !!id,
  })
}

export const useCreateUserRole = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: userRoleService.createUserRole,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-roles"] })
    },
  })
}

export const useUpdateUserRole = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserRoleRequest }) =>
      userRoleService.updateUserRole(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-roles"] })
    },
  })
}

export const useDeleteUserRole = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: userRoleService.deleteUserRole,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-roles"] })
    },
  })
}