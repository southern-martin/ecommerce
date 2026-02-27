import { useState } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Badge } from '@/shared/components/ui/badge';
import { Loader2, Plus, X, ImagePlus, Package, Settings2 } from 'lucide-react';
import { useCategories } from '../hooks/useAdminProducts';

const adminProductSchema = z.object({
  name: z.string().min(1, 'Product name is required'),
  description: z.string().min(10, 'Description must be at least 10 characters'),
  base_price_cents: z.coerce.number().min(1, 'Price is required'),
  currency: z.string().optional(),
  category_id: z.string().min(1, 'Category is required'),
  status: z.enum(['draft', 'active', 'inactive', 'archived']).optional(),
  product_type: z.string().optional(),
  stock_quantity: z.coerce.number().min(0).optional(),
});

type AdminProductFormValues = z.infer<typeof adminProductSchema>;

interface AdminProductFormProps {
  defaultValues?: Partial<AdminProductFormValues> & {
    tags?: string[];
    image_urls?: string[];
    product_type?: string;
    stock_quantity?: number;
  };
  onSubmit: (data: AdminProductFormValues & { tags: string[]; image_urls: string[]; product_type: string; stock_quantity: number }) => void;
  isPending?: boolean;
  submitLabel?: string;
  isEditing?: boolean;
}

