import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Separator } from '@/shared/components/ui/separator';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ArrowLeft, ShoppingCart } from 'lucide-react';
import { CartItem } from '../components/CartItem';
import { CartSummary } from '../components/CartSummary';
import { useCart } from '../hooks/useCart';

export default function CartPage() {
  const { cart, isLoading, updateQuantity, removeFromCart } = useCart();

  const handleUpdateQuantity = (itemId: string, quantity: number) => {
    updateQuantity.mutate({ itemId, quantity });
  };

  const handleRemove = (itemId: string) => {
    removeFromCart.mutate(itemId);
  };

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={i} className="h-32 w-full" />
        ))}
      </div>
    );
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-24">
        <ShoppingCart className="h-16 w-16 text-muted-foreground/50" />
        <h2 className="mt-6 text-2xl font-bold">Your cart is empty</h2>
        <p className="mt-2 text-muted-foreground">
          Looks like you haven&apos;t added anything to your cart yet.
        </p>
        <Button asChild className="mt-6">
          <Link to="/products">Start Shopping</Link>
        </Button>
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6 flex items-center gap-4">
        <Button asChild variant="ghost" size="sm">
          <Link to="/products">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Continue Shopping
          </Link>
        </Button>
      </div>

      <h1 className="mb-8 text-3xl font-bold">Shopping Cart</h1>

      <div className="grid gap-8 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <div className="divide-y rounded-lg border">
            <div className="px-6">
              {cart.items.map((item) => (
                <div key={item.id}>
                  <CartItem
                    item={item}
                    onUpdateQuantity={handleUpdateQuantity}
                    onRemove={handleRemove}
                  />
                  <Separator />
                </div>
              ))}
            </div>
          </div>
        </div>

        <div>
          <CartSummary subtotal={cart.subtotal} itemCount={cart.item_count} />
        </div>
      </div>
    </div>
  );
}
