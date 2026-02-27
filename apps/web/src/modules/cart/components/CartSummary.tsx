import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Separator } from '@/shared/components/ui/separator';
import { formatPrice } from '@/shared/lib/utils';
import { ShieldCheck } from 'lucide-react';

interface CartSummaryProps {
  subtotal: number;
  itemCount: number;
  estimatedShipping?: number;
  showCheckoutButton?: boolean;
}

export function CartSummary({
  subtotal,
  itemCount,
  estimatedShipping = 0,
  showCheckoutButton = true,
}: CartSummaryProps) {
  const total = subtotal + estimatedShipping;

  return (
    <div className="rounded-2xl border bg-card p-6">
      <h3 className="text-lg font-semibold tracking-tight">Order Summary</h3>

      <div className="mt-5 space-y-4">
        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">
            Subtotal ({itemCount} item{itemCount !== 1 ? 's' : ''})
          </span>
          <span className="font-medium">{formatPrice(subtotal)}</span>
        </div>

        <div className="flex justify-between text-sm">
          <span className="text-muted-foreground">Estimated Shipping</span>
          <span className="font-medium">
            {estimatedShipping > 0 ? formatPrice(estimatedShipping) : 'Free'}
          </span>
        </div>

        <Separator />

        <div className="flex justify-between text-lg font-bold text-primary">
          <span>Total</span>
          <span>{formatPrice(total)}</span>
        </div>
      </div>

      {showCheckoutButton && (
        <div className="mt-6 space-y-3">
          <Button asChild className="w-full rounded-xl font-semibold" size="lg">
            <Link to="/checkout">Proceed to Checkout</Link>
          </Button>
          <div className="flex items-center justify-center gap-1.5 text-xs text-muted-foreground">
            <ShieldCheck className="h-3.5 w-3.5" />
            <span>Secure checkout</span>
          </div>
        </div>
      )}
    </div>
  );
}
