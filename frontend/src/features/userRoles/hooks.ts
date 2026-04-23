import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { userRoleService } from "@/features/userRoles/services"
import type { UpdateUserRoleRequest, UserRoleParams } from "@/features/userRoles/types"

export const useUserRoles = (params: UserRoleParams) => {
  return useQuery({
    queryKey: ["user-roles", params],
    queryFn: () => userRoleService.getRoles(params),
  })
}

export const useUserRole = (id: number) => {
  return useQuery({
    queryKey: ["user-roles", id],
    queryFn: () => userRoleService.getRole(id),
    enabled: !!id,
  })
}

export const useCreateUserRole = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: userRoleService.createRole,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-roles"] })
    },
  })
}

export const useUpdateUserRole = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateUserRoleRequest }) =>
      userRoleService.updateRole(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-roles"] })
    },
  })
}

export const useDeleteUserRole = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: userRoleService.deleteRole,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user-roles"] })
    },
  })
}
