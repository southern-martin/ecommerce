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

const attributeSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  type: z.enum(['text', 'number', 'select', 'boolean']),
  required: z.boolean().default(false),
  filterable: z.boolean().default(false),
  options: z.string().optional(),
});

type AttributeFormValues = z.infer<typeof attributeSchema>;

interface AttributeFormProps {
  defaultValues?: Partial<AttributeFormValues>;
  onSubmit: (data: {
    name: string;
    type: 'text' | 'number' | 'select' | 'boolean';
    required?: boolean;
    filterable?: boolean;
    options?: string[];
  }) => void;
  isPending?: boolean;
  submitLabel?: string;
}

export function AttributeForm({
  defaultValues,
  onSubmit,
  isPending,
  submitLabel = 'Create Attribute',
}: AttributeFormProps) {
  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<AttributeFormValues>({
    resolver: zodResolver(attributeSchema),
    defaultValues: {
      name: '',
      type: 'text',
      required: false,
      filterable: false,
      options: '',
      ...defaultValues,
    },
  });

  const type = watch('type');

  const handleFormSubmit = (values: AttributeFormValues) => {
    const { options, ...rest } = values;
    onSubmit({
      ...rest,
      options:
        values.type === 'select' && options
          ? options.split(',').map((o) => o.trim()).filter(Boolean)
          : undefined,
    });
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="attr-name">Name</Label>
        <Input id="attr-name" {...register('name')} />
        {errors.name && <p className="text-sm text-destructive">{errors.name.message}</p>}
      </div>

      <div className="space-y-2">
        <Label>Type</Label>
        <Select
          defaultValue={defaultValues?.type || 'text'}
          onValueChange={(value) => setValue('type', value as any)}
        >
          <SelectTrigger>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="text">Text</SelectItem>
            <SelectItem value="number">Number</SelectItem>
            <SelectItem value="select">Select</SelectItem>
            <SelectItem value="boolean">Boolean</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <div className="flex items-center gap-6">
        <label className="flex items-center gap-2 text-sm">
          <input type="checkbox" {...register('required')} className="rounded" />
          Required
        </label>
        <label className="flex items-center gap-2 text-sm">
          <input type="checkbox" {...register('filterable')} className="rounded" />
          Filterable
        </label>
      </div>

      {type === 'select' && (
        <div className="space-y-2">
          <Label htmlFor="options">Options (comma-separated)</Label>
          <Input
            id="options"
            {...register('options')}
            placeholder="e.g. Red, Blue, Green"
          />
        </div>
      )}

      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}
