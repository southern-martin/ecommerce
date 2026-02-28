import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Loader2, Package, Settings2 } from 'lucide-react';
import { sellerProductApi } from '../services/seller-product.api';

const productSchema = z.object({
  name: z.string().min(1, 'Product name is required'),
  description: z.string().min(10, 'Description must be at least 10 characters'),
  price: z.coerce.number().min(1, 'Price is required'),
  compare_at_price: z.coerce.number().optional(),
  category_id: z.string().min(1, 'Category is required'),
  attribute_group_id: z.string().optional(),
  product_type: z.string().default('simple'),
  stock_quantity: z.coerce.number().min(0, 'Stock must be 0 or more').default(0),
});

type ProductFormValues = z.infer<typeof productSchema>;

interface ProductFormProps {
  defaultValues?: Partial<ProductFormValues>;
  onSubmit: (data: ProductFormValues) => void;
  isPending?: boolean;
  submitLabel?: string;
  showProductTypeSelector?: boolean;
}

export function ProductForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Save Product',
  showProductTypeSelector = false,
}: ProductFormProps) {
  const { data: attributeGroups = [] } = useQuery({
    queryKey: ['attribute-groups'],
    queryFn: () => sellerProductApi.getAttributeGroups(),
  });

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<ProductFormValues>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      product_type: 'simple',
      stock_quantity: 0,
      ...defaultValues,
    },
  });

  const productType = watch('product_type');

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {/* Product Type Selector */}
      {showProductTypeSelector && (
        <div className="space-y-2">
          <Label>Product Type</Label>
          <div className="grid grid-cols-2 gap-3">
            <button
              type="button"
              onClick={() => setValue('product_type', 'simple')}
              className={`flex items-start gap-3 rounded-lg border-2 p-4 text-left transition-all ${
                productType === 'simple'
                  ? 'border-primary bg-primary/5'
                  : 'border-muted hover:border-muted-foreground/30'
              }`}
            >
              <Package className={`mt-0.5 h-5 w-5 flex-shrink-0 ${productType === 'simple' ? 'text-primary' : 'text-muted-foreground'}`} />
              <div>
                <p className="font-medium">Simple Product</p>
                <p className="text-xs text-muted-foreground">
                  A single product with its own price and stock
                </p>
              </div>
            </button>
            <button
              type="button"
              onClick={() => setValue('product_type', 'configurable')}
              className={`flex items-start gap-3 rounded-lg border-2 p-4 text-left transition-all ${
                productType === 'configurable'
                  ? 'border-primary bg-primary/5'
                  : 'border-muted hover:border-muted-foreground/30'
              }`}
            >
              <Settings2 className={`mt-0.5 h-5 w-5 flex-shrink-0 ${productType === 'configurable' ? 'text-primary' : 'text-muted-foreground'}`} />
              <div>
                <p className="font-medium">Configurable Product</p>
                <p className="text-xs text-muted-foreground">
                  Product with options (size, color) and variants
                </p>
              </div>
            </button>
          </div>
          <input type="hidden" {...register('product_type')} />
        </div>
      )}

      {/* Show current type badge when not selectable */}
      {!showProductTypeSelector && defaultValues?.product_type && (
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span>Type:</span>
          <span className="rounded-md bg-muted px-2 py-0.5 font-medium capitalize">
            {defaultValues.product_type}
          </span>
        </div>
      )}

      <div className="space-y-2">
        <Label htmlFor="name">Product Name</Label>
        <Input id="name" {...register('name')} />
        {errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <textarea
          id="description"
          {...register('description')}
          rows={4}
          className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
        {errors.description && <p className="text-sm text-destructive">{errors.description.message}</p>}
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        <div className="space-y-2">
          <Label htmlFor="price">Price ($)</Label>
          <Input id="price" type="number" step="0.01" {...register('price')} />
          {errors.price && <p className="text-sm text-destructive">{errors.price.message}</p>}
        </div>
        <div className="space-y-2">
          <Label htmlFor="compare_at_price">Compare at Price</Label>
          <Input id="compare_at_price" type="number" step="0.01" {...register('compare_at_price')} />
        </div>
        {productType === 'simple' ? (
          <div className="space-y-2">
            <Label htmlFor="stock_quantity">Stock Quantity</Label>
            <Input id="stock_quantity" type="number" {...register('stock_quantity')} />
            {errors.stock_quantity && <p className="text-sm text-destructive">{errors.stock_quantity.message}</p>}
          </div>
        ) : (
          <div className="space-y-2">
            <Label className="text-muted-foreground">Stock</Label>
            <p className="flex h-10 items-center text-sm text-muted-foreground">
              Managed per variant
            </p>
          </div>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="category_id">Category ID</Label>
        <Input id="category_id" {...register('category_id')} />
        {errors.category_id && <p className="text-sm text-destructive">{errors.category_id.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="attribute_group_id">Attribute Group</Label>
        <select
          id="attribute_group_id"
          {...register('attribute_group_id')}
          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <option value="">No attribute group</option>
          {attributeGroups.map((group) => (
            <option key={group.id} value={group.id}>
              {group.name}
            </option>
          ))}
        </select>
        <p className="text-xs text-muted-foreground">
          Determines which specification attributes are available for this product.
        </p>
      </div>

      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}
