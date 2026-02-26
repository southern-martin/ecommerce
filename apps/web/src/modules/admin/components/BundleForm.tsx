import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Loader2 } from 'lucide-react';

const bundleSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  description: z.string().optional(),
  discount_percentage: z.coerce.number().min(1).max(100, 'Must be between 1 and 100'),
});

type BundleFormValues = z.infer<typeof bundleSchema>;

interface BundleFormProps {
  defaultValues?: Partial<BundleFormValues>;
  onSubmit: (data: BundleFormValues) => void;
  isPending?: boolean;
  submitLabel?: string;
}

export function BundleForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Create Bundle',
}: BundleFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<BundleFormValues>({
    resolver: zodResolver(bundleSchema),
    defaultValues: {
      name: '',
      description: '',
      discount_percentage: 10,
      ...defaultValues,
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="bundle-name">Name</Label>
        <Input id="bundle-name" {...register('name')} placeholder="Starter Bundle" />
        {errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="bundle-description">Description</Label>
        <textarea
          id="bundle-description"
          {...register('description')}
          rows={3}
          className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="bundle-discount">Discount Percentage</Label>
        <Input id="bundle-discount" type="number" {...register('discount_percentage')} />
        {errors.discount_percentage && (
          <p className="text-sm text-destructive">{errors.discount_percentage.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}
