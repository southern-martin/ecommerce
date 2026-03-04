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
import { AttributeGroupTable } from '../components/AttributeGroupTable';
import { AttributeGroupForm } from '../components/AttributeGroupForm';
import { AttributeGroupManager } from '../components/AttributeGroupManager';
import { AdminProductTable, getProductStatus, getProductTags } from '../components/AdminProductTable';
import type { ProductStatus } from '../components/AdminProductTable';
import { AdminProductForm } from '../components/AdminProductForm';
import {
  useCategories,
  useCreateCategory,
  useUpdateCategory,
  useDeleteCategory,
  useAdminAttributes,
  useCreateAttribute,
  useUpdateAttribute,
  useDeleteAttribute,
  useAttributeGroups,
  useCreateAttributeGroup,
  useUpdateAttributeGroup,
  useDeleteAttributeGroup,
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
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [editingCategory, setEditingCategory] = useState<any>(null);

  // Attribute state
  const [attributeDialogOpen, setAttributeDialogOpen] = useState(false);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [editingAttribute, setEditingAttribute] = useState<any>(null);

  // Attribute Group state
  const [attrGroupDialogOpen, setAttrGroupDialogOpen] = useState(false);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [editingAttrGroup, setEditingAttrGroup] = useState<any>(null);
  const [managingGroup, setManagingGroup] = useState<{ id: string; name: string } | null>(null);

  // Product state
  const [productDialogOpen, setProductDialogOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<Product | null>(null);
  const [productPage, setProductPage] = useState(1);
  const [productSearch, setProductSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');

  // Category hooks
  const { data: categories, isLoading: categoriesLoading } = useCategories();
  const createCategory = useCreateCategory();
  const updateCategory = useUpdateCategory();
  const deleteCategory = useDeleteCategory();

  // Attribute hooks
  const { data: attributes, isLoading: attributesLoading } = useAdminAttributes();
  const createAttribute = useCreateAttribute();
  const updateAttribute = useUpdateAttribute();
  const deleteAttribute = useDeleteAttribute();

  // Attribute Group hooks
  const { data: attributeGroups, isLoading: attrGroupsLoading } = useAttributeGroups();
  const createAttrGroup = useCreateAttributeGroup();
  const updateAttrGroup = useUpdateAttributeGroup();
  const deleteAttrGroup = useDeleteAttributeGroup();

  // Product hooks
  const { data: productData, isLoading: productsLoading } = useAdminProductList({
    page: productPage,
    page_size: 20,
    search: productSearch || undefined,
    status: statusFilter || undefined,
  });
  const createProduct = useAdminCreateProduct();
  const updateProduct = useAdminUpdateProduct();
  const deleteProduct = useAdminDeleteProduct();

  const totalProductPages = productData ? Math.ceil(productData.total / 20) : 0;

  // Category handlers
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const handleCreateCategory = (data: any) => {
    createCategory.mutate(data, {
      onSuccess: () => setCategoryDialogOpen(false),
    });
  };

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const handleUpdateCategory = (data: any) => {
    if (!editingCategory) return;
    updateCategory.mutate(
      { id: editingCategory.id, data: { name: data.name, parent_id: data.parent_id || undefined } },
      { onSuccess: () => setEditingCategory(null) }
    );
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

  // Attribute Group handlers
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const handleCreateAttrGroup = (data: any) => {
    createAttrGroup.mutate(data, {
      onSuccess: () => setAttrGroupDialogOpen(false),
    });
  };

  // Product handlers
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const handleCreateProduct = (data: any) => {
    const payload = {
      name: data.name,
      description: data.description,
      base_price_cents: data.base_price_cents,
      currency: data.currency || 'USD',
      category_id: data.category_id,
      attribute_group_id: data.attribute_group_id || undefined,
      product_type: data.product_type || 'simple',
      stock_quantity: data.product_type === 'simple' ? (data.stock_quantity ?? 0) : 0,
      tags: data.tags || [],
      image_urls: data.image_urls || [],
    };
    createProduct.mutate(payload, {
      onSuccess: () => setProductDialogOpen(false),
    });
  };

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const handleUpdateProduct = (data: any) => {
    if (!editingProduct) return;
    updateProduct.mutate(
      {
        id: editingProduct.id,
        data: {
          name: data.name,
          description: data.description,
          base_price_cents: data.base_price_cents,
          currency: data.currency,
          status: data.status,
          category_id: data.category_id,
          attribute_group_id: data.attribute_group_id || undefined,
          tags: data.tags || [],
          image_urls: data.image_urls || [],
        },
      },
      { onSuccess: () => setEditingProduct(null) }
    );
  };

  const handleStatusChange = (id: string, status: ProductStatus) => {
    updateProduct.mutate({
      id,
      data: { status },
    });
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
          <TabsTrigger value="attribute-groups">Attribute Groups</TabsTrigger>
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
              <DialogContent className="sm:max-w-lg">
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
                onStatusChange={handleStatusChange}
                isDeleting={deleteProduct.isPending}
                isUpdating={updateProduct.isPending}
                searchValue={productSearch}
                onSearchChange={setProductSearch}
                statusFilter={statusFilter}
                onStatusFilterChange={setStatusFilter}
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
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>Edit Product</DialogTitle>
              </DialogHeader>
              {editingProduct && (
                <AdminProductForm
                  isEditing
                  defaultValues={{
                    name: editingProduct.name,
                    description: editingProduct.description,
                    base_price_cents: editingProduct.price,
                    currency: (editingProduct as any)._currency || 'USD',
                    category_id: editingProduct.category?.id || '',
                    attribute_group_id: (editingProduct as any)._attribute_group_id || '',
                    status: getProductStatus(editingProduct),
                    tags: getProductTags(editingProduct),
                    image_urls: editingProduct.images?.map((img) => img.url) || [],
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
              onEdit={(cat) => setEditingCategory(cat)}
              onDelete={(id) => deleteCategory.mutate(id)}
              isDeleting={deleteCategory.isPending}
            />
          )}

          {/* Edit Category Dialog */}
          <Dialog
            open={!!editingCategory}
            onOpenChange={(open) => !open && setEditingCategory(null)}
          >
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Category</DialogTitle>
              </DialogHeader>
              {editingCategory && (
                <CategoryForm
                  categories={(categories || []).filter((c) => c.id !== editingCategory.id)}
                  defaultValues={{
                    name: editingCategory.name,
                    slug: editingCategory.slug,
                    description: editingCategory.description || '',
                    parent_id: editingCategory.parent_id || '',
                  }}
                  onSubmit={handleUpdateCategory}
                  isPending={updateCategory.isPending}
                  submitLabel="Update Category"
                />
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>

        {/* Attribute Groups Tab */}
        <TabsContent value="attribute-groups" className="space-y-4">
          <div className="flex items-center justify-between">
            <p className="text-sm text-muted-foreground">
              Group related attributes together and assign them to products
            </p>
            <Dialog open={attrGroupDialogOpen} onOpenChange={setAttrGroupDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Attribute Group
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create Attribute Group</DialogTitle>
                </DialogHeader>
                <AttributeGroupForm
                  onSubmit={handleCreateAttrGroup}
                  isPending={createAttrGroup.isPending}
                />
              </DialogContent>
            </Dialog>
          </div>

          {attrGroupsLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : (
            <AttributeGroupTable
              groups={attributeGroups || []}
              onEdit={(group) => setEditingAttrGroup(group)}
              onManageAttributes={(group) => setManagingGroup({ id: group.id, name: group.name })}
              onDelete={(id) => deleteAttrGroup.mutate(id)}
              isDeleting={deleteAttrGroup.isPending}
            />
          )}

          {/* Attribute Group Manager Dialog */}
          {managingGroup && (
            <AttributeGroupManager
              groupId={managingGroup.id}
              groupName={managingGroup.name}
              open={!!managingGroup}
              onOpenChange={(open) => !open && setManagingGroup(null)}
            />
          )}

          {/* Edit Attribute Group Dialog */}
          <Dialog
            open={!!editingAttrGroup}
            onOpenChange={(open) => !open && setEditingAttrGroup(null)}
          >
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Attribute Group</DialogTitle>
              </DialogHeader>
              {editingAttrGroup && (
                <AttributeGroupForm
                  defaultValues={{
                    name: editingAttrGroup.name,
                    description: editingAttrGroup.description || '',
                  }}
                  onSubmit={(data) => {
                    updateAttrGroup.mutate(
                      { id: editingAttrGroup.id, data: { name: data.name, description: data.description } },
                      { onSuccess: () => setEditingAttrGroup(null) }
                    );
                  }}
                  isPending={updateAttrGroup.isPending}
                  submitLabel="Update Group"
                />
              )}
            </DialogContent>
          </Dialog>
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
                  }}
                  defaultOptionValues={editingAttribute.option_values?.map((ov: { value: string; color_hex?: string }) => ({
                    value: ov.value,
                    color_hex: ov.color_hex,
                  }))}
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