export function AdminProductForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Save Product',
  isEditing = false,
}: AdminProductFormProps) {
  const { data: categories } = useCategories();
  const [tags, setTags] = useState<string[]>(defaultValues?.tags || []);
  const [tagInput, setTagInput] = useState('');
  const [imageUrls, setImageUrls] = useState<string[]>(defaultValues?.image_urls || []);
  const [imageInput, setImageInput] = useState('');

  const {
    register,
    handleSubmit,
    control,
    setValue,
    formState: { errors },
  } = useForm<AdminProductFormValues>({
    resolver: zodResolver(adminProductSchema),
    defaultValues: {
      name: defaultValues?.name || '',
      description: defaultValues?.description || '',
      base_price_cents: defaultValues?.base_price_cents || 0,
      currency: defaultValues?.currency || 'USD',
      category_id: defaultValues?.category_id || '',
      status: defaultValues?.status || 'draft',
      product_type: defaultValues?.product_type || 'simple',
      stock_quantity: defaultValues?.stock_quantity ?? 0,
    },
  });

  const watchedProductType = useWatch({ control, name: 'product_type' }) || 'simple';

  const handleAddTag = () => {
    const trimmed = tagInput.trim();
    if (trimmed && !tags.includes(trimmed)) {
      setTags([...tags, trimmed]);
      setTagInput('');
    }
  };

  const handleRemoveTag = (tag: string) => {
    setTags(tags.filter((t) => t !== tag));
  };

  const handleAddImage = () => {
    const trimmed = imageInput.trim();
    if (trimmed && !imageUrls.includes(trimmed)) {
      setImageUrls([...imageUrls, trimmed]);
      setImageInput('');
    }
  };

  const handleRemoveImage = (url: string) => {
    setImageUrls(imageUrls.filter((u) => u !== url));
  };

  const handleFormSubmit = (data: AdminProductFormValues) => {
    onSubmit({
      ...data,
      tags,
      image_urls: imageUrls,
      product_type: data.product_type || 'simple',
      stock_quantity: data.product_type === 'simple' ? (data.stock_quantity ?? 0) : 0,
    });
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4 max-h-[70vh] overflow-y-auto pr-2">
      {/* Product Name */}
      <div className="space-y-2">
        <Label htmlFor="name">Product Name</Label>
        <Input id="name" {...register('name')} placeholder="e.g. Premium Wireless Headphones" />
        {errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}
      </div>

      {/* Description */}
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

      {/* Price + Currency + Status */}
      <div className="grid gap-4 sm:grid-cols-3">
        <div className="space-y-2">
          <Label htmlFor="base_price_cents">Price (cents)</Label>
          <Input id="base_price_cents" type="number" {...register('base_price_cents')} placeholder="9999" />
          {errors.base_price_cents && (
            <p className="text-sm text-destructive">{errors.base_price_cents.message}</p>
          )}
        </div>
        <div className="space-y-2">
          <Label htmlFor="currency">Currency</Label>
          <select
            id="currency"
            {...register('currency')}
            className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <option value="USD">USD</option>
            <option value="EUR">EUR</option>
            <option value="GBP">GBP</option>
            <option value="AUD">AUD</option>
          </select>
        </div>
        {isEditing && (
          <div className="space-y-2">
            <Label htmlFor="status">Status</Label>
            <select
              id="status"
              {...register('status')}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              <option value="draft">Draft</option>
              <option value="active">Active</option>
              <option value="inactive">Inactive</option>
              <option value="archived">Archived</option>
            </select>
          </div>
        )}
      </div>

      {/* Category */}
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

      {/* Product Type (create only) */}
      {!isEditing && (
        <div className="space-y-2">
          <Label>Product Type</Label>
          <div className="grid grid-cols-2 gap-3">
            <button
              type="button"
              onClick={() => setValue('product_type', 'simple')}
              className={`flex items-center gap-2 rounded-lg border-2 p-3 text-left text-sm transition-colors ${
                watchedProductType === 'simple'
                  ? 'border-primary bg-primary/5'
                  : 'border-muted hover:border-muted-foreground/30'
              }`}
            >
              <Package className="h-5 w-5 text-muted-foreground" />
              <div>
                <p className="font-medium">Simple</p>
                <p className="text-xs text-muted-foreground">Direct stock & price</p>
              </div>
            </button>
            <button
              type="button"
              onClick={() => setValue('product_type', 'configurable')}
              className={`flex items-center gap-2 rounded-lg border-2 p-3 text-left text-sm transition-colors ${
                watchedProductType === 'configurable'
                  ? 'border-primary bg-primary/5'
                  : 'border-muted hover:border-muted-foreground/30'
              }`}
            >
              <Settings2 className="h-5 w-5 text-muted-foreground" />
              <div>
                <p className="font-medium">Configurable</p>
                <p className="text-xs text-muted-foreground">Options & variants</p>
              </div>
            </button>
          </div>
        </div>
      )}

      {/* Stock Quantity (simple only) */}
      {watchedProductType === 'simple' && (
        <div className="space-y-2">
          <Label htmlFor="stock_quantity">Stock Quantity</Label>
          <Input id="stock_quantity" type="number" {...register('stock_quantity')} placeholder="0" />
        </div>
      )}

      {watchedProductType === 'configurable' && !isEditing && (
        <div className="rounded-lg bg-muted/50 p-3 text-sm text-muted-foreground">
          After creating, use the <strong>Manage</strong> button to set up options (Size, Color) and generate variants with individual pricing and stock.
        </div>
      )}

      {/* Tags */}
      <div className="space-y-2">
        <Label>Tags</Label>
        <div className="flex gap-2">
          <Input
            placeholder="Add a tag..."
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.preventDefault();
                handleAddTag();
              }
            }}
          />
          <Button type="button" variant="outline" size="icon" onClick={handleAddTag}>
            <Plus className="h-4 w-4" />
          </Button>
        </div>
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-1.5 mt-2">
            {tags.map((tag) => (
              <Badge key={tag} variant="secondary" className="gap-1 pr-1">
                {tag}
                <button
                  type="button"
                  onClick={() => handleRemoveTag(tag)}
                  className="ml-1 rounded-full p-0.5 hover:bg-muted"
                >
                  <X className="h-3 w-3" />
                </button>
              </Badge>
            ))}
          </div>
        )}
      </div>

      {/* Image URLs */}
      <div className="space-y-2">
        <Label>Image URLs</Label>
        <div className="flex gap-2">
          <Input
            placeholder="https://images.unsplash.com/..."
            value={imageInput}
            onChange={(e) => setImageInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') {
                e.preventDefault();
                handleAddImage();
              }
            }}
          />
          <Button type="button" variant="outline" size="icon" onClick={handleAddImage}>
            <ImagePlus className="h-4 w-4" />
          </Button>
        </div>
        {imageUrls.length > 0 && (
          <div className="grid grid-cols-4 gap-2 mt-2">
            {imageUrls.map((url, idx) => (
              <div key={idx} className="group relative aspect-square overflow-hidden rounded-lg border bg-muted">
                <img
                  src={url}
                  alt={`Product image ${idx + 1}`}
                  className="h-full w-full object-cover"
                  onError={(e) => {
                    (e.target as HTMLImageElement).src = '';
                    (e.target as HTMLImageElement).style.display = 'none';
                  }}
                />
                <button
                  type="button"
                  onClick={() => handleRemoveImage(url)}
                  className="absolute right-1 top-1 rounded-full bg-destructive/80 p-1 text-white opacity-0 transition-opacity group-hover:opacity-100"
                >
                  <X className="h-3 w-3" />
                </button>
                {idx === 0 && (
                  <Badge className="absolute bottom-1 left-1 text-xs" variant="secondary">
                    Primary
                  </Badge>
                )}
              </div>
            ))}
          </div>
        )}
        <p className="text-xs text-muted-foreground">First image will be used as the primary product image</p>
      </div>

      <Button type="submit" disabled={isPending} className="w-full">
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}
