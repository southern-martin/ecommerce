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

const returnSchema = z.object({
  order_id: z.string().min(1, 'Order is required'),
  reason: z.string().min(1, 'Reason is required'),
  description: z.string().min(10, 'Please provide more details'),
});

type ReturnFormValues = z.infer<typeof returnSchema>;

interface ReturnRequestFormProps {
  orderId?: string;
  onSubmit: (data: ReturnFormValues) => void;
  isPending?: boolean;
}

const returnReasons = [
  'Defective product',
  'Wrong item received',
  'Item not as described',
  'Changed my mind',
  'Better price elsewhere',
  'Other',
];

export function ReturnRequestForm({ orderId, onSubmit, isPending }: ReturnRequestFormProps) {
  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<ReturnFormValues>({
    resolver: zodResolver(returnSchema),
    defaultValues: { order_id: orderId ?? '' },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="order_id">Order ID</Label>
        <Input id="order_id" {...register('order_id')} disabled={!!orderId} />
        {errors.order_id && <p className="text-sm text-destructive">{errors.order_id.message}</p>}
      </div>

      <div className="space-y-2">
        <Label>Reason</Label>
        <Select onValueChange={(value) => setValue('reason', value)}>
          <SelectTrigger>
            <SelectValue placeholder="Select a reason" />
          </SelectTrigger>
          <SelectContent>
            {returnReasons.map((reason) => (
              <SelectItem key={reason} value={reason}>
                {reason}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        {errors.reason && <p className="text-sm text-destructive">{errors.reason.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <textarea
          id="description"
          {...register('description')}
          rows={4}
          placeholder="Describe the issue in detail..."
          className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
        {errors.description && <p className="text-sm text-destructive">{errors.description.message}</p>}
      </div>

      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Submit Return Request
      </Button>
    </form>
  );
}
