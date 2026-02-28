import { useState } from 'react';
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
import { Loader2, Plus, X } from 'lucide-react';

const attributeSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  type: z.enum(['text', 'number', 'select', 'multi_select', 'color', 'boolean']),
  required: z.boolean().default(false),
  filterable: z.boolean().default(false),
});

type AttributeFormValues = z.infer<typeof attributeSchema>;

interface OptionValueRow {
  value: string;
  color_hex: string;
}

interface AttributeFormProps {
  defaultValues?: Partial<AttributeFormValues>;
  defaultOptionValues?: { value: string; color_hex?: string }[];
  onSubmit: (data: {
    name: string;
    type: 'text' | 'number' | 'select' | 'boolean' | 'multi_select' | 'color';
    required?: boolean;
    filterable?: boolean;
    option_values?: { value: string; color_hex?: string; sort_order?: number }[];
  }) => void;
  isPending?: boolean;
  submitLabel?: string;
}

export function AttributeForm({
  defaultValues,
  defaultOptionValues,
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
      ...defaultValues,
    },
  });

  const type = watch('type');
  const showOptions = type === 'select' || type === 'multi_select' || type === 'color';

  const [optionRows, setOptionRows] = useState<OptionValueRow[]>(
    defaultOptionValues?.map((ov) => ({ value: ov.value, color_hex: ov.color_hex || '' })) || [{ value: '', color_hex: '' }]
  );

  const addOptionRow = () => {
    setOptionRows([...optionRows, { value: '', color_hex: '' }]);
  };

  const removeOptionRow = (index: number) => {
    setOptionRows(optionRows.filter((_, i) => i !== index));
  };

  const updateOptionRow = (index: number, field: keyof OptionValueRow, val: string) => {
    const updated = [...optionRows];
    updated[index] = { ...updated[index], [field]: val };
    setOptionRows(updated);
  };

  const handleFormSubmit = (values: AttributeFormValues) => {
    const optionValues = showOptions
      ? optionRows
          .filter((row) => row.value.trim() !== '')
          .map((row, i) => ({
            value: row.value.trim(),
            color_hex: row.color_hex || undefined,
            sort_order: i,
          }))
      : undefined;

    onSubmit({
      ...values,
      option_values: optionValues,
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
            <SelectItem value="multi_select">Multi-Select</SelectItem>
            <SelectItem value="color">Color</SelectItem>
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

      {showOptions && (
        <div className="space-y-2">
          <Label>Option Values</Label>
          <div className="space-y-2">
            {optionRows.map((row, index) => (
              <div key={index} className="flex items-center gap-2">
                <Input
                  value={row.value}
                  onChange={(e) => updateOptionRow(index, 'value', e.target.value)}
                  placeholder={`Option ${index + 1}`}
                  className="flex-1"
                />
                {type === 'color' && (
                  <input
                    type="color"
                    value={row.color_hex || '#000000'}
                    onChange={(e) => updateOptionRow(index, 'color_hex', e.target.value)}
                    className="h-10 w-10 cursor-pointer rounded border"
                  />
                )}
                {optionRows.length > 1 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => removeOptionRow(index)}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                )}
              </div>
            ))}
          </div>
          <Button type="button" variant="outline" size="sm" onClick={addOptionRow}>
            <Plus className="mr-1 h-3 w-3" /> Add Option
          </Button>
        </div>
      )}

      <Button type="submit" disabled={isPending}>
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        {submitLabel}
      </Button>
    </form>
  );
}
