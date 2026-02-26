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
import { Loader2 } from 'lucide-react';

const COUPON_TYPES = ['percentage', 'fixed_amount', 'free_shipping'] as const;

const couponSchema = z.object({
  code: z.string().min(1, 'Code is required'),
  type: z.enum(COUPON_TYPES),
  value: z.coerce.number().min(0, 'Value must be 0 or more'),
  min_order_amount: z.coerce.number().optional(),
  max_uses: z.coerce.number().optional(),
  starts_at: z.string().min(1, 'Start date is required'),
  expires_at: z.string().min(1, 'Expiry date is required'),
});

type CouponFormValues = z.infer<typeof couponSchema>;

interface AdminCouponFormProps {
  defaultValues?: Partial<CouponFormValues>;
  onSubmit: (data: CouponFormValues) => void;
  isPending?: boolean;
  submitLabel?: string;
}

export function AdminCouponForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Create Coupon',
}: AdminCouponFormProps) {
  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<CouponFormValues>({
    resolver: zodResolver(couponSchema),
    defaultValues: {
      code: '',
      type: 'percentage',
      value: 0,
      starts_at: '',
      expires_at: '',
      ...defaultValues,
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="code">Code</Label>
        <Input id="code" {...register('code')} placeholder="SUMMER20" />
        {errors.code && <p className="text-sm text-destructive">{errors.code.message}</p>}
      </div>

      <div className="space-y-2">
        <Label>Type</Label>
        <Select
          defaultValue={defaultValues?.type || 'percentage'}
          onValueChange={(value) => setValue('type', value as any)}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="percentage">Percentage</SelectItem>
            <SelectItem value="fixed_amount">Fixed Amount</SelectItem>
            <SelectItem value="free_shipping">Free Shipping</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="value">Value</Label>
          <Input id="value" type="number" {...register('value')} />
          {errors.value && <p className="text-sm text-destructive">{errors.value.message}</p>}
        </div>
        <div className="space-y-2">
          <Label htmlFor="min_order_amount">Min Order Amount</Label>
          <Input id="min_order_amount" type="number" {...register('min_order_amount')} />
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="max_uses">Max Uses</Label>
        <Input id="max_uses" type="number" {...register('max_uses')} />
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="starts_at">Starts At</Label>
          <Input id="starts_at" type="datetime-local" {...register('starts_at')} />
          {errors.starts_at && (
            <p className="text-sm text-destructive">{errors.starts_at.message}</p>
          )}
        </div>
        <div className="space-y-2">
          <Label htmlFor="expires_at">Expires At</Label>
          <Input id="expires_at" type="datetime-local" {...register('expires_at')} />
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
