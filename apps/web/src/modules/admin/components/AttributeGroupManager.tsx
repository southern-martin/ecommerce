import { useState } from 'react';
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
import { useAdminAttributes, useGroupAttributes, useAssignAttributeToGroup, useRemoveAttributeFromGroup } from '../hooks/useAdminProducts';
import type { Attribute } from '../services/admin-product.api';

interface AttributeGroupManagerProps {
  groupId: string;
  groupName: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function AttributeGroupManager({
  groupId,
  groupName,
  open,
  onOpenChange,
}: AttributeGroupManagerProps) {
  const [addingAttrId, setAddingAttrId] = useState<string | null>(null);

  // Fetch all available attribute definitions
  const { data: allAttributes = [], isLoading: allLoading } = useAdminAttributes();

  // Fetch currently assigned attributes for this group
  const { data: groupAttributes = [], isLoading: groupLoading } = useGroupAttributes(groupId);

  // Assign mutation
  const assignMutation = useAssignAttributeToGroup();

  // Remove mutation
  const removeMutation = useRemoveAttributeFromGroup();

  // Get attribute IDs already assigned
  const assignedIds = new Set(
    groupAttributes.map((attr: Attribute) => attr.id)
  );

  // Available (not yet assigned) attributes
  const availableAttributes = allAttributes.filter(
    (attr: Attribute) => !assignedIds.has(attr.id)
  );

  const isLoading = allLoading || groupLoading;

  const handleAssign = (attrId: string) => {
    setAddingAttrId(attrId);
    assignMutation.mutate(
      { groupId, attribute_id: attrId },
      { onSettled: () => setAddingAttrId(null) }
    );
  };

  const handleRemove = (attrId: string) => {
    removeMutation.mutate({ groupId, attrId });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Tags className="h-5 w-5" />
            Manage Attributes: {groupName}
          </DialogTitle>
          <DialogDescription>
            Assign attributes to this group. Grouped attributes are displayed together on product pages.
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
            ) : groupAttributes.length === 0 ? (
              <p className="text-sm text-muted-foreground">No attributes assigned yet.</p>
            ) : (
              <div className="space-y-2">
                {groupAttributes.map((attr: Attribute) => (
                  <div
                    key={attr.id}
                    className="flex items-center justify-between rounded-lg border px-3 py-2"
                  >
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium">{attr.name}</span>
                      <Badge variant="outline" className="text-xs">
                        {attr.type}
                      </Badge>
                      {attr.required && (
                        <Badge variant="destructive" className="text-xs">
                          Required
                        </Badge>
                      )}
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-7 w-7 text-destructive hover:text-destructive"
                      onClick={() => handleRemove(attr.id)}
                      disabled={removeMutation.isPending}
                    >
                      <X className="h-3.5 w-3.5" />
                    </Button>
                  </div>
                ))}
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
                  : 'All attributes are already assigned to this group.'}
              </p>
            ) : (
              <div className="flex flex-wrap gap-2">
                {availableAttributes.map((attr: Attribute) => (
                  <Button
                    key={attr.id}
                    variant="outline"
                    size="sm"
                    className="gap-1"
                    onClick={() => handleAssign(attr.id)}
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
