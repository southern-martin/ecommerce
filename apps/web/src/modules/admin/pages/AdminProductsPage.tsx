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
import { Plus } from 'lucide-react';
import { CategoryTable } from '../components/CategoryTable';
import { CategoryForm } from '../components/CategoryForm';
import { AttributeTable } from '../components/AttributeTable';
import { AttributeForm } from '../components/AttributeForm';
import {
  useCategories,
  useCreateCategory,
  useAdminAttributes,
  useCreateAttribute,
  useUpdateAttribute,
  useDeleteAttribute,
} from '../hooks/useAdminProducts';
import type { CreateAttributeData } from '../services/admin-product.api';

export default function AdminProductsPage() {
  const [categoryDialogOpen, setCategoryDialogOpen] = useState(false);
  const [attributeDialogOpen, setAttributeDialogOpen] = useState(false);
  const [editingAttribute, setEditingAttribute] = useState<any>(null);

  const { data: categories, isLoading: categoriesLoading } = useCategories();
  const createCategory = useCreateCategory();

  const { data: attributes, isLoading: attributesLoading } = useAdminAttributes();
  const createAttribute = useCreateAttribute();
  const updateAttribute = useUpdateAttribute();
  const deleteAttribute = useDeleteAttribute();

  const handleCreateCategory = (data: any) => {
    createCategory.mutate(data, {
      onSuccess: () => setCategoryDialogOpen(false),
    });
  };

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

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Products Management</h1>

      <Tabs defaultValue="categories">
        <TabsList>
          <TabsTrigger value="categories">Categories</TabsTrigger>
          <TabsTrigger value="attributes">Attributes</TabsTrigger>
        </TabsList>

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
