import { Link } from 'react-router-dom';
import { Card, CardContent } from '@/shared/components/ui/card';
import { Button } from '@/shared/components/ui/button';
import { Star, ShoppingCart } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import type { Product } from '../types/shop.types';

interface ProductCardProps {
  product: Product;
  onAddToCart?: (productId: string) => void;
}

export function ProductCard({ product, onAddToCart }: ProductCardProps) {
  const images = product.images || [];
  const primaryImage = images.find((img) => img.is_primary) ?? images[0];

  return (
    <Card className="group overflow-hidden transition-shadow hover:shadow-lg">
      <Link to={`/products/${product.slug}`}>
        <div className="aspect-square overflow-hidden bg-muted">
          {primaryImage?.url ? (
            <img
              src={primaryImage.url}
              alt={primaryImage.alt ?? product.name}
              className="h-full w-full object-cover transition-transform group-hover:scale-105"
            />
          ) : (
            <div className="flex h-full w-full items-center justify-center text-muted-foreground/40">
              <ShoppingCart className="h-12 w-12" />
            </div>
          )}
        </div>
      </Link>
      <CardContent className="p-4">
        <Link to={`/products/${product.slug}`}>
          <h3 className="line-clamp-2 text-sm font-medium hover:text-primary">
            {product.name}
          </h3>
        </Link>

        <div className="mt-1 flex items-center gap-1">
          {Array.from({ length: 5 }).map((_, i) => (
            <Star
              key={i}
              className={`h-3.5 w-3.5 ${
                i < Math.round(product.rating)
                  ? 'fill-yellow-400 text-yellow-400'
                  : 'text-muted-foreground/30'
              }`}
            />
          ))}
          <span className="ml-1 text-xs text-muted-foreground">
            ({product.review_count})
          </span>
        </div>

        <div className="mt-2 flex items-center gap-2">
          <span className="text-lg font-bold">{formatPrice(product.price)}</span>
          {product.compare_at_price && product.compare_at_price > product.price && (
            <span className="text-sm text-muted-foreground line-through">
              {formatPrice(product.compare_at_price)}
            </span>
          )}
        </div>

        <Button
          size="sm"
          className="mt-3 w-full"
          onClick={() => onAddToCart?.(product.id)}
          disabled={!product.in_stock}
        >
          <ShoppingCart className="mr-2 h-4 w-4" />
          {product.in_stock ? 'Add to Cart' : 'Out of Stock'}
        </Button>
      </CardContent>
    </Card>
  );
}
