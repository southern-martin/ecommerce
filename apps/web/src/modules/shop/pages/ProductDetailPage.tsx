import { useState, useMemo, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Badge } from '@/shared/components/ui/badge';
import { Separator } from '@/shared/components/ui/separator';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Star, ShoppingCart, Minus, Plus, Truck, ShieldCheck, RefreshCcw } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import { PageLayout } from '@/shared/components/layout/PageLayout';
import { useProduct } from '../hooks/useProduct';
import { useCartStore } from '@/shared/stores/cart.store';
import type { ProductVariant } from '../types/shop.types';

export default function ProductDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const { data: product, isLoading } = useProduct(slug!);
  const [selectedImage, setSelectedImage] = useState(0);
  const [quantity, setQuantity] = useState(1);
  const [selectedOptions, setSelectedOptions] = useState<Record<string, string>>({});
  const addItem = useCartStore((s) => s.addItem);

  // Auto-select default variant options on load
  useEffect(() => {
    if (product?.options && product.variants) {
      const defaultVariant = product.variants.find((v) => v.is_default) || product.variants[0];
      if (defaultVariant?.option_values) {
        const defaults: Record<string, string> = {};
        defaultVariant.option_values.forEach((ov) => {
          defaults[ov.option_name] = ov.option_value_id;
        });
        setSelectedOptions(defaults);
      }
    }
  }, [product]);

  // Find the variant matching the selected options
  const activeVariant: ProductVariant | undefined = useMemo(() => {
    if (!product?.variants || !product?.options) return undefined;
    const optionCount = product.options.length;
    return product.variants.find((v) => {
      if (!v.option_values || v.option_values.length !== optionCount) return false;
      return v.option_values.every(
        (ov) => selectedOptions[ov.option_name] === ov.option_value_id
      );
    });
  }, [product, selectedOptions]);

  // Effective price: variant price or base product price
  const effectivePrice = activeVariant ? activeVariant.price_cents : product?.price ?? 0;
  const effectiveCompareAt = activeVariant?.compare_at_cents;
  const effectiveStock = activeVariant ? activeVariant.stock : product?.stock_quantity ?? 0;
  const isInStock = activeVariant ? activeVariant.stock > 0 && activeVariant.is_active : product?.in_stock ?? false;

  if (isLoading) {
    return (
      <PageLayout breadcrumbs={[{ label: 'Shop', href: '/products' }, { label: 'Loading...' }]}>
        <div className="grid gap-8 md:grid-cols-2">
          <Skeleton className="aspect-square w-full rounded-2xl" />
          <div className="space-y-4">
            <Skeleton className="h-8 w-3/4" />
            <Skeleton className="h-6 w-1/4" />
            <Skeleton className="h-24 w-full" />
            <Skeleton className="h-10 w-48" />
          </div>
        </div>
      </PageLayout>
    );
  }

  if (!product) {
    return (
      <PageLayout breadcrumbs={[{ label: 'Shop', href: '/products' }, { label: 'Not Found' }]}>
        <div className="py-16 text-center">
          <p className="text-lg text-muted-foreground">Product not found</p>
        </div>
      </PageLayout>
    );
  }

  const hasOptions = product.options && product.options.length > 0;

  return (
    <PageLayout
      breadcrumbs={[{ label: 'Shop', href: '/products' }, { label: product.name }]}
    >
      <div className="space-y-12">
        <div className="grid gap-8 md:grid-cols-2">
          {/* Image Gallery */}
          <div className="space-y-4">
            <div className="aspect-square overflow-hidden rounded-2xl bg-muted">
              {(product.images || [])[selectedImage]?.url ? (
                <img
                  src={product.images[selectedImage].url}
                  alt={product.images[selectedImage]?.alt ?? product.name}
                  className="h-full w-full object-cover"
                />
              ) : (
                <div className="flex h-full w-full items-center justify-center text-muted-foreground/40">
                  <ShoppingCart className="h-16 w-16" />
                </div>
              )}
            </div>
            {(product.images || []).length > 1 && (
              <div className="flex gap-2 overflow-x-auto">
                {product.images.map((image, index) => (
                  <button
                    key={image.id}
                    onClick={() => setSelectedImage(index)}
                    className={`h-20 w-20 flex-shrink-0 overflow-hidden rounded-lg border-2 ${
                      index === selectedImage ? 'border-primary' : 'border-transparent'
                    }`}
                  >
                    <img
                      src={image.url}
                      alt={image.alt}
                      className="h-full w-full object-cover"
                    />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Product Info */}
          <div className="space-y-6">
            <div>
              <h1 className="text-3xl font-bold">{product.name}</h1>
              <p className="mt-1 text-sm text-muted-foreground">
                Sold by {product.seller.name || 'Seller'}
              </p>
            </div>

            <div className="flex items-center gap-2">
              <div className="flex items-center">
                {Array.from({ length: 5 }).map((_, i) => (
                  <Star
                    key={i}
                    className={`h-5 w-5 ${
                      i < Math.round(product.rating)
                        ? 'fill-yellow-400 text-yellow-400'
                        : 'text-muted-foreground/30'
                    }`}
                  />
                ))}
              </div>
              <span className="text-sm text-muted-foreground">
                ({product.review_count} reviews)
              </span>
            </div>

            {/* Price */}
            <div className="flex items-baseline gap-3">
              <span className="text-3xl font-bold text-primary">{formatPrice(effectivePrice)}</span>
              {effectiveCompareAt && effectiveCompareAt > effectivePrice && (
                <>
                  <span className="text-xl text-muted-foreground line-through">
                    {formatPrice(effectiveCompareAt)}
                  </span>
                  <Badge variant="destructive">
                    {Math.round(
                      ((effectiveCompareAt - effectivePrice) / effectiveCompareAt) * 100
                    )}
                    % OFF
                  </Badge>
                </>
              )}
            </div>

            <Separator />

            <p className="leading-relaxed text-muted-foreground">{product.description}</p>

            {/* Option-based Variant Selector */}
            {hasOptions && (
              <div className="space-y-4">
                {product.options!.map((option) => (
                  <div key={option.id} className="space-y-2">
                    <span className="text-sm font-medium">{option.name}</span>
                    <div className="flex flex-wrap gap-2">
                      {option.values.map((optVal) => {
                        const isSelected = selectedOptions[option.name] === optVal.id;
                        const isColor = !!optVal.color_hex;

                        if (isColor) {
                          return (
                            <button
                              key={optVal.id}
                              title={optVal.value}
                              onClick={() =>
                                setSelectedOptions((prev) => ({ ...prev, [option.name]: optVal.id }))
                              }
                              className={`h-10 w-10 rounded-full border-2 transition-all ${
                                isSelected
                                  ? 'border-primary ring-2 ring-primary ring-offset-2'
                                  : 'border-muted-foreground/30 hover:border-muted-foreground'
                              }`}
                              style={{ backgroundColor: optVal.color_hex }}
                            />
                          );
                        }

                        return (
                          <Button
                            key={optVal.id}
                            variant={isSelected ? 'default' : 'outline'}
                            size="sm"
                            onClick={() =>
                              setSelectedOptions((prev) => ({ ...prev, [option.name]: optVal.id }))
                            }
                          >
                            {optVal.value}
                          </Button>
                        );
                      })}
                    </div>
                  </div>
                ))}

                {/* Variant info */}
                {activeVariant && (
                  <div className="text-sm text-muted-foreground">
                    SKU: {activeVariant.sku}
                    {activeVariant.stock > 0 && activeVariant.stock <= 5 && (
                      <span className="ml-3 text-orange-500 font-medium">
                        Only {activeVariant.stock} left!
                      </span>
                    )}
                  </div>
                )}
                {hasOptions && !activeVariant && Object.keys(selectedOptions).length > 0 && (
                  <p className="text-sm text-destructive">
                    This combination is not available.
                  </p>
                )}
              </div>
            )}

            {/* Product Attributes */}
            {product.attributes && product.attributes.length > 0 && (
              <>
                <Separator />
                <div className="space-y-2">
                  <span className="text-sm font-medium">Specifications</span>
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    {product.attributes.map((attr) => (
                      <div key={attr.id} className="flex gap-2">
                        <span className="text-muted-foreground">{attr.attribute_name}:</span>
                        <span className="font-medium">
                          {attr.values && attr.values.length > 0
                            ? attr.values.join(', ')
                            : attr.value}
                        </span>
                      </div>
                    ))}
                  </div>
                </div>
              </>
            )}

            {/* Quantity + Add to Cart */}
            <div className="flex items-center gap-4">
              <div className="flex items-center rounded-xl border">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                >
                  <Minus className="h-4 w-4" />
                </Button>
                <span className="w-12 text-center text-sm font-medium">{quantity}</span>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setQuantity((q) => Math.min(q + 1, effectiveStock || 99))}
                >
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
              <Button
                size="lg"
                disabled={!isInStock || (hasOptions && !activeVariant)}
                className="flex-1 rounded-xl"
                onClick={() => {
                  addItem({
                    id: `${product.id}-${activeVariant?.id || 'default'}`,
                    product_id: product.id,
                    product_name: product.name,
                    product_slug: product.slug,
                    price_cents: effectivePrice,
                    quantity,
                    image_url: (product.images || [])[0]?.url,
                    variant_id: activeVariant?.id || undefined,
                    seller_id: product.seller?.id,
                  });
                }}
              >
                <ShoppingCart className="mr-2 h-5 w-5" />
                {!isInStock
                  ? 'Out of Stock'
                  : hasOptions && !activeVariant
                    ? 'Select Options'
                    : 'Add to Cart'}
              </Button>
            </div>

            {/* Shipping & Return Info Badges */}
            <div className="flex flex-wrap gap-3">
              <div className="flex items-center gap-2 rounded-lg bg-muted/50 px-3 py-2 text-sm text-muted-foreground">
                <Truck className="h-4 w-4 text-primary" />
                Free shipping over $50
              </div>
              <div className="flex items-center gap-2 rounded-lg bg-muted/50 px-3 py-2 text-sm text-muted-foreground">
                <ShieldCheck className="h-4 w-4 text-primary" />
                Secure payment
              </div>
              <div className="flex items-center gap-2 rounded-lg bg-muted/50 px-3 py-2 text-sm text-muted-foreground">
                <RefreshCcw className="h-4 w-4 text-primary" />
                30-day returns
              </div>
            </div>

            {isInStock ? (
              <Badge variant="secondary" className="text-green-600">
                In Stock{effectiveStock > 0 ? ` (${effectiveStock})` : ''}
              </Badge>
            ) : (
              <Badge variant="destructive">Out of Stock</Badge>
            )}
          </div>
        </div>
      </div>
    </PageLayout>
  );
}
