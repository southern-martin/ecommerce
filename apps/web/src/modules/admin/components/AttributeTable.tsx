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
import { Pencil, Trash2 } from 'lucide-react';
import { ConfirmDialog } from '@/shared/components/data/ConfirmDialog';

interface Attribute {
  id: string;
  name: string;
  type: string;
  required: boolean;
  filterable: boolean;
  options?: string[];
}

interface AttributeTableProps {
  attributes: Attribute[];
  onEdit: (attribute: Attribute) => void;
  onDelete: (id: string) => void;
  isDeleting?: boolean;
}

export function AttributeTable({ attributes, onEdit, onDelete, isDeleting }: AttributeTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Flags</TableHead>
          <TableHead>Options</TableHead>
          <TableHead className="w-[120px]">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {attributes.length === 0 ? (
          <TableRow>
            <TableCell colSpan={5} className="text-center text-muted-foreground">
              No attributes found.
            </TableCell>
          </TableRow>
        ) : (
          attributes.map((attr) => (
            <TableRow key={attr.id}>
              <TableCell className="font-medium">{attr.name}</TableCell>
              <TableCell className="capitalize">{attr.type}</TableCell>
              <TableCell>
                <div className="flex gap-1">
                  {attr.required && (
                    <Badge variant="outline" className="bg-blue-50 text-blue-700">
                      Required
                    </Badge>
                  )}
                  {attr.filterable && (
                    <Badge variant="outline" className="bg-green-50 text-green-700">
                      Filterable
                    </Badge>
                  )}
                </div>
              </TableCell>
              <TableCell>{attr.options?.length ?? 0}</TableCell>
              <TableCell>
                <div className="flex gap-1">
                  <Button variant="ghost" size="sm" onClick={() => onEdit(attr)}>
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <ConfirmDialog
                    title="Delete Attribute"
                    description={`Are you sure you want to delete "${attr.name}"? This action cannot be undone.`}
                    onConfirm={() => onDelete(attr.id)}
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
