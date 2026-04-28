import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { authService, userService } from "@/features/auth/services"
import { useAuthStore } from "@/stores/authStore"
import type {
  LoginRequest,
  UserRequest,
  UserParams,
} from "@/features/auth/types"

// Auth mutations
export const useLogin = () => {
  const { setAuth } = useAuthStore()

  return useMutation({
    mutationFn: async (data: LoginRequest) => {
      const response = await authService.login(data)
      const { accessToken, refreshToken } = response
      
      localStorage.setItem("accessToken", accessToken)
      localStorage.setItem("refreshToken", refreshToken || "")
      
      const user = await authService.me()
      return { ...response, user }
    },
    onSuccess: ({ accessToken, refreshToken, user }) => {
      setAuth(user, accessToken, refreshToken || "")
    },
  })
}

export const useRegister = () => {
  return useMutation({
    mutationFn: authService.register,
  })
}

export const useLogout = () => {
  const { logout } = useAuthStore()

  return () => {
    authService.logout()
    logout()
  }
}

// User queries
export const useUsers = (params: UserParams) => {
  return useQuery({
    queryKey: ["users", params],
    queryFn: () => userService.getUsers(params),
  })
}

export const useCreateUser = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: userService.createUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] })
    },
  })
}

export const useUpdateUser = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<UserRequest> }) =>
      userService.updateUser(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] })
    },
  })
}

export const useDeleteUser = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: userService.deleteUser,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] })
    },
  })
}
