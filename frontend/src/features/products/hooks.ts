import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { productService } from "@/features/products/services"
import type {
  UpdateProductRequest,
  ProductParams,
} from "@/features/products/types"

export const useProducts = (params: ProductParams) => {
  return useQuery({
    queryKey: ["products", params],
    queryFn: () => productService.getProducts(params),
  })
}

export const useProduct = (id: number) => {
  return useQuery({
    queryKey: ["products", id],
    queryFn: () => productService.getProduct(id),
    enabled: !!id,
  })
}

export const useCreateProduct = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: productService.createProduct,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] })
    },
  })
}

export const useUpdateProduct = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateProductRequest }) =>
      productService.updateProduct(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] })
    },
  })
}

export const useDeleteProduct = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: productService.deleteProduct,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] })
    },
  })
}
