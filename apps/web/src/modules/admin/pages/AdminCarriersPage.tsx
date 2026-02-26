import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Plus, Pencil } from 'lucide-react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Skeleton } from '@/shared/components/ui/skeleton';
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
import { Badge } from '@/shared/components/ui/badge';
import { useAdminCarriers, useCreateCarrier, useUpdateCarrier } from '../hooks/useAdminShipping';
import type { Carrier } from '../services/admin-shipping.api';

const carrierSchema = z.object({
  code: z.string().min(1, 'Code is required'),
  name: z.string().min(1, 'Name is required'),
  tracking_url_template: z.string().min(1, 'Tracking URL template is required'),
  is_active: z.boolean(),
});

type CarrierFormValues = z.infer<typeof carrierSchema>;

export default function AdminCarriersPage() {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingCarrier, setEditingCarrier] = useState<Carrier | null>(null);

  const { data: carriers, isLoading } = useAdminCarriers();
  const createCarrier = useCreateCarrier();
  const updateCarrier = useUpdateCarrier();

  const form = useForm<CarrierFormValues>({
    resolver: zodResolver(carrierSchema),
    defaultValues: {
      code: '',
      name: '',
      tracking_url_template: '',
      is_active: true,
    },
  });

  const openCreate = () => {
    setEditingCarrier(null);
    form.reset({ code: '', name: '', tracking_url_template: '', is_active: true });
    setDialogOpen(true);
  };

  const openEdit = (carrier: Carrier) => {
    setEditingCarrier(carrier);
    form.reset({
      code: carrier.code,
      name: carrier.name,
      tracking_url_template: carrier.tracking_url_template,
      is_active: carrier.is_active,
    });
    setDialogOpen(true);
  };

  const onSubmit = (values: CarrierFormValues) => {
    if (editingCarrier) {
      updateCarrier.mutate(
        { code: editingCarrier.code, data: values },
        { onSuccess: () => setDialogOpen(false) },
      );
    } else {
      createCarrier.mutate(values, { onSuccess: () => setDialogOpen(false) });
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
        <h1 className="text-2xl font-bold">Shipping Carriers</h1>
        <Button onClick={openCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Add Carrier
        </Button>
      </div>

      {carriers && carriers.length > 0 ? (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Code</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Tracking URL Template</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-20">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {carriers.map((carrier) => (
              <TableRow key={carrier.code}>
                <TableCell className="font-mono text-sm">{carrier.code}</TableCell>
                <TableCell>{carrier.name}</TableCell>
                <TableCell className="max-w-xs truncate text-sm text-muted-foreground">
                  {carrier.tracking_url_template}
                </TableCell>
                <TableCell>
                  <Badge variant={carrier.is_active ? 'default' : 'secondary'}>
                    {carrier.is_active ? 'Active' : 'Inactive'}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Button variant="ghost" size="icon" onClick={() => openEdit(carrier)}>
                    <Pencil className="h-4 w-4" />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      ) : (
        <p className="py-8 text-center text-muted-foreground">No carriers configured.</p>
      )}

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{editingCarrier ? 'Edit Carrier' : 'Add Carrier'}</DialogTitle>
          </DialogHeader>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="code">Code</Label>
              <Input
                id="code"
                {...form.register('code')}
                disabled={!!editingCarrier}
                placeholder="e.g. USPS"
              />
              {form.formState.errors.code && (
                <p className="text-sm text-destructive">{form.formState.errors.code.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" {...form.register('name')} placeholder="e.g. US Postal Service" />
              {form.formState.errors.name && (
                <p className="text-sm text-destructive">{form.formState.errors.name.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="tracking_url_template">Tracking URL Template</Label>
              <Input
                id="tracking_url_template"
                {...form.register('tracking_url_template')}
                placeholder="https://track.example.com/{tracking_number}"
              />
              {form.formState.errors.tracking_url_template && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.tracking_url_template.message}
                </p>
              )}
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
              <Button
                type="submit"
                disabled={createCarrier.isPending || updateCarrier.isPending}
              >
                {editingCarrier ? 'Update' : 'Create'}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
}
