import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query"
import { productCategoryService } from "@/features/product_categories/services"
import type {
  ProductCategoryParams,
  UpdateProductCategoryRequest,
} from "@/features/product_categories/types"

export const useProductCategories = (params: ProductCategoryParams) => {
  return useQuery({
    queryKey: ["product-categories", params],
    queryFn: () => productCategoryService.getProductCategories(params),
  })
}

export const useProductCategory = (id: number) => {
  return useQuery({
    queryKey: ["product-categories", id],
    queryFn: () => productCategoryService.getProductCategory(id),
    enabled: !!id,
  })
}

export const useCreateProductCategory = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: productCategoryService.createProductCategory,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] })
    },
  })
}

export const useUpdateProductCategory = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateProductCategoryRequest }) =>
      productCategoryService.updateProductCategory(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] })
    },
  })
}

export const useDeleteProductCategory = () => {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: productCategoryService.deleteProductCategory,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] })
    },
  })
}