import { create } from "zustand"
import { persist } from "zustand/middleware"
import type { User } from "@/features/auth/types"

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  setAuth: (user: User, accessToken: string, refreshToken: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      isAuthenticated: false,
      setAuth: (user, accessToken, refreshToken) => {
        console.log("[AuthStore] setAuth called", { hasToken: !!accessToken })
        localStorage.setItem("accessToken", accessToken)
        localStorage.setItem("refreshToken", refreshToken)
        set({ user, accessToken, refreshToken, isAuthenticated: true })
      },
      logout: () => {
        console.log("[AuthStore] logout called")
        localStorage.removeItem("accessToken")
        localStorage.removeItem("refreshToken")
        set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false })
      },
    }),
    {
      name: "auth-storage",
      partialize: (state) => ({ user: state.user, accessToken: state.accessToken, refreshToken: state.refreshToken, isAuthenticated: state.isAuthenticated }),
      onRehydrateStorage: () => (state) => {
        console.log("[AuthStore] Rehydrated", { 
          isAuthenticated: state?.isAuthenticated, 
          hasToken: !!localStorage.getItem("accessToken"),
          zustandHasToken: !!state?.accessToken 
        })
      },
    }
  )
)

// Debug: Check initial state on module load
console.log("[AuthStore] Initial state check - isAuthenticated:", useAuthStore.getState().isAuthenticated)
