import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';

interface PaymentFormProps {
  onBack: () => void;
  onContinue: (paymentMethod: string) => void;
  couponCode: string;
  onCouponChange: (code: string) => void;
}

export function PaymentForm({ onBack, onContinue, couponCode, onCouponChange }: PaymentFormProps) {
  const [paymentMethod, setPaymentMethod] = useState<'card' | 'cod'>('card');

  return (
    <div className="space-y-6">
      <div className="space-y-3">
        <Label className="text-base font-semibold">Payment Method</Label>
        <div className="grid grid-cols-2 gap-3">
          <button
            type="button"
            onClick={() => setPaymentMethod('card')}
            className={`rounded-lg border-2 p-4 text-left transition-colors ${
              paymentMethod === 'card' ? 'border-primary bg-primary/5' : 'border-muted hover:border-muted-foreground/30'
            }`}
          >
            <div className="font-medium">Credit Card</div>
            <div className="text-sm text-muted-foreground">Pay with Visa, Mastercard, etc.</div>
          </button>
          <button
            type="button"
            onClick={() => setPaymentMethod('cod')}
            className={`rounded-lg border-2 p-4 text-left transition-colors ${
              paymentMethod === 'cod' ? 'border-primary bg-primary/5' : 'border-muted hover:border-muted-foreground/30'
            }`}
          >
            <div className="font-medium">Cash on Delivery</div>
            <div className="text-sm text-muted-foreground">Pay when you receive</div>
          </button>
        </div>
      </div>

      {paymentMethod === 'card' && (
        <div className="space-y-4 rounded-lg border p-4">
          <div className="space-y-2">
            <Label htmlFor="cardNumber">Card Number</Label>
            <Input id="cardNumber" placeholder="4242 4242 4242 4242" />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="expiry">Expiry Date</Label>
              <Input id="expiry" placeholder="MM/YY" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cvv">CVV</Label>
              <Input id="cvv" placeholder="123" />
            </div>
          </div>
          <p className="text-xs text-muted-foreground">
            Demo mode — no real payment processing.
          </p>
        </div>
      )}

      {paymentMethod === 'cod' && (
        <div className="rounded-lg border bg-muted/30 p-4">
          <p className="text-sm text-muted-foreground">
            You will pay the full amount when your order is delivered to your doorstep.
          </p>
        </div>
      )}

      {/* Coupon */}
      <div className="space-y-2">
        <Label>Coupon Code (optional)</Label>
        <div className="flex gap-2">
          <Input
            placeholder="Enter coupon code"
            value={couponCode}
            onChange={(e) => onCouponChange(e.target.value)}
          />
          <Button variant="outline" type="button">Apply</Button>
        </div>
      </div>

      <div className="flex gap-4">
        <Button variant="outline" onClick={onBack}>Back</Button>
        <Button className="flex-1" onClick={() => onContinue(paymentMethod)}>
          Continue to Review
        </Button>
      </div>
    </div>
  );
}
