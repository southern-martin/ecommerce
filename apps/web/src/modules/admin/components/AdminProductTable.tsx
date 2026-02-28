import { useState } from 'react';
import { Badge } from '@/shared/components/ui/badge';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
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
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog';
import { Link } from 'react-router-dom';
import {
  Edit,
  Trash2,
  Search,
  Package,
  CheckCircle,
  XCircle,
  Archive,
  FileText,
  Eye,
  Settings2,
} from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import type { Product } from '@/modules/shop/types/shop.types';

export type ProductStatus = 'draft' | 'active' | 'inactive' | 'archived';

const STATUS_CONFIG: Record<ProductStatus, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline'; icon: typeof CheckCircle }> = {
  draft: { label: 'Draft', variant: 'secondary', icon: FileText },
  active: { label: 'Active', variant: 'default', icon: CheckCircle },
  inactive: { label: 'Inactive', variant: 'destructive', icon: XCircle },
  archived: { label: 'Archived', variant: 'outline', icon: Archive },
};

const STATUS_FILTERS: { label: string; value: string }[] = [
  { label: 'All', value: '' },
  { label: 'Draft', value: 'draft' },
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Archived', value: 'archived' },
];

interface AdminProductTableProps {
  products: Product[];
  onEdit: (product: Product) => void;
  onDelete: (id: string) => void;
  onStatusChange: (id: string, status: ProductStatus) => void;
  isDeleting?: boolean;
  isUpdating?: boolean;
  searchValue?: string;
  onSearchChange?: (value: string) => void;
  statusFilter?: string;
  onStatusFilterChange?: (status: string) => void;
}

// Helper to get status from the augmented product
export function getProductStatus(product: Product): ProductStatus {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const status = (product as any)._status;
  if (status && ['draft', 'active', 'inactive', 'archived'].includes(status)) {
    return status as ProductStatus;
  }
  return product.in_stock ? 'active' : 'inactive';
}

// Helper to get tags from the augmented product
export function getProductTags(product: Product): string[] {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  return (product as any)._tags || [];
}

// Helper to get product type
export function getProductType(product: Product): string {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  return (product as any)._product_type || product.product_type || 'simple';
}

