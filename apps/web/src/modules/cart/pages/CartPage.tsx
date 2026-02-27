import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ShoppingCart, ShoppingBag } from 'lucide-react';
import { CartItem } from '../components/CartItem';
import { CartSummary } from '../components/CartSummary';
import { useCart } from '../hooks/useCart';
import { PageLayout } from '@/shared/components/layout/PageLayout';

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
      <PageLayout
        title="Shopping Cart"
        icon={ShoppingCart}
        breadcrumbs={[{ label: 'Cart' }]}
      >
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-32 w-full rounded-2xl" />
          ))}
        </div>
      </PageLayout>
    );
  }

  if (!cart || cart.items.length === 0) {
    return (
      <PageLayout
        title="Shopping Cart"
        icon={ShoppingCart}
        breadcrumbs={[{ label: 'Cart' }]}
      >
        <div className="flex flex-col items-center justify-center py-24">
          <div className="flex h-24 w-24 items-center justify-center rounded-full bg-gradient-to-br from-muted/50 to-muted">
            <ShoppingCart className="h-12 w-12 text-muted-foreground/40" />
          </div>
          <h2 className="mt-6 text-2xl font-bold">Your cart is empty</h2>
          <p className="mt-2 text-muted-foreground">
            Looks like you haven&apos;t added anything to your cart yet.
          </p>
          <Button asChild className="mt-6 rounded-xl" size="lg">
            <Link to="/products">
              <ShoppingBag className="mr-2 h-4 w-4" />
              Start Shopping
            </Link>
          </Button>
        </div>
      </PageLayout>
    );
  }

  return (
    <PageLayout
      title="Shopping Cart"
      icon={ShoppingCart}
      breadcrumbs={[{ label: 'Cart' }]}
    >
      <div className="grid gap-8 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <div className="divide-y rounded-2xl border bg-card">
            <div className="px-6">
              {cart.items.map((item) => (
                <CartItem
                  key={item.id}
                  item={item}
                  onUpdateQuantity={handleUpdateQuantity}
                  onRemove={handleRemove}
                />
              ))}
            </div>
          </div>
        </div>

        <div>
          <CartSummary subtotal={cart.subtotal} itemCount={cart.item_count} />
        </div>
      </div>
    </PageLayout>
  );
}
