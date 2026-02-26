import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Separator } from '@/shared/components/ui/separator';
import { ShoppingCart, X } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import { CartItem } from './CartItem';
import { useCart } from '../hooks/useCart';

interface CartDrawerProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CartDrawer({ isOpen, onClose }: CartDrawerProps) {
  const { cart, updateQuantity, removeFromCart } = useCart();

  const handleUpdateQuantity = (itemId: string, quantity: number) => {
    updateQuantity.mutate({ itemId, quantity });
  };

  const handleRemove = (itemId: string) => {
    removeFromCart.mutate(itemId);
  };

  return (
    <>
      {isOpen && (
        <div className="fixed inset-0 z-40 bg-black/50" onClick={onClose} />
      )}
      <div
        className={`fixed right-0 top-0 z-50 h-full w-full max-w-md transform bg-background shadow-xl transition-transform duration-300 ${
          isOpen ? 'translate-x-0' : 'translate-x-full'
        }`}
      >
        <div className="flex h-full flex-col">
          <div className="flex items-center justify-between border-b px-6 py-4">
            <div className="flex items-center gap-2">
              <ShoppingCart className="h-5 w-5" />
              <h2 className="text-lg font-semibold">Cart</h2>
              {cart && (
                <span className="text-sm text-muted-foreground">
                  ({cart.item_count} item{cart.item_count !== 1 ? 's' : ''})
                </span>
              )}
            </div>
            <Button variant="ghost" size="icon" onClick={onClose}>
              <X className="h-5 w-5" />
            </Button>
          </div>

          <div className="flex-1 overflow-y-auto px-6">
            {!cart || cart.items.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-16">
                <ShoppingCart className="h-12 w-12 text-muted-foreground/50" />
                <p className="mt-4 text-muted-foreground">Your cart is empty</p>
                <Button asChild variant="outline" className="mt-4" onClick={onClose}>
                  <Link to="/products">Continue Shopping</Link>
                </Button>
              </div>
            ) : (
              <div className="divide-y">
                {cart.items.map((item) => (
                  <CartItem
                    key={item.id}
                    item={item}
                    onUpdateQuantity={handleUpdateQuantity}
                    onRemove={handleRemove}
                  />
                ))}
              </div>
            )}
          </div>

          {cart && cart.items.length > 0 && (
            <div className="border-t px-6 py-4">
              <div className="flex justify-between text-lg font-semibold">
                <span>Subtotal</span>
                <span>{formatPrice(cart.subtotal)}</span>
              </div>
              <Separator className="my-3" />
              <Button asChild className="w-full" size="lg" onClick={onClose}>
                <Link to="/checkout">Checkout</Link>
              </Button>
              <Button asChild variant="outline" className="mt-2 w-full" onClick={onClose}>
                <Link to="/cart">View Cart</Link>
              </Button>
            </div>
          )}
        </div>
      </div>
    </>
  );
}
