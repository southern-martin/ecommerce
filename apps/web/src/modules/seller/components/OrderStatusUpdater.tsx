import { useState } from 'react';
import { Button } from '@/shared/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';
import { StatusBadge } from '@/shared/components/data/StatusBadge';
import { ORDER_STATUSES } from '@/shared/lib/constants';
import { Loader2 } from 'lucide-react';

interface OrderStatusUpdaterProps {
  currentStatus: string;
  onUpdate: (status: string) => void;
  isPending: boolean;
}

export function OrderStatusUpdater({ currentStatus, onUpdate, isPending }: OrderStatusUpdaterProps) {
  const [selectedStatus, setSelectedStatus] = useState(currentStatus);

  const handleUpdate = () => {
    if (selectedStatus !== currentStatus) {
      onUpdate(selectedStatus);
    }
  };

  return (
    <div className="flex items-center gap-4">
      <div className="flex items-center gap-2">
        <span className="text-sm text-muted-foreground">Current:</span>
        <StatusBadge status={currentStatus} />
      </div>
      <Select value={selectedStatus} onValueChange={setSelectedStatus}>
        <SelectTrigger className="w-[180px]">
          <SelectValue placeholder="Select status" />
        </SelectTrigger>
        <SelectContent>
          {ORDER_STATUSES.map((status) => (
            <SelectItem key={status} value={status}>
              {status.replace(/_/g, ' ')}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <Button
        onClick={handleUpdate}
        disabled={isPending || selectedStatus === currentStatus}
        size="sm"
      >
        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
        Update
      </Button>
    </div>
  );
}
