import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Star, ShoppingCart, Heart, Eye } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import { useCartStore } from '@/shared/stores/cart.store';
import type { Product } from '../types/shop.types';

interface ProductCardProps {
  product: Product;
}

export function ProductCard({ product }: ProductCardProps) {
  const addItem = useCartStore((s) => s.addItem);
  const images = product.images || [];
  const primaryImage = images.find((img) => img.is_primary) ?? images[0];
  const discount = product.compare_at_price && product.compare_at_price > product.price
    ? Math.round((1 - product.price / product.compare_at_price) * 100)
    : 0;

  const handleAddToCart = () => {
    addItem({
      id: `${product.id}-default`,
      product_id: product.id,
      product_name: product.name,
      product_slug: product.slug,
      price_cents: product.price,
      quantity: 1,
      image_url: primaryImage?.url,
      seller_id: product.seller?.id,
    });
  };

  return (
    <div className="group relative overflow-hidden rounded-2xl border bg-card transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
      {/* Image Container */}
      <Link to={`/products/${product.slug}`} className="relative block overflow-hidden">
        <div className="aspect-square overflow-hidden bg-muted">
          {primaryImage?.url ? (
            <img
              src={primaryImage.url}
              alt={primaryImage.alt ?? product.name}
              className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
            />
          ) : (
            <div className="flex h-full w-full items-center justify-center bg-gradient-to-br from-muted to-muted/50">
              <ShoppingCart className="h-16 w-16 text-muted-foreground/20" />
            </div>
          )}
        </div>

        {/* Badges */}
        <div className="absolute left-3 top-3 flex flex-col gap-1.5">
          {discount > 0 && (
            <span className="rounded-lg bg-red-500 px-2.5 py-1 text-xs font-bold text-white shadow-sm">
              -{discount}%
            </span>
          )}
          {!product.in_stock && (
            <span className="rounded-lg bg-gray-900/80 px-2.5 py-1 text-xs font-medium text-white backdrop-blur-sm">
              Sold Out
            </span>
          )}
        </div>

        {/* Quick Actions (visible on hover) */}
        <div className="absolute right-3 top-3 flex flex-col gap-2 opacity-0 transition-all duration-300 group-hover:opacity-100 translate-x-2 group-hover:translate-x-0">
          <button className="flex h-9 w-9 items-center justify-center rounded-full bg-white/90 text-gray-700 shadow-md backdrop-blur-sm transition-colors hover:bg-primary hover:text-white">
            <Heart className="h-4 w-4" />
          </button>
          <button className="flex h-9 w-9 items-center justify-center rounded-full bg-white/90 text-gray-700 shadow-md backdrop-blur-sm transition-colors hover:bg-primary hover:text-white">
            <Eye className="h-4 w-4" />
          </button>
        </div>
      </Link>

      {/* Content */}
      <div className="p-4">
        <Link to={`/products/${product.slug}`}>
          <h3 className="line-clamp-2 text-sm font-semibold leading-snug transition-colors hover:text-primary">
            {product.name}
          </h3>
        </Link>

        {/* Rating */}
        <div className="mt-2 flex items-center gap-1.5">
          <div className="flex items-center">
            {Array.from({ length: 5 }).map((_, i) => (
              <Star
                key={i}
                className={`h-3.5 w-3.5 ${
                  i < Math.round(product.rating)
                    ? 'fill-amber-400 text-amber-400'
                    : 'fill-muted text-muted'
                }`}
              />
            ))}
          </div>
          <span className="text-xs text-muted-foreground">
            ({product.review_count})
          </span>
        </div>

        {/* Price */}
        <div className="mt-3 flex items-baseline gap-2">
          <span className="text-xl font-bold text-primary">{formatPrice(product.price)}</span>
          {product.compare_at_price && product.compare_at_price > product.price && (
            <span className="text-sm text-muted-foreground line-through">
              {formatPrice(product.compare_at_price)}
            </span>
          )}
        </div>

        {/* Add to Cart */}
        <Button
          size="sm"
          className="mt-3 w-full rounded-xl font-medium"
          onClick={handleAddToCart}
          disabled={!product.in_stock}
        >
          <ShoppingCart className="mr-2 h-4 w-4" />
          {product.in_stock ? 'Add to Cart' : 'Out of Stock'}
        </Button>
      </div>
    </div>
  );
}
