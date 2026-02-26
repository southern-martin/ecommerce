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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { Plus, Pencil, Trash2 } from 'lucide-react';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { ConfirmDialog } from '@/shared/components/data/ConfirmDialog';
import { AdminCouponForm } from '../components/AdminCouponForm';
import { FlashSaleForm } from '../components/FlashSaleForm';
import { BundleForm } from '../components/BundleForm';
import {
  useAdminCoupons,
  useCreateCoupon,
  useUpdateCoupon,
  useDeleteCoupon,
  useAdminFlashSales,
  useCreateFlashSale,
  useUpdateFlashSale,
  useAdminBundles,
  useCreateBundle,
  useUpdateBundle,
  useDeleteBundle,
} from '../hooks/useAdminPromotions';
import { formatDate } from '@/shared/lib/utils';

export default function AdminPromotionsPage() {
  const [couponDialogOpen, setCouponDialogOpen] = useState(false);
  const [editingCoupon, setEditingCoupon] = useState<any>(null);
  const [flashSaleDialogOpen, setFlashSaleDialogOpen] = useState(false);
  const [editingFlashSale, setEditingFlashSale] = useState<any>(null);
  const [bundleDialogOpen, setBundleDialogOpen] = useState(false);
  const [editingBundle, setEditingBundle] = useState<any>(null);

  const { data: coupons, isLoading: couponsLoading } = useAdminCoupons();
  const createCoupon = useCreateCoupon();
  const updateCoupon = useUpdateCoupon();
  const deleteCoupon = useDeleteCoupon();

  const { data: flashSales, isLoading: flashSalesLoading } = useAdminFlashSales();
  const createFlashSale = useCreateFlashSale();
  const updateFlashSale = useUpdateFlashSale();

  const { data: bundles, isLoading: bundlesLoading } = useAdminBundles();
  const createBundle = useCreateBundle();
  const updateBundle = useUpdateBundle();
  const deleteBundle = useDeleteBundle();

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Promotions</h1>

      <Tabs defaultValue="coupons">
        <TabsList>
          <TabsTrigger value="coupons">Coupons</TabsTrigger>
          <TabsTrigger value="flash-sales">Flash Sales</TabsTrigger>
          <TabsTrigger value="bundles">Bundles</TabsTrigger>
        </TabsList>

        {/* Coupons Tab */}
        <TabsContent value="coupons" className="space-y-4">
          <div className="flex justify-end">
            <Dialog open={couponDialogOpen} onOpenChange={setCouponDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Coupon
                </Button>
              </DialogTrigger>
              <DialogContent className="max-w-lg">
                <DialogHeader>
                  <DialogTitle>Create Coupon</DialogTitle>
                </DialogHeader>
                <AdminCouponForm
                  onSubmit={(data) =>
                    createCoupon.mutate(data, {
                      onSuccess: () => setCouponDialogOpen(false),
                    })
                  }
                  isPending={createCoupon.isPending}
                />
              </DialogContent>
            </Dialog>
          </div>

          {couponsLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Code</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Value</TableHead>
                  <TableHead>Used</TableHead>
                  <TableHead>Expires</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[100px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(!coupons || coupons.length === 0) ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center text-muted-foreground">
                      No coupons found.
                    </TableCell>
                  </TableRow>
                ) : (
                  coupons.map((coupon) => (
                    <TableRow key={coupon.id}>
                      <TableCell className="font-mono font-medium">{coupon.code}</TableCell>
                      <TableCell className="capitalize">
                        {coupon.type.replace(/_/g, ' ')}
                      </TableCell>
                      <TableCell>{coupon.value}</TableCell>
                      <TableCell>
                        {coupon.used_count}
                        {coupon.max_uses ? `/${coupon.max_uses}` : ''}
                      </TableCell>
                      <TableCell>{formatDate(coupon.expires_at)}</TableCell>
                      <TableCell>
                        <StatusBadge status={coupon.is_active ? 'active' : 'inactive'} />
                      </TableCell>
                      <TableCell>
                        <div className="flex gap-1">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setEditingCoupon(coupon)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <ConfirmDialog
                            title="Delete Coupon"
                            description={`Delete coupon "${coupon.code}"? This cannot be undone.`}
                            onConfirm={() => deleteCoupon.mutate(coupon.id)}
                            isPending={deleteCoupon.isPending}
                            trigger={
                              <Button variant="ghost" size="sm">
                                <Trash2 className="h-4 w-4 text-destructive" />
                              </Button>
                            }
                          />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          )}

          {/* Edit Coupon Dialog */}
          <Dialog open={!!editingCoupon} onOpenChange={(open) => !open && setEditingCoupon(null)}>
            <DialogContent className="max-w-lg">
              <DialogHeader>
                <DialogTitle>Edit Coupon</DialogTitle>
              </DialogHeader>
              {editingCoupon && (
                <AdminCouponForm
                  defaultValues={editingCoupon}
                  onSubmit={(data) =>
                    updateCoupon.mutate(
                      { id: editingCoupon.id, data },
                      { onSuccess: () => setEditingCoupon(null) }
                    )
                  }
                  isPending={updateCoupon.isPending}
                  submitLabel="Update Coupon"
                />
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>

        {/* Flash Sales Tab */}
        <TabsContent value="flash-sales" className="space-y-4">
          <div className="flex justify-end">
            <Dialog open={flashSaleDialogOpen} onOpenChange={setFlashSaleDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Flash Sale
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create Flash Sale</DialogTitle>
                </DialogHeader>
                <FlashSaleForm
                  onSubmit={(data) =>
                    createFlashSale.mutate(data, {
                      onSuccess: () => setFlashSaleDialogOpen(false),
                    })
                  }
                  isPending={createFlashSale.isPending}
                />
              </DialogContent>
            </Dialog>
          </div>

          {flashSalesLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Discount</TableHead>
                  <TableHead>Starts</TableHead>
                  <TableHead>Ends</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[80px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(!flashSales || flashSales.length === 0) ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-muted-foreground">
                      No flash sales found.
                    </TableCell>
                  </TableRow>
                ) : (
                  flashSales.map((sale) => (
                    <TableRow key={sale.id}>
                      <TableCell className="font-medium">{sale.name}</TableCell>
                      <TableCell>{sale.discount_percentage}%</TableCell>
                      <TableCell>{formatDate(sale.starts_at)}</TableCell>
                      <TableCell>{formatDate(sale.ends_at)}</TableCell>
                      <TableCell>
                        <StatusBadge status={sale.is_active ? 'active' : 'inactive'} />
                      </TableCell>
                      <TableCell>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => setEditingFlashSale(sale)}
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          )}

          {/* Edit Flash Sale Dialog */}
          <Dialog
            open={!!editingFlashSale}
            onOpenChange={(open) => !open && setEditingFlashSale(null)}
          >
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Flash Sale</DialogTitle>
              </DialogHeader>
              {editingFlashSale && (
                <FlashSaleForm
                  defaultValues={editingFlashSale}
                  onSubmit={(data) =>
                    updateFlashSale.mutate(
                      { id: editingFlashSale.id, data },
                      { onSuccess: () => setEditingFlashSale(null) }
                    )
                  }
                  isPending={updateFlashSale.isPending}
                  submitLabel="Update Flash Sale"
                />
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>

        {/* Bundles Tab */}
        <TabsContent value="bundles" className="space-y-4">
          <div className="flex justify-end">
            <Dialog open={bundleDialogOpen} onOpenChange={setBundleDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Bundle
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create Bundle</DialogTitle>
                </DialogHeader>
                <BundleForm
                  onSubmit={(data) =>
                    createBundle.mutate(data, {
                      onSuccess: () => setBundleDialogOpen(false),
                    })
                  }
                  isPending={createBundle.isPending}
                />
              </DialogContent>
            </Dialog>
          </div>

          {bundlesLoading ? (
            <Skeleton className="h-64 w-full" />
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Discount</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[100px]">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {(!bundles || bundles.length === 0) ? (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center text-muted-foreground">
                      No bundles found.
                    </TableCell>
                  </TableRow>
                ) : (
                  bundles.map((bundle) => (
                    <TableRow key={bundle.id}>
                      <TableCell className="font-medium">{bundle.name}</TableCell>
                      <TableCell className="max-w-[200px] truncate">
                        {bundle.description || '-'}
                      </TableCell>
                      <TableCell>{bundle.discount_percentage}%</TableCell>
                      <TableCell>
                        <StatusBadge status={bundle.is_active ? 'active' : 'inactive'} />
                      </TableCell>
                      <TableCell>
                        <div className="flex gap-1">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setEditingBundle(bundle)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <ConfirmDialog
                            title="Delete Bundle"
                            description={`Delete bundle "${bundle.name}"? This cannot be undone.`}
                            onConfirm={() => deleteBundle.mutate(bundle.id)}
                            isPending={deleteBundle.isPending}
                            trigger={
                              <Button variant="ghost" size="sm">
                                <Trash2 className="h-4 w-4 text-destructive" />
                              </Button>
                            }
                          />
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          )}

          {/* Edit Bundle Dialog */}
          <Dialog
            open={!!editingBundle}
            onOpenChange={(open) => !open && setEditingBundle(null)}
          >
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Edit Bundle</DialogTitle>
              </DialogHeader>
              {editingBundle && (
                <BundleForm
                  defaultValues={editingBundle}
                  onSubmit={(data) =>
                    updateBundle.mutate(
                      { id: editingBundle.id, data },
                      { onSuccess: () => setEditingBundle(null) }
                    )
                  }
                  isPending={updateBundle.isPending}
                  submitLabel="Update Bundle"
                />
              )}
            </DialogContent>
          </Dialog>
        </TabsContent>
      </Tabs>
    </div>
  );
}
