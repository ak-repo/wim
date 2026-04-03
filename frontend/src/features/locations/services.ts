import { apiService } from "@/lib/api"
import type {
  Location,
  CreateLocationRequest,
  UpdateLocationRequest,
  LocationParams,
} from "@/features/locations/types"
import type { PaginatedResponse } from "@/types"

export const locationService = {
  getLocations: async (params: LocationParams): Promise<PaginatedResponse<Location>> => {
    const response = await apiService.get<PaginatedResponse<Location>>("/admin/locations", params as unknown as Record<string, unknown>)
    return response.data
  },

  getLocation: async (id: string): Promise<Location> => {
    const response = await apiService.get<Location>(`/admin/locations/${id}`)
    return response.data
  },

  createLocation: async (data: CreateLocationRequest): Promise<Location> => {
    const response = await apiService.post<Location>("/admin/locations", data)
    return response.data
  },

  updateLocation: async (id: string, data: UpdateLocationRequest): Promise<Location> => {
    const response = await apiService.patch<Location>(`/admin/locations/${id}`, data)
    return response.data
  },

  deleteLocation: async (id: string): Promise<void> => {
    await apiService.delete(`/admin/locations/${id}`)
  },
}
