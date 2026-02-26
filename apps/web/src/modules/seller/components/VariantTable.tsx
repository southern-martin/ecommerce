import { useState } from 'react';
import { Save } from 'lucide-react';
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
import { Badge } from '@/shared/components/ui/badge';
import { useUpdateVariant, useUpdateVariantStock } from '../hooks/useSellerVariants';
import { formatPrice } from '@/shared/lib/utils';
import type { Variant } from '../services/seller-variant.api';

interface VariantTableProps {
  productId: string;
  variants: Variant[];
}

export function VariantTable({ productId, variants }: VariantTableProps) {
  const [editingPrices, setEditingPrices] = useState<Record<string, string>>({});
  const [editingStocks, setEditingStocks] = useState<Record<string, string>>({});
  const updateVariant = useUpdateVariant();
  const updateStock = useUpdateVariantStock();

  const handleSaveVariant = (variant: Variant) => {
    const priceCents = editingPrices[variant.id]
      ? Math.round(parseFloat(editingPrices[variant.id]) * 100)
      : variant.price_cents;
    const stock = editingStocks[variant.id]
      ? parseInt(editingStocks[variant.id], 10)
      : variant.stock;

    if (priceCents !== variant.price_cents) {
      updateVariant.mutate({ productId, variantId: variant.id, data: { price_cents: priceCents } });
    }
    if (stock !== variant.stock) {
      updateStock.mutate({ productId, variantId: variant.id, stock });
    }
  };

  if (variants.length === 0) {
    return <p className="text-sm text-muted-foreground">No variants generated yet. Add options above and click &quot;Generate Variants&quot;.</p>;
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Variant</TableHead>
          <TableHead>SKU</TableHead>
          <TableHead>Price</TableHead>
          <TableHead>Stock</TableHead>
          <TableHead>Status</TableHead>
          <TableHead className="w-16" />
        </TableRow>
      </TableHeader>
      <TableBody>
        {variants.map((variant) => (
          <TableRow key={variant.id}>
            <TableCell className="font-medium">{variant.name}</TableCell>
            <TableCell className="text-muted-foreground">{variant.sku}</TableCell>
            <TableCell>
              <Input
                type="number"
                step="0.01"
                className="w-24"
                defaultValue={(variant.price_cents / 100).toFixed(2)}
                onChange={(e) => setEditingPrices((p) => ({ ...p, [variant.id]: e.target.value }))}
              />
            </TableCell>
            <TableCell>
              <Input
                type="number"
                className="w-20"
                defaultValue={variant.stock}
                onChange={(e) => setEditingStocks((p) => ({ ...p, [variant.id]: e.target.value }))}
              />
            </TableCell>
            <TableCell>
              <Badge variant={variant.is_active ? 'default' : 'secondary'}>
                {variant.is_active ? 'Active' : 'Inactive'}
              </Badge>
            </TableCell>
            <TableCell>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleSaveVariant(variant)}
                disabled={updateVariant.isPending || updateStock.isPending}
              >
                <Save className="h-4 w-4" />
              </Button>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
