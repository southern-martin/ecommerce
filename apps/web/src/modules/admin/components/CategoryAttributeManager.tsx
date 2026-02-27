import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Badge } from '@/shared/components/ui/badge';
import { Button } from '@/shared/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/shared/components/ui/dialog';
import { Plus, X, Tags, Loader2 } from 'lucide-react';
import { adminProductApi } from '../services/admin-product.api';
import type { Attribute } from '../services/admin-product.api';

interface CategoryAttributeManagerProps {
  categoryId: string;
  categoryName: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CategoryAttributeManager({
  categoryId,
  categoryName,
  open,
  onOpenChange,
}: CategoryAttributeManagerProps) {
  const queryClient = useQueryClient();
  const [addingAttrId, setAddingAttrId] = useState<string | null>(null);

  // Fetch all available attribute definitions
  const { data: allAttributes = [], isLoading: allLoading } = useQuery({
    queryKey: ['admin-attributes'],
    queryFn: () => adminProductApi.getAttributes(),
    enabled: open,
  });

  // Fetch currently assigned attributes for this category
  const { data: assignedAttributes = [], isLoading: assignedLoading } = useQuery({
    queryKey: ['category-attributes', categoryId],
    queryFn: () => adminProductApi.getCategoryAttributes(categoryId),
    enabled: open && !!categoryId,
  });

  // Assign mutation
  const assignMutation = useMutation({
    mutationFn: (attributeId: string) =>
      adminProductApi.assignAttribute(categoryId, { attribute_id: attributeId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['category-attributes', categoryId] });
      setAddingAttrId(null);
    },
  });

  // Remove mutation
  const removeMutation = useMutation({
    mutationFn: (attributeId: string) =>
      adminProductApi.removeAttribute(categoryId, attributeId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['category-attributes', categoryId] });
    },
  });

  // Get attribute IDs already assigned
  const assignedIds = new Set(
    assignedAttributes.map((ca: any) => ca.attribute_id || ca.id)
  );

  // Available (not yet assigned) attributes
  const availableAttributes = allAttributes.filter(
    (attr: Attribute) => !assignedIds.has(attr.id)
  );

  const isLoading = allLoading || assignedLoading;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Tags className="h-5 w-5" />
            Manage Attributes: {categoryName}
          </DialogTitle>
          <DialogDescription>
            Assign attribute definitions to this category. Products in this category will require
            values for required attributes.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Currently Assigned */}
          <div>
            <h4 className="mb-2 text-sm font-medium text-muted-foreground">Assigned Attributes</h4>
            {isLoading ? (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <Loader2 className="h-4 w-4 animate-spin" /> Loading...
              </div>
            ) : assignedAttributes.length === 0 ? (
              <p className="text-sm text-muted-foreground">No attributes assigned yet.</p>
            ) : (
              <div className="space-y-2">
                {assignedAttributes.map((ca: any) => {
                  const attr = ca.attribute || allAttributes.find((a: Attribute) => a.id === ca.attribute_id);
                  const attrName = attr?.name || ca.attribute_id?.slice(0, 8) || 'Unknown';
                  const attrType = attr?.type || '';
                  const isRequired = attr?.required;

                  return (
                    <div
                      key={ca.attribute_id || ca.id}
                      className="flex items-center justify-between rounded-lg border px-3 py-2"
                    >
                      <div className="flex items-center gap-2">
                        <span className="text-sm font-medium">{attrName}</span>
                        <Badge variant="outline" className="text-xs">
                          {attrType}
                        </Badge>
                        {isRequired && (
                          <Badge variant="destructive" className="text-xs">
                            Required
                          </Badge>
                        )}
                      </div>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-7 w-7 text-destructive hover:text-destructive"
                        onClick={() => removeMutation.mutate(ca.attribute_id || ca.id)}
                        disabled={removeMutation.isPending}
                      >
                        <X className="h-3.5 w-3.5" />
                      </Button>
                    </div>
                  );
                })}
              </div>
            )}
          </div>

          {/* Add Attribute */}
          <div>
            <h4 className="mb-2 text-sm font-medium text-muted-foreground">Add Attribute</h4>
            {availableAttributes.length === 0 ? (
              <p className="text-sm text-muted-foreground">
                {allAttributes.length === 0
                  ? 'No attribute definitions exist. Create some in the Attributes tab first.'
                  : 'All attributes are already assigned to this category.'}
              </p>
            ) : (
              <div className="flex flex-wrap gap-2">
                {availableAttributes.map((attr: Attribute) => (
                  <Button
                    key={attr.id}
                    variant="outline"
                    size="sm"
                    className="gap-1"
                    onClick={() => {
                      setAddingAttrId(attr.id);
                      assignMutation.mutate(attr.id);
                    }}
                    disabled={assignMutation.isPending && addingAttrId === attr.id}
                  >
                    {assignMutation.isPending && addingAttrId === attr.id ? (
                      <Loader2 className="h-3 w-3 animate-spin" />
                    ) : (
                      <Plus className="h-3 w-3" />
                    )}
                    {attr.name}
                    <Badge variant="secondary" className="ml-1 text-xs">
                      {attr.type}
                    </Badge>
                  </Button>
                ))}
              </div>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
