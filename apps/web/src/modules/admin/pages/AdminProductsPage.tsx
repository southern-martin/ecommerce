import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/shared/components/ui/dialog';
import { Plus, ChevronLeft, ChevronRight } from 'lucide-react';
import { CategoryTable } from '../components/CategoryTable';
import { CategoryForm } from '../components/CategoryForm';
import { AttributeTable } from '../components/AttributeTable';
import { AttributeForm } from '../components/AttributeForm';
import { AdminProductTable } from '../components/AdminProductTable';
import { AdminProductForm } from '../components/AdminProductForm';
import {
  useCategories,
  useCreateCategory,
  useAdminAttributes,
  useCreateAttribute,
  useUpdateAttribute,
  useDeleteAttribute,
} from '../hooks/useAdminProducts';
import {
  useAdminProductList,
  useAdminCreateProduct,
  useAdminUpdateProduct,
  useAdminDeleteProduct,
} from '../hooks/useAdminProductMgmt';
import type { CreateAttributeData } from '../services/admin-product.api';
import type { Product } from '@/modules/shop/types/shop.types';

export default function AdminProductsPage() {
  // Category state
  const [categoryDialogOpen, setCategoryDialogOpen] = useState(false);
  // Attribute state
  const [attributeDialogOpen, setAttributeDialogOpen] = useState(false);
  const [editingAttribute, setEditingAttribute] = useState<any>(null);
  // Product state
  const [productDialogOpen, setProductDialogOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [productPage, setProductPage] = useState(1);
  const [productSearch, setProductSearch] = useState('');

  // Category hooks
  const { data: categories, isLoading: categoriesLoading } = useCategories();
  const createCategory = useCreateCategory();

  // Attribute hooks
  const { data: attributes, isLoading: attributesLoading } = useAdminAttributes();
  const createAttribute = useCreateAttribute();
  const updateAttribute = useUpdateAttribute();
  const deleteAttribute = useDeleteAttribute();

  // Product hooks
  const { data: productData, isLoading: productsLoading } = useAdminProductList({
    page: productPage,
    page_size: 20,
    search: productSearch || undefined,
  });
  const createProduct = useAdminCreateProduct();
  const updateProduct = useAdminUpdateProduct();
  const deleteProduct = useAdminDeleteProduct();

  const totalProductPages = productData ? Math.ceil(productData.total / 20) : 0;

  // Category handlers
  const handleCreateCategory = (data: any) => {
    createCategory.mutate(data, {
      onSuccess: () => setCategoryDialogOpen(false),
    });
  };

  // Attribute handlers
  const handleCreateAttribute = (data: CreateAttributeData) => {
    createAttribute.mutate(data, {
      onSuccess: () => setAttributeDialogOpen(false),
    });
  };

  const handleUpdateAttribute = (data: CreateAttributeData) => {
    if (!editingAttribute) return;
    updateAttribute.mutate(
      { id: editingAttribute.id, data },
      { onSuccess: () => setEditingAttribute(null) }
    );
  };

  const handleManageAttributes = (_category: any) => {
    // Placeholder for category attribute management navigation
  };

  // Product handlers
  const handleCreateProduct = (data: any) => {
    const payload = {
      name: data.name,
      description: data.description,
      price: data.price,
      compare_at_price: data.compare_at_price,
      category_id: data.category_id,
      stock_quantity: data.stock_quantity,
      images: data.image_url ? [{ url: data.image_url, alt: data.name, is_primary: true }] : [],
    };
    createProduct.mutate(payload, {
      onSuccess: () => setProductDialogOpen(false),
    });
  };

  const handleUpdateProduct = (data: any) => {
    if (!editingProduct) return;
    const payload: any = {
      name: data.name,
      description: data.description,
      price: data.price,
      compare_at_price: data.compare_at_price,
      category_id: data.category_id,
      stock_quantity: data.stock_quantity,
    };
    if (data.image_url) {
      payload.images = [{ url: data.image_url, alt: data.name, is_primary: true }];
    }
    updateProduct.mutate(
      { id: editingProduct.id, data: payload },
      { onSuccess: () => setEditingProduct(null) }
    );
  };

  const handleDeleteProduct = (id: string) => {
    deleteProduct.mutate(id);
  };

  const handleEditProduct = (product: Product) => {
    setEditingProduct(product);
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Products Management</h1>

      <Tabs defaultValue="products">
        <TabsList>
          <TabsTrigger value="products">Products</TabsTrigger>
          <TabsTrigger value="categories">Categories</TabsTrigger>
          <TabsTrigger value="attributes">Attributes</TabsTrigger>
        </TabsList>

        {/* Products Tab */}
        <TabsContent value="products" className="space-y-4">
          <div className="flex items-center justify-between">
            <p className="text-sm text-muted-foreground">
              {productData ? `${productData.total} products from all vendors` : 'Loading...'}
            </p>
            <Dialog open={productDialogOpen} onOpenChange={setProductDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Add Product
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create Product</DialogTitle>
                </DialogHeader>
                <AdminProductForm
                  onSubmit={handleCreateProduct}
                  isPending={createProduct.isPending}
                  submitLabel="Create Product"
                />
              </DialogContent>
            </Dialog>
          </div>

          {productsLoading ? (
            <div className="space-y-3">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-16 w-full rounded-xl" />
              ))}
            </div>
          ) : (
            <>
              <AdminProductTable
                products={productData?.data ?? []}
                onEdit={handleEditProduct}
                onDelete={handleDeleteProduct}
                isDeleting={deleteProduct.isPending}
                searchValue={productSearch}
                onSearchChange={setProductSearch}
              />

              {totalProductPages > 1 && (
                <div className="flex items-center justify-center gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={productPage === 1}
                    onClick={() => setProductPage((p) => p - 1)}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <span className="text-sm text-muted-foreground">
                    Page {productPage} of {totalProductPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={productPage === totalProductPages}
                    onClick={() => setProductPage((p) => p + 1)}
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              )}
            </>
          )}

          {/* Edit Product Dialog */}
          <Dialog
            open={!!editingProduct}
            onOpenChange={(open) => !open && setEditingProduct(null)}
          >
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Product</DialogTitle>
              </DialogHeader>
              {editingProduct && (
                <AdminProductForm
                  defaultValues={{
                    name: editingProduct.name,
                    description: editingProduct.description,
                    price: editingProduct.price,
                    compare_at_price: editingProduct.compare_at_price,
                    category_id: editingProduct.category?.id || '',
                    stock_quantity: editingProduct.stock_quantity,
                    image_url: editingProduct.images?.[0]?.url || '',
                  }}
                  onSubmit={handleUpdateProduct}
                  isPending={updateProduct.isPending}
                  submitLabel="Update Product"
                />
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>

        {/* Categories Tab */}
        <TabsContent value="categories" className="space-y-4">
          <div className="flex justify-end">
            <Dialog open={categoryDialogOpen} onOpenChange={setCategoryDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Category
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create Category</DialogTitle>
                </DialogHeader>
                <CategoryForm
                  categories={categories || []}
                  onSubmit={handleCreateCategory}
                  isPending={createCategory.isPending}
                />
              </DialogContent>
            </Dialog>
          </div>

          {categoriesLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : (
            <CategoryTable
              categories={categories || []}
              onManageAttributes={handleManageAttributes}
            />
          )}
        </TabsContent>

        {/* Attributes Tab */}
        <TabsContent value="attributes" className="space-y-4">
          <div className="flex justify-end">
            <Dialog open={attributeDialogOpen} onOpenChange={setAttributeDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Attribute
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create Attribute</DialogTitle>
                </DialogHeader>
                <AttributeForm
                  onSubmit={handleCreateAttribute}
                  isPending={createAttribute.isPending}
                />
              </DialogContent>
            </Dialog>
          </div>

          {attributesLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : (
            <AttributeTable
              attributes={attributes || []}
              onEdit={(attr) => setEditingAttribute(attr)}
              onDelete={(id) => deleteAttribute.mutate(id)}
              isDeleting={deleteAttribute.isPending}
            />
          )}

          {/* Edit Attribute Dialog */}
          <Dialog
            open={!!editingAttribute}
            onOpenChange={(open) => !open && setEditingAttribute(null)}
          >
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Attribute</DialogTitle>
              </DialogHeader>
              {editingAttribute && (
                <AttributeForm
                  defaultValues={{
                    name: editingAttribute.name,
                    type: editingAttribute.type,
                    required: editingAttribute.required,
                    filterable: editingAttribute.filterable,
                    options: editingAttribute.options?.join(', ') || '',
                  }}
                  onSubmit={handleUpdateAttribute}
                  isPending={updateAttribute.isPending}
                  submitLabel="Update Attribute"
                />
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>
      </Tabs>
    </div>
  );
}
