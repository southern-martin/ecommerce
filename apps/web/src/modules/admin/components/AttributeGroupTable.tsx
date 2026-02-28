import { Button } from '@/shared/components/ui/button';
import { Badge } from '@/shared/components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/components/ui/table';
import { Pencil, Trash2, Tags } from 'lucide-react';
import { ConfirmDialog } from '@/shared/components/data/ConfirmDialog';
import type { AttributeGroup } from '../services/admin-product.api';

interface AttributeGroupTableProps {
  groups: AttributeGroup[];
  onEdit: (group: AttributeGroup) => void;
  onManageAttributes: (group: AttributeGroup) => void;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

export function AttributeGroupTable({
  groups,
  onEdit,
  onManageAttributes,
  onDelete,
  isDeleting,
}: AttributeGroupTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Description</TableHead>
          <TableHead>Attributes</TableHead>
          <TableHead className="w-[180px]">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {groups.length === 0 ? (
          <TableRow>
            <TableCell colSpan={4} className="text-center text-muted-foreground">
              No attribute groups found.
            </TableCell>
          </TableRow>
        ) : (
          groups.map((group) => (
            <TableRow key={group.id}>
              <TableCell className="font-medium">{group.name}</TableCell>
              <TableCell className="text-muted-foreground">
                {group.description || '-'}
              </TableCell>
              <TableCell>
                <Badge variant="secondary">
                  {group.attributes?.length ?? 0}
                </Badge>
              </TableCell>
              <TableCell>
                <div className="flex gap-1">
                  <Button variant="ghost" size="sm" onClick={() => onEdit(group)}>
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => onManageAttributes(group)}>
                    <Tags className="h-4 w-4" />
                  </Button>
                  <ConfirmDialog
                    title="Delete Attribute Group"
                    description={`Are you sure you want to delete "${group.name}"? This action cannot be undone.`}
                    onConfirm={() => onDelete(group.id)}
                    isPending={isDeleting}
                    trigger={
                      <Button variant="ghost" size="sm">
                        <Trash2 className="h-4 w-4 text-destructive" />
                      </Button>
                    }
                  />
                </div>
              </TableCell>
            </TableRow>
          ))
        )}
      </TableBody>
    </Table>
  );
}
