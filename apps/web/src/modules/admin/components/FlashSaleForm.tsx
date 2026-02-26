import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Loader2 } from 'lucide-react';

const flashSaleSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  discount_percentage: z.coerce.number().min(1).max(100, 'Must be between 1 and 100'),
  starts_at: z.string().min(1, 'Start date is required'),
  ends_at: z.string().min(1, 'End date is required'),
});

type FlashSaleFormValues = z.infer<typeof flashSaleSchema>;

interface FlashSaleFormProps {
  defaultValues?: Partial<FlashSaleFormValues>;
  onSubmit: (data: FlashSaleFormValues) => void;
  isPending?: boolean;
  submitLabel?: string;
}

export function FlashSaleForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Create Flash Sale',
}: FlashSaleFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<FlashSaleFormValues>({
    resolver: zodResolver(flashSaleSchema),
    defaultValues: {
      name: '',
      discount_percentage: 10,
      starts_at: '',
      ends_at: '',
      ...defaultValues,
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="flash-name">Name</Label>
        <Input id="flash-name" {...register('name')} placeholder="Weekend Flash Sale" />
        {errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="discount_percentage">Discount Percentage</Label>
        <Input id="discount_percentage" type="number" {...register('discount_percentage')} />
        {errors.discount_percentage && (
          <p className="text-sm text-destructive">{errors.discount_percentage.message}</p>
        )}
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="flash-starts">Starts At</Label>
          <Input id="flash-starts" type="datetime-local" {...register('starts_at')} />
          {errors.starts_at && (
            <p className="text-sm text-destructive">{errors.starts_at.message}</p>
          )}
        </div>
        <div className="space-y-2">
          <Label htmlFor="flash-ends">Ends At</Label>
          <Input id="flash-ends" type="datetime-local" {...register('ends_at')} />
          {errors.ends_at && (
            <p className="text-sm text-destructive">{errors.ends_at.message}</p>
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
