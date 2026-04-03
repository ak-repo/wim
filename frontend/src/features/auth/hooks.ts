import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { authService, userService } from "@/features/auth/services"
import { useAuthStore } from "@/stores/authStore"
import type {
  LoginRequest,
  UserRequest,
  UserParams,
} from "@/features/auth/types"
import type { User } from "@/features/auth/types"

// Auth mutations
export const useLogin = () => {
  const { setAuth } = useAuthStore()

  return useMutation({
    mutationFn: async (data: LoginRequest) => {
      const response = await authService.login(data)
      // Get user info - we'll need to parse JWT or make a separate call
      // For now, we'll store just the tokens
      return response
    },
    onSuccess: (data) => {
      // Mock user data since backend doesn't return it directly
      // In a real app, you'd decode the JWT or fetch user profile
      const mockUser: User = {
        id: "",
        username: "",
        email: "",
        role: "admin",
        isActive: true,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }
      setAuth(mockUser, data.accessToken, data.refreshToken || "")
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
    mutationFn: ({ id, data }: { id: string; data: Partial<UserRequest> }) =>
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
