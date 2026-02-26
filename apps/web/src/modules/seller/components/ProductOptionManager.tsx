import { useState } from 'react';
import { Plus, X } from 'lucide-react';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Badge } from '@/shared/components/ui/badge';
import { useAddOption, useRemoveOption } from '../hooks/useSellerVariants';
import type { ProductOption } from '../services/seller-variant.api';

interface ProductOptionManagerProps {
  productId: string;
  options: ProductOption[];
}

export function ProductOptionManager({ productId, options }: ProductOptionManagerProps) {
  const [optionName, setOptionName] = useState('');
  const [optionValues, setOptionValues] = useState('');
  const addOption = useAddOption();
  const removeOption = useRemoveOption();

  const handleAddOption = () => {
    if (!optionName.trim() || !optionValues.trim()) return;
    const values = optionValues.split(',').map((v) => v.trim()).filter(Boolean);
    addOption.mutate(
      { productId, data: { name: optionName.trim(), values } },
      {
        onSuccess: () => {
          setOptionName('');
          setOptionValues('');
        },
      }
    );
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Add Option</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div>
            <Label>Option Name</Label>
            <Input
              placeholder="e.g. Color, Size"
              value={optionName}
              onChange={(e) => setOptionName(e.target.value)}
            />
          </div>
          <div>
            <Label>Values (comma separated)</Label>
            <Input
              placeholder="e.g. Red, Blue, Green"
              value={optionValues}
              onChange={(e) => setOptionValues(e.target.value)}
            />
          </div>
          <Button onClick={handleAddOption} disabled={addOption.isPending} size="sm">
            <Plus className="mr-2 h-4 w-4" />
            {addOption.isPending ? 'Adding...' : 'Add Option'}
          </Button>
        </CardContent>
      </Card>

      {options.length > 0 && (
        <div className="space-y-2">
          <h4 className="text-sm font-medium">Current Options</h4>
          {options.map((option) => (
            <div key={option.id} className="flex items-center justify-between rounded-lg border p-3">
              <div>
                <span className="font-medium">{option.name}:</span>
                <span className="ml-2">
                  {option.values.map((v) => (
                    <Badge key={v} variant="secondary" className="mr-1">
                      {v}
                    </Badge>
                  ))}
                </span>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => removeOption.mutate({ productId, optionId: option.id })}
                disabled={removeOption.isPending}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
