import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Minus, Plus, Trash2 } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import type { CartItem as CartItemType } from '../services/cart.api';

interface CartItemProps {
  item: CartItemType;
  onUpdateQuantity: (itemId: string, quantity: number) => void;
  onRemove: (itemId: string) => void;
}

export function CartItem({ item, onUpdateQuantity, onRemove }: CartItemProps) {
  return (
    <div className="flex gap-4 py-4">
      <Link to={`/products/${item.slug}`} className="h-24 w-24 flex-shrink-0 overflow-hidden rounded-md bg-muted">
        <img src={item.image_url} alt={item.name} className="h-full w-full object-cover" />
      </Link>

      <div className="flex flex-1 flex-col justify-between">
        <div className="flex justify-between">
          <div>
            <Link to={`/products/${item.slug}`} className="text-sm font-medium hover:text-primary">
              {item.name}
            </Link>
            {item.variant_name && (
              <p className="mt-0.5 text-xs text-muted-foreground">{item.variant_name}</p>
            )}
          </div>
          <span className="text-sm font-medium">{formatPrice(item.price * item.quantity)}</span>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center rounded-md border">
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8"
              onClick={() => onUpdateQuantity(item.id, Math.max(1, item.quantity - 1))}
            >
              <Minus className="h-3 w-3" />
            </Button>
            <span className="w-8 text-center text-sm">{item.quantity}</span>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8"
              onClick={() => onUpdateQuantity(item.id, item.quantity + 1)}
            >
              <Plus className="h-3 w-3" />
            </Button>
          </div>
          <Button variant="ghost" size="icon" className="h-8 w-8 text-destructive" onClick={() => onRemove(item.id)}>
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
