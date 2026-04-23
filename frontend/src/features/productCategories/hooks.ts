import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { productCategoryService } from "@/features/productCategories/services"
import type {
  UpdateProductCategoryRequest,
  ProductCategoryParams,
} from "@/features/productCategories/types"

export const useProductCategories = (params: ProductCategoryParams) => {
  return useQuery({
    queryKey: ["product-categories", params],
    queryFn: () => productCategoryService.getCategories(params),
  })
}

export const useProductCategory = (id: number) => {
  return useQuery({
    queryKey: ["product-categories", id],
    queryFn: () => productCategoryService.getCategory(id),
    enabled: !!id,
  })
}

export const useCreateProductCategory = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: productCategoryService.createCategory,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] })
    },
  })
}

export const useUpdateProductCategory = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateProductCategoryRequest }) =>
      productCategoryService.updateCategory(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] })
    },
  })
}

export const useDeleteProductCategory = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: productCategoryService.deleteCategory,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] })
    },
  })
}
