import axios, { type AxiosError, type AxiosInstance, type AxiosResponse } from "axios"
import { useAuthStore } from "@/stores/authStore"

const API_BASE_URL = import.meta.env.VITE_API_URL || "http://localhost:8090/api/v1"

class ApiService {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        "Content-Type": "application/json",
      },
      withCredentials: true,
    })

    this.client.interceptors.request.use(
      (config) => {
        const token = useAuthStore.getState().accessToken
        const localToken = localStorage.getItem("accessToken")
        console.log("[API] Token from store:", token ? "present" : "MISSING", "| Token from localStorage:", localToken ? "present" : "MISSING")
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        } else if (localToken) {
          config.headers.Authorization = `Bearer ${localToken}`
        }
        return config
      },
      (error) => Promise.reject(error)
    )

    // Response interceptor
    this.client.interceptors.response.use(
      (response: AxiosResponse) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config

        if (error.response?.status === 401 && originalRequest) {
          // Try to refresh token
          const refreshToken = localStorage.getItem("refreshToken")
          if (refreshToken) {
            try {
              const response = await this.post<{ accessToken: string }>("/adminPublic/refresh", {
                refreshToken,
              })
              localStorage.setItem("accessToken", response.data.accessToken)
              originalRequest.headers.Authorization = `Bearer ${response.data.accessToken}`
              return this.client(originalRequest)
            } catch {
              localStorage.removeItem("accessToken")
              localStorage.removeItem("refreshToken")
              window.location.href = "/login"
            }
          } else {
            localStorage.removeItem("accessToken")
            localStorage.removeItem("refreshToken")
            window.location.href = "/login"
          }
        }

        return Promise.reject(error)
      }
    )
  }

  async get<T>(url: string, params?: Record<string, unknown>): Promise<AxiosResponse<T>> {
    return this.client.get<T>(url, { params })
  }

  async post<T>(url: string, data?: unknown): Promise<AxiosResponse<T>> {
    return this.client.post<T>(url, data)
  }

  async put<T>(url: string, data?: unknown): Promise<AxiosResponse<T>> {
    return this.client.put<T>(url, data)
  }

  async patch<T>(url: string, data?: unknown): Promise<AxiosResponse<T>> {
    return this.client.patch<T>(url, data)
  }

  async delete<T>(url: string): Promise<AxiosResponse<T>> {
    return this.client.delete<T>(url)
  }
}

export const apiService = new ApiService()
