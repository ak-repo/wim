import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { locationService } from "@/features/locations/services"
import type {
  UpdateLocationRequest,
  LocationParams,
} from "@/features/locations/types"

export const useLocations = (params: LocationParams) => {
  return useQuery({
    queryKey: ["locations", params],
    queryFn: () => locationService.getLocations(params),
  })
}

export const useLocation = (id: string) => {
  return useQuery({
    queryKey: ["locations", id],
    queryFn: () => locationService.getLocation(id),
    enabled: !!id,
  })
}

export const useCreateLocation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: locationService.createLocation,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["locations"] })
    },
  })
}

export const useUpdateLocation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateLocationRequest }) =>
      locationService.updateLocation(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["locations"] })
    },
  })
}

export const useDeleteLocation = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: locationService.deleteLocation,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["locations"] })
    },
  })
}
