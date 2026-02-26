import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/shared/components/ui/dialog';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ChevronLeft, ChevronRight, Plus } from 'lucide-react';
import { CouponForm } from '../components/CouponForm';
import { CouponTable } from '../components/CouponTable';
import {
  useSellerCoupons,
  useCreateCoupon,
  useUpdateCoupon,
  useDeleteCoupon,
} from '../hooks/useSellerCoupons';
import type { Coupon } from '../services/seller-coupon.api';

export default function SellerCouponsPage() {
  const [page, setPage] = useState(1);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingCoupon, setEditingCoupon] = useState<Coupon | null>(null);

  const { data, isLoading } = useSellerCoupons(page);
  const createCoupon = useCreateCoupon();
  const updateCoupon = useUpdateCoupon();
  const deleteCoupon = useDeleteCoupon();

  const totalPages = data ? Math.ceil(data.total / data.page_size) : 0;
  const coupons: Coupon[] = data?.data ?? [];

  const handleOpenCreate = () => {
    setEditingCoupon(null);
    setDialogOpen(true);
  };

  const handleEdit = (coupon: Coupon) => {
    setEditingCoupon(coupon);
    setDialogOpen(true);
  };

  const handleSubmit = (formData: Partial<Coupon>) => {
    if (editingCoupon) {
      updateCoupon.mutate(
        { id: editingCoupon.id, data: formData },
        { onSuccess: () => setDialogOpen(false) }
      );
    } else {
      createCoupon.mutate(formData, {
        onSuccess: () => setDialogOpen(false),
      });
    }
  };

  const handleDelete = (id: string) => {
    deleteCoupon.mutate(id);
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
        <h1 className="text-2xl font-bold">Coupons</h1>
        <Button onClick={handleOpenCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Create Coupon
        </Button>
      </div>

      {coupons.length > 0 ? (
        <>
          <CouponTable
            coupons={coupons}
            onEdit={handleEdit}
            onDelete={handleDelete}
            isDeleting={deleteCoupon.isPending}
          />
          {totalPages > 1 && (
            <div className="mt-6 flex items-center justify-center gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={page === 1}
                onClick={() => setPage((p) => p - 1)}
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <span className="text-sm text-muted-foreground">
                Page {page} of {totalPages}
              </span>
              <Button
                variant="outline"
                size="sm"
                disabled={page === totalPages}
                onClick={() => setPage((p) => p + 1)}
              >
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          )}
        </>
      ) : (
        <p className="py-8 text-center text-muted-foreground">No coupons yet.</p>
      )}

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>{editingCoupon ? 'Edit Coupon' : 'Create Coupon'}</DialogTitle>
          </DialogHeader>
          <CouponForm
            onSubmit={handleSubmit}
            isPending={createCoupon.isPending || updateCoupon.isPending}
            defaultValues={editingCoupon ?? undefined}
            submitLabel={editingCoupon ? 'Update Coupon' : 'Create Coupon'}
          />
        </DialogContent>
      </Dialog>
    </div>
  );
}
