import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Loader2 } from 'lucide-react';

const attributeGroupSchema = z.object({
  name: z.string().min(1, 'Group name is required'),
  description: z.string().optional(),
});

type AttributeGroupFormValues = z.infer<typeof attributeGroupSchema>;

interface AttributeGroupFormProps {
  defaultValues?: { name: string; description?: string };
  onSubmit: (data: { name: string; description?: string }) => void;
  isPending?: boolean;
  submitLabel?: string;
}

export function AttributeGroupForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Save',
}: AttributeGroupFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<AttributeGroupFormValues>({
    resolver: zodResolver(attributeGroupSchema),
    defaultValues: {
      name: defaultValues?.name || '',
      description: defaultValues?.description || '',
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Name</Label>
        <Input id="name" {...register('name')} placeholder="e.g. Physical Attributes" />
        {errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <textarea
          id="description"
          {...register('description')}
          rows={3}
          placeholder="Optional description..."
          className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>

      <Button type="submit" disabled={isPending} className="w-full">
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}