export function AdminProductTable({
  products,
  onEdit,
  onDelete,
  onStatusChange,
  isDeleting,
  isUpdating,
  searchValue = '',
  onSearchChange,
  statusFilter = '',
  onStatusFilterChange,
}: AdminProductTableProps) {
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const [statusAction, setStatusAction] = useState<{ id: string; status: ProductStatus } | null>(null);

  const handleConfirmDelete = () => {
    if (deleteId) {
      onDelete(deleteId);
      setDeleteId(null);
    }
  };

  const handleConfirmStatusChange = () => {
    if (statusAction) {
      onStatusChange(statusAction.id, statusAction.status);
      setStatusAction(null);
    }
  };

  if (products.length === 0 && !searchValue && !statusFilter) {
    return (
      <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed py-16">
        <div className="flex h-14 w-14 items-center justify-center rounded-full bg-muted">
          <Package className="h-7 w-7 text-muted-foreground" />
        </div>
        <p className="mt-4 text-sm font-medium text-muted-foreground">No products yet</p>
        <p className="mt-1 text-xs text-muted-foreground/70">Create your first product to get started</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Search + Status Filter Bar */}
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        {onSearchChange && (
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder="Search products by name..."
              value={searchValue}
              onChange={(e) => onSearchChange(e.target.value)}
              className="pl-10"
            />
          </div>
        )}

        {onStatusFilterChange && (
          <div className="flex gap-1">
            {STATUS_FILTERS.map((f) => (
              <Button
                key={f.value}
                variant={statusFilter === f.value ? 'default' : 'outline'}
                size="sm"
                onClick={() => onStatusFilterChange(f.value)}
                className="text-xs"
              >
                {f.label}
              </Button>
            ))}
          </div>
        )}
      </div>

      <div className="rounded-xl border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[300px]">Product</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Price</TableHead>
              <TableHead>Seller</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {products.length === 0 ? (
              <TableRow>
                <TableCell colSpan={8} className="h-24 text-center text-muted-foreground">
                  No products match your filters.
                </TableCell>
              </TableRow>
            ) : (
              products.map((product) => {
                const primaryImage = product.images?.find((img) => img.is_primary) ?? product.images?.[0];
                const status = getProductStatus(product);
                const productType = getProductType(product);
                const isConfigurable = productType === 'configurable';
                const statusCfg = STATUS_CONFIG[status];
                const StatusIcon = statusCfg.icon;

                return (
                  <TableRow key={product.id}>
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <div className="h-12 w-12 flex-shrink-0 overflow-hidden rounded-lg bg-muted">
                          {primaryImage?.url ? (
                            <img
                              src={primaryImage.url}
                              alt={product.name}
                              className="h-full w-full object-cover"
                            />
                          ) : (
                            <div className="flex h-full w-full items-center justify-center">
                              <Package className="h-5 w-5 text-muted-foreground/40" />
                            </div>
                          )}
                        </div>
                        <div className="min-w-0">
                          <Link
                            to={`/admin/products/${product.id}/edit`}
                            className="truncate font-medium hover:underline"
                          >
                            {product.name}
                          </Link>
                          <p className="truncate text-xs text-muted-foreground">{product.slug}</p>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant={isConfigurable ? 'default' : 'outline'} className="text-xs">
                        {isConfigurable ? 'Configurable' : 'Simple'}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-medium">{formatPrice(product.price)}</TableCell>
                    <TableCell>
                      <span className="inline-flex items-center gap-1 rounded-md bg-muted px-2 py-1 text-xs font-mono">
                        {product.seller?.id ? product.seller.id.slice(0, 8) + '...' : 'N/A'}
                      </span>
                    </TableCell>
                    <TableCell>
                      <span className="text-sm text-muted-foreground">
                        {product.category?.name || product.category?.id?.slice(0, 8) || '\u2014'}
                      </span>
                    </TableCell>
                    <TableCell>
                      <Badge variant={statusCfg.variant} className="inline-flex items-center gap-1 text-xs">
                        <StatusIcon className="h-3 w-3" />
                        {statusCfg.label}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-1">
                        {/* Quick status actions */}
                        {status === 'draft' && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-green-600 hover:text-green-700"
                            onClick={() => setStatusAction({ id: product.id, status: 'active' })}
                            title="Activate product"
                            disabled={isUpdating}
                          >
                            <CheckCircle className="h-4 w-4" />
                          </Button>
                        )}
                        {status === 'active' && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-orange-600 hover:text-orange-700"
                            onClick={() => setStatusAction({ id: product.id, status: 'inactive' })}
                            title="Deactivate product"
                            disabled={isUpdating}
                          >
                            <XCircle className="h-4 w-4" />
                          </Button>
                        )}
                        {status === 'inactive' && (
                          <>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-8 w-8 text-green-600 hover:text-green-700"
                              onClick={() => setStatusAction({ id: product.id, status: 'active' })}
                              title="Reactivate product"
                              disabled={isUpdating}
                            >
                              <CheckCircle className="h-4 w-4" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-8 w-8 text-muted-foreground"
                              onClick={() => setStatusAction({ id: product.id, status: 'archived' })}
                              title="Archive product"
                              disabled={isUpdating}
                            >
                              <Archive className="h-4 w-4" />
                            </Button>
                          </>
                        )}
                        <Button asChild variant="ghost" size="sm" className="h-8">
                          <Link to={`/admin/products/${product.id}/edit`} title="Manage product details, options, variants & attributes">
                            <Settings2 className="mr-1 h-4 w-4" />
                            Manage
                          </Link>
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          onClick={() => window.open(`/products/${product.slug}`, '_blank')}
                          title="View in store"
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          onClick={() => onEdit(product)}
                          title="Quick edit"
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 text-destructive hover:text-destructive"
                          onClick={() => setDeleteId(product.id)}
                          title="Delete product"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </div>

      {/* Status Change Confirmation Dialog */}
      <Dialog open={!!statusAction} onOpenChange={(open) => !open && setStatusAction(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Change Product Status</DialogTitle>
            <DialogDescription>
              {statusAction && (
                <>
                  Are you sure you want to change this product&apos;s status to{' '}
                  <Badge variant={STATUS_CONFIG[statusAction.status].variant} className="mx-1">
                    {STATUS_CONFIG[statusAction.status].label}
                  </Badge>
                  ?
                  {statusAction.status === 'active' && (
                    <span className="mt-2 block text-sm">
                      This will make the product visible to customers in the storefront.
                    </span>
                  )}
                  {statusAction.status === 'inactive' && (
                    <span className="mt-2 block text-sm">
                      This will hide the product from the storefront.
                    </span>
                  )}
                  {statusAction.status === 'archived' && (
                    <span className="mt-2 block text-sm">
                      Archived products are hidden and cannot be purchased.
                    </span>
                  )}
                </>
              )}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setStatusAction(null)}>
              Cancel
            </Button>
            <Button onClick={handleConfirmStatusChange} disabled={isUpdating}>
              {isUpdating ? 'Updating...' : 'Confirm'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={!!deleteId} onOpenChange={(open) => !open && setDeleteId(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete Product</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this product? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteId(null)}>
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
              disabled={isDeleting}
            >
              {isDeleting ? 'Deleting...' : 'Delete'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
