import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Tag, Loader2, Check, X } from 'lucide-react';
import { useValidateCoupon } from '../hooks/useCoupons';
import { formatPrice } from '@/shared/lib/utils';

interface CouponInputProps {
  orderTotal: number;
  onApply: (code: string, discount: number) => void;
}

export function CouponInput({ orderTotal, onApply }: CouponInputProps) {
  const [code, setCode] = useState('');
  const validateCoupon = useValidateCoupon();

  const handleApply = () => {
    if (!code.trim()) return;
    validateCoupon.mutate(
      { code: code.trim(), orderTotal },
      {
        onSuccess: (result) => {
          if (result.valid) {
            onApply(code.trim(), result.discount_amount);
          }
        },
      }
    );
  };

  return (
    <div className="space-y-2">
      <div className="flex gap-2">
        <div className="relative flex-1">
          <Tag className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={code}
            onChange={(e) => setCode(e.target.value.toUpperCase())}
            placeholder="Enter coupon code"
            className="pl-10"
          />
        </div>
        <Button onClick={handleApply} disabled={!code.trim() || validateCoupon.isPending}>
          {validateCoupon.isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            'Apply'
          )}
        </Button>
      </div>
      {validateCoupon.data && (
        <p
          className={`flex items-center gap-1 text-sm ${
            validateCoupon.data.valid ? 'text-green-600' : 'text-destructive'
          }`}
        >
          {validateCoupon.data.valid ? (
            <>
              <Check className="h-3.5 w-3.5" />
              Discount: {formatPrice(validateCoupon.data.discount_amount)}
            </>
          ) : (
            <>
              <X className="h-3.5 w-3.5" />
              {validateCoupon.data.message ?? 'Invalid coupon code'}
            </>
          )}
        </p>
      )}
    </div>
  );
}
