import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Loader2 } from 'lucide-react';
import { useCategories } from '../hooks/useAdminProducts';

const adminProductSchema = z.object({
  name: z.string().min(1, 'Product name is required'),
  description: z.string().min(10, 'Description must be at least 10 characters'),
  price: z.coerce.number().min(1, 'Price is required'),
  compare_at_price: z.coerce.number().optional(),
  category_id: z.string().min(1, 'Category is required'),
  stock_quantity: z.coerce.number().min(0, 'Stock must be 0 or more'),
  image_url: z.string().optional(),
});

type AdminProductFormValues = z.infer<typeof adminProductSchema>;

interface AdminProductFormProps {
  defaultValues?: Partial<AdminProductFormValues>;
  onSubmit: (data: AdminProductFormValues) => void;
  isPending?: boolean;
  submitLabel?: string;
}

export function AdminProductForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Save Product',
}: AdminProductFormProps) {
  const { data: categories } = useCategories();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<AdminProductFormValues>({
    resolver: zodResolver(adminProductSchema),
    defaultValues,
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Product Name</Label>
        <Input id="name" {...register('name')} placeholder="e.g. Premium Wireless Headphones" />
        {errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <textarea
          id="description"
          {...register('description')}
          rows={4}
          placeholder="Describe the product..."
          className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
        {errors.description && <p className="text-sm text-destructive">{errors.description.message}</p>}
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        <div className="space-y-2">
          <Label htmlFor="price">Price (cents)</Label>
          <Input id="price" type="number" {...register('price')} placeholder="9999" />
          {errors.price && <p className="text-sm text-destructive">{errors.price.message}</p>}
        </div>
        <div className="space-y-2">
          <Label htmlFor="compare_at_price">Compare at Price</Label>
          <Input id="compare_at_price" type="number" {...register('compare_at_price')} placeholder="Optional" />
        </div>
        <div className="space-y-2">
          <Label htmlFor="stock_quantity">Stock Quantity</Label>
          <Input id="stock_quantity" type="number" {...register('stock_quantity')} placeholder="100" />
          {errors.stock_quantity && <p className="text-sm text-destructive">{errors.stock_quantity.message}</p>}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="category_id">Category</Label>
        <select
          id="category_id"
          {...register('category_id')}
          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <option value="">Select a category...</option>
          {(categories || []).map((cat) => (
            <option key={cat.id} value={cat.id}>
              {cat.name}
            </option>
          ))}
        </select>
        {errors.category_id && <p className="text-sm text-destructive">{errors.category_id.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="image_url">Image URL</Label>
        <Input id="image_url" {...register('image_url')} placeholder="https://images.unsplash.com/..." />
        <p className="text-xs text-muted-foreground">Enter a direct URL to the product image</p>
      </div>

      <Button type="submit" disabled={isPending} className="w-full">
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}
