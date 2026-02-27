import { Link } from 'react-router-dom';
import { Badge } from '@/shared/components/ui/badge';
import { Button } from '@/shared/components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { Edit, Trash2, Settings2, Package } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import type { SellerProduct } from '../services/seller-product.api';

interface SellerProductTableProps {
  products: SellerProduct[];
  onDelete?: (id: string) => void;
}

export function SellerProductTable({ products, onDelete }: SellerProductTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Product</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Price</TableHead>
          <TableHead>Stock / Variants</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {products.map((product) => {
          const isConfigurable = product.product_type === 'configurable';
          return (
            <TableRow key={product.id}>
              <TableCell>
                <div className="flex items-center gap-3">
                  {product.image_urls?.[0] ? (
                    <img
                      src={product.image_urls[0]}
                      alt={product.name}
                      className="h-10 w-10 rounded-md bg-muted object-cover"
                    />
                  ) : (
                    <div className="flex h-10 w-10 items-center justify-center rounded-md bg-muted">
                      <Package className="h-5 w-5 text-muted-foreground/40" />
                    </div>
                  )}
                  <div>
                    <Link
                      to={`/seller/products/${product.id}/edit`}
                      className="font-medium hover:underline"
                    >
                      {product.name}
                    </Link>
                  </div>
                </div>
              </TableCell>
              <TableCell>
                <Badge variant={isConfigurable ? 'default' : 'outline'} className="text-xs">
                  {isConfigurable ? 'Configurable' : 'Simple'}
                </Badge>
              </TableCell>
              <TableCell>{formatPrice(product.base_price_cents)}</TableCell>
              <TableCell>
                {isConfigurable ? (
                  <div className="flex items-center gap-2">
                    <Badge variant="outline">
                      {product.variants?.length || 0} variant{(product.variants?.length || 0) !== 1 ? 's' : ''}
                    </Badge>
                    <Link
                      to={`/seller/products/${product.id}/edit`}
                      className="text-xs text-primary hover:underline"
                    >
                      Manage
                    </Link>
                  </div>
                ) : (
                  <span className="text-sm">
                    {product.stock_quantity ?? 0} in stock
                  </span>
                )}
              </TableCell>
              <TableCell>
                <Badge variant={product.status === 'active' ? 'default' : 'secondary'}>
                  {product.status}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-1">
                  <Button asChild variant="ghost" size="sm">
                    <Link to={`/seller/products/${product.id}/edit`}>
                      <Settings2 className="mr-1 h-4 w-4" />
                      Edit
                    </Link>
                  </Button>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-destructive"
                    onClick={() => onDelete?.(product.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          );
        })}
      </TableBody>
    </Table>
  );
}
