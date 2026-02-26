import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Plus, Pencil, Trash2 } from 'lucide-react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Badge } from '@/shared/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog';
import {
  useAdminTaxRules,
  useCreateTaxRule,
  useUpdateTaxRule,
  useDeleteTaxRule,
} from '../hooks/useAdminTax';
import type { TaxRule } from '../services/admin-tax.api';

const taxRuleSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  country: z.string().min(1, 'Country is required'),
  state: z.string(),
  tax_rate: z.coerce.number().min(0).max(100),
  product_category: z.string().min(1, 'Product category is required'),
  is_active: z.boolean(),
});

type TaxRuleFormValues = z.infer<typeof taxRuleSchema>;

export default function AdminTaxPage() {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingRule, setEditingRule] = useState<TaxRule | null>(null);
  const [deleteId, setDeleteId] = useState<string | null>(null);

  const { data: rules, isLoading } = useAdminTaxRules();
  const createRule = useCreateTaxRule();
  const updateRule = useUpdateTaxRule();
  const deleteRule = useDeleteTaxRule();

  const form = useForm<TaxRuleFormValues>({
    resolver: zodResolver(taxRuleSchema),
    defaultValues: {
      name: '',
      country: '',
      state: '',
      tax_rate: 0,
      product_category: '',
      is_active: true,
    },
  });

  const openCreate = () => {
    setEditingRule(null);
    form.reset({
      name: '',
      country: '',
      state: '',
      tax_rate: 0,
      product_category: '',
      is_active: true,
    });
    setDialogOpen(true);
  };

  const openEdit = (rule: TaxRule) => {
    setEditingRule(rule);
    form.reset({
      name: rule.name,
      country: rule.country,
      state: rule.state,
      tax_rate: rule.tax_rate,
      product_category: rule.product_category,
      is_active: rule.is_active,
    });
    setDialogOpen(true);
  };

  const onSubmit = (values: TaxRuleFormValues) => {
    if (editingRule) {
      updateRule.mutate(
        { id: editingRule.id, data: values },
        { onSuccess: () => setDialogOpen(false) },
      );
    } else {
      createRule.mutate(values, { onSuccess: () => setDialogOpen(false) });
    }
  };

  const handleDelete = () => {
    if (deleteId) {
      deleteRule.mutate(deleteId, { onSuccess: () => setDeleteId(null) });
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold">Tax Rules</h1>
        <Button onClick={openCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Add Tax Rule
        </Button>
      </div>

      {rules && rules.length > 0 ? (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Country</TableHead>
              <TableHead>State</TableHead>
              <TableHead>Rate</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-24">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rules.map((rule) => (
              <TableRow key={rule.id}>
                <TableCell className="font-medium">{rule.name}</TableCell>
                <TableCell>{rule.country}</TableCell>
                <TableCell>{rule.state || '-'}</TableCell>
                <TableCell>{rule.tax_rate}%</TableCell>
                <TableCell>{rule.product_category}</TableCell>
                <TableCell>
                  <Badge variant={rule.is_active ? 'default' : 'secondary'}>
                    {rule.is_active ? 'Active' : 'Inactive'}
                  </Badge>
                </TableCell>
                <TableCell>
                  <div className="flex gap-1">
                    <Button variant="ghost" size="icon" onClick={() => openEdit(rule)}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" onClick={() => setDeleteId(rule.id)}>
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      ) : (
        <p className="py-8 text-center text-muted-foreground">No tax rules configured.</p>
      )}

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editingRule ? 'Edit Tax Rule' : 'Add Tax Rule'}</DialogTitle>
          </DialogHeader>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" {...form.register('name')} placeholder="e.g. US Sales Tax" />
              {form.formState.errors.name && (
                <p className="text-sm text-destructive">{form.formState.errors.name.message}</p>
              )}
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="country">Country</Label>
                <Input id="country" {...form.register('country')} placeholder="e.g. US" />
                {form.formState.errors.country && (
                  <p className="text-sm text-destructive">
                    {form.formState.errors.country.message}
                  </p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="state">State</Label>
                <Input id="state" {...form.register('state')} placeholder="e.g. CA" />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="tax_rate">Tax Rate (%)</Label>
                <Input
                  id="tax_rate"
                  type="number"
                  step="0.01"
                  {...form.register('tax_rate')}
                />
                {form.formState.errors.tax_rate && (
                  <p className="text-sm text-destructive">
                    {form.formState.errors.tax_rate.message}
                  </p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="product_category">Product Category</Label>
                <Input
                  id="product_category"
                  {...form.register('product_category')}
                  placeholder="e.g. electronics"
                />
                {form.formState.errors.product_category && (
                  <p className="text-sm text-destructive">
                    {form.formState.errors.product_category.message}
                  </p>
                )}
              </div>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="is_active"
                {...form.register('is_active')}
                className="h-4 w-4 rounded border-gray-300"
              />
              <Label htmlFor="is_active">Active</Label>
            </div>
            <div className="flex justify-end gap-2">
              <Button type="button" variant="outline" onClick={() => setDialogOpen(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={createRule.isPending || updateRule.isPending}>
                {editingRule ? 'Update' : 'Create'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Tax Rule</DialogTitle>
          </DialogHeader>
          <p className="text-sm text-muted-foreground">
            Are you sure you want to delete this tax rule? This action cannot be undone.
          </p>
          <div className="flex justify-end gap-2">
            <Button variant="outline" onClick={() => setDeleteId(null)} disabled={deleteRule.isPending}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDelete} disabled={deleteRule.isPending}>
              {deleteRule.isPending ? 'Deleting...' : 'Delete'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
