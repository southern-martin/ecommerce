import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';
import { COUPON_TYPES } from '@/shared/lib/constants';
import { Loader2 } from 'lucide-react';
import type { Coupon } from '../services/seller-coupon.api';

const couponSchema = z.object({
  code: z.string().min(1, 'Coupon code is required'),
  type: z.enum(['percentage', 'fixed_amount', 'free_shipping']),
  value: z.coerce.number().min(0, 'Value must be 0 or more'),
  min_order_amount: z.coerce.number().min(0, 'Minimum order amount must be 0 or more'),
  max_uses: z.coerce.number().min(1, 'Max uses must be at least 1'),
  starts_at: z.string().min(1, 'Start date is required'),
  expires_at: z.string().min(1, 'Expiry date is required'),
});

type CouponFormValues = z.infer<typeof couponSchema>;

interface CouponFormProps {
  onSubmit: (data: CouponFormValues) => void;
  isPending: boolean;
  defaultValues?: Partial<Coupon>;
  submitLabel?: string;
}

export function CouponForm({
  onSubmit,
  isPending,
  defaultValues,
  submitLabel = 'Create Coupon',
}: CouponFormProps) {
  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CouponFormValues>({
    resolver: zodResolver(couponSchema),
    defaultValues: {
      code: defaultValues?.code ?? '',
      type: defaultValues?.type ?? 'percentage',
      value: defaultValues?.value ?? 0,
      min_order_amount: defaultValues?.min_order_amount ?? 0,
      max_uses: defaultValues?.max_uses ?? 100,
      starts_at: defaultValues?.starts_at?.slice(0, 10) ?? '',
      expires_at: defaultValues?.expires_at?.slice(0, 10) ?? '',
    },
  });

  const typeValue = watch('type');

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="code">Coupon Code</Label>
        <Input id="code" {...register('code')} placeholder="e.g. SAVE20" />
        {errors.code && <p className="text-sm text-destructive">{errors.code.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="type">Type</Label>
        <Select value={typeValue} onValueChange={(val) => setValue('type', val as CouponFormValues['type'])}>
          <SelectTrigger>
            <SelectValue placeholder="Select type" />
          </SelectTrigger>
          <SelectContent>
            {COUPON_TYPES.map((type) => (
              <SelectItem key={type} value={type}>
                {type.replace(/_/g, ' ')}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {errors.type && <p className="text-sm text-destructive">{errors.type.message}</p>}
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="value">Value</Label>
          <Input id="value" type="number" {...register('value')} />
          {errors.value && <p className="text-sm text-destructive">{errors.value.message}</p>}
        </div>
        <div className="space-y-2">
          <Label htmlFor="min_order_amount">Min Order Amount (cents)</Label>
          <Input id="min_order_amount" type="number" {...register('min_order_amount')} />
          {errors.min_order_amount && (
            <p className="text-sm text-destructive">{errors.min_order_amount.message}</p>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="max_uses">Max Uses</Label>
        <Input id="max_uses" type="number" {...register('max_uses')} />
        {errors.max_uses && <p className="text-sm text-destructive">{errors.max_uses.message}</p>}
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="starts_at">Start Date</Label>
          <Input id="starts_at" type="date" {...register('starts_at')} />
          {errors.starts_at && <p className="text-sm text-destructive">{errors.starts_at.message}</p>}
        </div>
        <div className="space-y-2">
          <Label htmlFor="expires_at">Expiry Date</Label>
          <Input id="expires_at" type="date" {...register('expires_at')} />
          {errors.expires_at && (
            <p className="text-sm text-destructive">{errors.expires_at.message}</p>
          )}
        </div>
      </div>

      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}
