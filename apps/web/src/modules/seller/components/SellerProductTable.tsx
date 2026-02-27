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
import { Edit, Trash2 } from 'lucide-react';
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
          <TableHead>Price</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Variants</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {products.map((product) => (
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
                  <div className="h-10 w-10 rounded-md bg-muted" />
                )}
                <span className="font-medium">{product.name}</span>
              </div>
            </TableCell>
            <TableCell>{formatPrice(product.base_price_cents)}</TableCell>
            <TableCell>
              <Badge variant={product.status === 'active' ? 'default' : 'secondary'}>
                {product.status}
              </Badge>
            </TableCell>
            <TableCell>
              {product.has_variants ? (
                <Badge variant="outline">{product.variants?.length || 0} variants</Badge>
              ) : (
                <span className="text-muted-foreground text-sm">None</span>
              )}
            </TableCell>
            <TableCell className="text-right">
              <div className="flex justify-end gap-2">
                <Button asChild variant="ghost" size="icon">
                  <Link to={`/seller/products/${product.id}/edit`}>
                    <Edit className="h-4 w-4" />
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
        ))}
      </TableBody>
    </Table>
  );
}
