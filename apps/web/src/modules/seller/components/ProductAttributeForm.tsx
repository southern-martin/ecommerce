import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Badge } from '@/shared/components/ui/badge';
import { Loader2, Save } from 'lucide-react';
import { adminProductApi, type CategoryAttribute } from '@/modules/admin/services/admin-product.api';
import { useProductAttributes, useSetProductAttributes } from '../hooks/useSellerVariants';

interface ProductAttributeFormProps {
  productId: string;
  categoryId: string;
}

export function ProductAttributeForm({ productId, categoryId }: ProductAttributeFormProps) {
  const [values, setValues] = useState<Record<string, { value: string; values: string[] }>>({});

  // Load category attributes (what attributes this category requires)
  const { data: categoryAttributes = [], isLoading: loadingCategoryAttrs } = useQuery({
    queryKey: ['category-attributes', categoryId],
    queryFn: () => adminProductApi.getCategoryAttributes(categoryId),
    enabled: !!categoryId,
  });

  // Load existing product attribute values
  const { data: productAttributes = [], isLoading: loadingProductAttrs } = useProductAttributes(productId);

  const setProductAttributes = useSetProductAttributes();

  // Initialize form values from existing product attributes
  useEffect(() => {
    if (productAttributes.length > 0) {
      const initial: Record<string, { value: string; values: string[] }> = {};
      productAttributes.forEach((attr) => {
        initial[attr.attribute_id] = {
          value: attr.value || '',
          values: attr.values || [],
        };
      });
      setValues(initial);
    }
  }, [productAttributes]);

  const isLoading = loadingCategoryAttrs || loadingProductAttrs;

  if (!categoryId) {
    return (
      <Card>
        <CardContent className="py-6 text-center text-muted-foreground">
          Please select a category first to see available attributes.
        </CardContent>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <Loader2 className="h-6 w-6 animate-spin" />
      </div>
    );
  }

  if (categoryAttributes.length === 0) {
    return (
      <Card>
        <CardContent className="py-6 text-center text-muted-foreground">
          No attributes defined for this category.
        </CardContent>
      </Card>
    );
  }

  const handleSave = () => {
    const attrs = categoryAttributes
      .map((ca) => ({
        attribute_id: ca.attribute_id,
        value: values[ca.attribute_id]?.value || '',
        values: values[ca.attribute_id]?.values || [],
      }))
      .filter((a) => a.value || a.values.length > 0);

    setProductAttributes.mutate({ productId, attributes: attrs });
  };

  const updateValue = (attrId: string, value: string) => {
    setValues((prev) => ({
      ...prev,
      [attrId]: { ...prev[attrId], value, values: prev[attrId]?.values || [] },
    }));
  };

  const toggleMultiValue = (attrId: string, option: string) => {
    setValues((prev) => {
      const current = prev[attrId]?.values || [];
      const updated = current.includes(option)
        ? current.filter((v) => v !== option)
        : [...current, option];
      return {
        ...prev,
        [attrId]: { ...prev[attrId], value: prev[attrId]?.value || '', values: updated },
      };
    });
  };

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-base">Product Attributes</CardTitle>
        <Button size="sm" onClick={handleSave} disabled={setProductAttributes.isPending}>
          {setProductAttributes.isPending ? (
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
          ) : (
            <Save className="mr-2 h-4 w-4" />
          )}
          Save Attributes
        </Button>
      </CardHeader>
      <CardContent className="space-y-4">
        {categoryAttributes.map((ca: CategoryAttribute) => {
          const attr = ca.attribute;
          const attrType = attr.type;
          const currentValue = values[attr.id]?.value || '';
          const currentValues = values[attr.id]?.values || [];

          return (
            <div key={attr.id} className="space-y-1.5">
              <Label className="flex items-center gap-2">
                {attr.name}
                {attr.required && (
                  <Badge variant="destructive" className="text-xs px-1.5 py-0">
                    Required
                  </Badge>
                )}
                {attr.unit && (
                  <span className="text-xs text-muted-foreground">({attr.unit})</span>
                )}
              </Label>

              {(attrType === 'text' || attrType === 'number') && (
                <Input
                  type={attrType === 'number' ? 'number' : 'text'}
                  value={currentValue}
                  onChange={(e) => updateValue(attr.id, e.target.value)}
                  placeholder={`Enter ${attr.name.toLowerCase()}`}
                />
              )}

              {attrType === 'select' && attr.options && (
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={currentValue}
                  onChange={(e) => updateValue(attr.id, e.target.value)}
                >
                  <option value="">Select {attr.name}</option>
                  {attr.options.map((opt) => (
                    <option key={opt} value={opt}>
                      {opt}
                    </option>
                  ))}
                </select>
              )}

              {attrType === 'multi_select' && attr.options && (
                <div className="flex flex-wrap gap-2">
                  {attr.options.map((opt) => {
                    const isSelected = currentValues.includes(opt);
                    return (
                      <Button
                        key={opt}
                        variant={isSelected ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => toggleMultiValue(attr.id, opt)}
                      >
                        {opt}
                      </Button>
                    );
                  })}
                </div>
              )}

              {(attrType === 'boolean' || attrType === 'bool') && (
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id={`attr-${attr.id}`}
                    checked={currentValue === 'true'}
                    onChange={(e) => updateValue(attr.id, e.target.checked ? 'true' : 'false')}
                    className="h-4 w-4 rounded border-gray-300"
                  />
                  <Label htmlFor={`attr-${attr.id}`} className="text-sm font-normal">
                    {attr.name}
                  </Label>
                </div>
              )}

              {attrType === 'color' && (
                <div className="flex items-center gap-2">
                  <input
                    type="color"
                    value={currentValue || '#000000'}
                    onChange={(e) => updateValue(attr.id, e.target.value)}
                    className="h-10 w-10 cursor-pointer rounded border"
                  />
                  <Input
                    value={currentValue}
                    onChange={(e) => updateValue(attr.id, e.target.value)}
                    placeholder="#000000"
                    className="w-32"
                  />
                </div>
              )}
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}
