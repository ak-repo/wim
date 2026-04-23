import * as React from "react"
import { useProducts } from "@/features/products/hooks"
import { Button } from "@/components/ui/Button"
import { Input } from "@/components/ui/Input"
import { Badge } from "@/components/ui/Badge"
import { Skeleton } from "@/components/ui/skeleton"
import { Tabs } from "@/components/ui/Tabs"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table"
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/Card"
import { ProductFormDialog } from "@/features/products/components/ProductFormDialog"
import { ProductDeleteDialog } from "@/features/products/components/ProductDeleteDialog"
import type { Product } from "@/features/products/types"
import { Plus, Search, Pencil, Trash2, Package, Sparkles, Layers } from "lucide-react"
import { useProductCategories } from "@/features/productCategories/hooks"
import type { ProductCategory } from "@/features/productCategories/types"
import { ProductCategoryFormDialog } from "@/features/productCategories/components/ProductCategoryFormDialog"
import { ProductCategoryDeleteDialog } from "@/features/productCategories/components/ProductCategoryDeleteDialog"

export default function ProductsPage() {
  const [search, setSearch] = React.useState("")
  const [selectedProduct, setSelectedProduct] = React.useState<Product | null>(null)
  const [isFormOpen, setIsFormOpen] = React.useState(false)
  const [isDeleteOpen, setIsDeleteOpen] = React.useState(false)
  const [activeTab, setActiveTab] = React.useState("products")
  const [selectedCategory, setSelectedCategory] = React.useState<ProductCategory | null>(null)
  const [isCategoryFormOpen, setIsCategoryFormOpen] = React.useState(false)
  const [isCategoryDeleteOpen, setIsCategoryDeleteOpen] = React.useState(false)

  const { data: productsData, isLoading } = useProducts({
    page: 1,
    limit: 10,
  })

  const { data: categoriesData, isLoading: categoriesLoading } = useProductCategories({
    page: 1,
    limit: 10,
  })

  const handleEdit = (product: Product) => {
    setSelectedProduct(product)
    setIsFormOpen(true)
  }

  const handleDelete = (product: Product) => {
    setSelectedProduct(product)
    setIsDeleteOpen(true)
  }

  const handleCreate = () => {
    setSelectedProduct(null)
    setIsFormOpen(true)
  }

  const handleEditCategory = (category: ProductCategory) => {
    setSelectedCategory(category)
    setIsCategoryFormOpen(true)
  }

  const handleDeleteCategory = (category: ProductCategory) => {
    setSelectedCategory(category)
    setIsCategoryDeleteOpen(true)
  }

  const handleCreateCategory = () => {
    setSelectedCategory(null)
    setIsCategoryFormOpen(true)
  }

  const filteredProducts =
    productsData?.data?.filter(
      (product) =>
        product.name.toLowerCase().includes(search.toLowerCase()) ||
        product.sku.toLowerCase().includes(search.toLowerCase()) ||
        (product.category?.toLowerCase() || "").includes(search.toLowerCase())
    ) || []

  const filteredCategories =
    categoriesData?.data?.filter((category) =>
      category.name.toLowerCase().includes(search.toLowerCase())
    ) || []

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-foreground">Product Master</h1>
          <p className="text-muted-foreground mt-1">Manage products and categories in one place.</p>
        </div>
        <div className="flex items-center gap-3">
          <Tabs
            value={activeTab}
            onChange={setActiveTab}
            tabs={[
              { id: "products", label: "Products", icon: Package },
              { id: "categories", label: "Categories", icon: Layers },
            ]}
          />
          <Button
            onClick={activeTab === "products" ? handleCreate : handleCreateCategory}
            size="lg"
          >
            <Plus className="h-4 w-4 mr-2" />
            {activeTab === "products" ? "Add Product" : "Add Category"}
          </Button>
        </div>
      </div>

      {activeTab === "products" ? (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Product Catalog</CardTitle>
                <CardDescription>
                  Search and manage your products
                </CardDescription>
              </div>
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                <Sparkles className="h-4 w-4 text-primary" />
              </div>
            </div>
            <div className="flex items-center gap-4 mt-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search by name, SKU, or category..."
                  className="pl-9 bg-muted/50 border-transparent focus:bg-background"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>SKU</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Category</TableHead>
                  <TableHead>UoM</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-48" /></TableCell>
                      <TableCell><Skeleton className="h-6 w-20" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-16" /></TableCell>
                      <TableCell><Skeleton className="h-6 w-16" /></TableCell>
                      <TableCell>
                        <div className="flex justify-end gap-2">
                          <Skeleton className="h-8 w-8" />
                          <Skeleton className="h-8 w-8" />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                ) : filteredProducts.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-40">
                      <div className="flex flex-col items-center justify-center text-center">
                        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                          <Package className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <p className="text-sm font-medium text-foreground">No products found</p>
                        <p className="text-xs text-muted-foreground mt-1">
                          {search ? "Try adjusting your search" : "Add a product to get started"}
                        </p>
                        {!search && (
                          <Button variant="outline" size="sm" className="mt-3" onClick={handleCreate}>
                            <Plus className="h-3 w-3 mr-1" />
                            Add Product
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredProducts.map((product) => (
                    <TableRow key={product.id}>
                      <TableCell>
                        <span className="font-mono text-xs font-medium">{product.sku}</span>
                      </TableCell>
                      <TableCell>
                        <div>
                          <p className="font-medium text-foreground">{product.name}</p>
                          {product.description && (
                            <p className="text-xs text-muted-foreground line-clamp-1">{product.description}</p>
                          )}
                        </div>
                      </TableCell>
                      <TableCell>
                        {product.category ? (
                          <Badge variant="outline">{product.category}</Badge>
                        ) : (
                          <span className="text-muted-foreground text-sm">-</span>
                        )}
                      </TableCell>
                      <TableCell className="text-muted-foreground text-sm">{product.unitOfMeasure}</TableCell>
                      <TableCell>
                        <Badge variant={product.isActive ? "success" : "secondary"}>
                          {product.isActive ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleEdit(product)}
                            className="h-8 w-8"
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDelete(product)}
                            className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <div>
                <CardTitle>Product Categories</CardTitle>
                <CardDescription>
                  Manage category definitions used by products
                </CardDescription>
              </div>
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                <Layers className="h-4 w-4 text-primary" />
              </div>
            </div>
            <div className="flex items-center gap-4 mt-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search by category name..."
                  className="pl-9 bg-muted/50 border-transparent focus:bg-background"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Ref Code</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {categoriesLoading ? (
                  Array.from({ length: 5 }).map((_, i) => (
                    <TableRow key={i}>
                      <TableCell><Skeleton className="h-4 w-40" /></TableCell>
                      <TableCell><Skeleton className="h-4 w-24" /></TableCell>
                      <TableCell><Skeleton className="h-6 w-16" /></TableCell>
                      <TableCell>
                        <div className="flex justify-end gap-2">
                          <Skeleton className="h-8 w-8" />
                          <Skeleton className="h-8 w-8" />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                ) : filteredCategories.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} className="h-40">
                      <div className="flex flex-col items-center justify-center text-center">
                        <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted mb-3">
                          <Layers className="h-6 w-6 text-muted-foreground" />
                        </div>
                        <p className="text-sm font-medium text-foreground">No categories found</p>
                        <p className="text-xs text-muted-foreground mt-1">
                          {search ? "Try adjusting your search" : "Add a category to get started"}
                        </p>
                        {!search && (
                          <Button variant="outline" size="sm" className="mt-3" onClick={handleCreateCategory}>
                            <Plus className="h-3 w-3 mr-1" />
                            Add Category
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  filteredCategories.map((category) => (
                    <TableRow key={category.id}>
                      <TableCell>
                        <p className="font-medium text-foreground">{category.name}</p>
                      </TableCell>
                      <TableCell>
                        <span className="font-mono text-xs font-medium text-muted-foreground">
                          {category.refCode}
                        </span>
                      </TableCell>
                      <TableCell>
                        <Badge variant={category.isActive ? "success" : "secondary"}>
                          {category.isActive ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleEditCategory(category)}
                            className="h-8 w-8"
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleDeleteCategory(category)}
                            className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}

      <ProductFormDialog
        product={selectedProduct}
        open={isFormOpen}
        onOpenChange={setIsFormOpen}
      />

      <ProductDeleteDialog
        product={selectedProduct}
        open={isDeleteOpen}
        onOpenChange={setIsDeleteOpen}
      />

      <ProductCategoryFormDialog
        category={selectedCategory}
        open={isCategoryFormOpen}
        onOpenChange={setIsCategoryFormOpen}
      />

      <ProductCategoryDeleteDialog
        category={selectedCategory}
        open={isCategoryDeleteOpen}
        onOpenChange={setIsCategoryDeleteOpen}
      />
    </div>
  )
}
