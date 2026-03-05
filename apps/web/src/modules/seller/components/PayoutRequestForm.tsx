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
import { formatPrice } from '@/shared/lib/utils';
import type { RequestPayoutInput } from '../services/seller-wallet.api';

interface PayoutRequestFormProps {
  availableBalance: number;
  onSubmit: (data: RequestPayoutInput) => void;
  isPending: boolean;
}

export function PayoutRequestForm({
  availableBalance,
  onSubmit,
  isPending,
}: PayoutRequestFormProps) {
  const maxAmount = availableBalance / 100;

  const payoutSchema = z.object({
    amount: z.coerce
      .number()
      .min(0.01, 'Amount must be at least $0.01')
      .max(maxAmount, `Amount cannot exceed ${formatPrice(availableBalance)}`),
    method: z.enum(['stripe_connect', 'bank_transfer']),
  });

  type PayoutFormValues = z.infer<typeof payoutSchema>;

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    reset,
    formState: { errors },
  } = useForm<PayoutFormValues>({
    resolver: zodResolver(payoutSchema),
    defaultValues: {
      amount: 0,
      method: 'stripe_connect',
    },
  });

  const methodValue = watch('method');

  const handleFormSubmit = (data: PayoutFormValues) => {
    onSubmit({
      amount_cents: Math.round(data.amount * 100),
      currency: 'USD',
      method: data.method,
    });
    reset();
  };

  return (
    <div className="space-y-6">
      <div className="rounded-lg border bg-muted/50 p-4">
        <p className="text-sm text-muted-foreground">Available for payout</p>
        <p className="text-3xl font-bold">{formatPrice(availableBalance)}</p>
      </div>

      <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="amount">Payout Amount ($)</Label>
          <Input
            id="amount"
            type="number"
            step="0.01"
            min="0.01"
            max={maxAmount}
            {...register('amount')}
            placeholder="0.00"
          />
          {errors.amount && (
            <p className="text-sm text-destructive">{errors.amount.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="method">Payout Method</Label>
          <Select
            value={methodValue}
            onValueChange={(val) => setValue('method', val as PayoutFormValues['method'])}
          >
            <SelectTrigger>
              <SelectValue placeholder="Select method" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="stripe_connect">Stripe Connect</SelectItem>
              <SelectItem value="bank_transfer">Bank Transfer</SelectItem>
            </SelectContent>
          </Select>
          {errors.method && (
            <p className="text-sm text-destructive">{errors.method.message}</p>
          )}
        </div>

        <Button type="submit" disabled={isPending || availableBalance <= 0}>
          {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Request Payout
        </Button>
      </form>
    </div>
  );
}
