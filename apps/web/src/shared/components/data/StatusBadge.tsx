import { Badge } from '@/shared/components/ui/badge';
import { cn } from '@/shared/lib/utils';

const statusColorMap: Record<string, string> = {
  // Order statuses
  pending: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  confirmed: 'bg-blue-100 text-blue-800 border-blue-200',
  processing: 'bg-indigo-100 text-indigo-800 border-indigo-200',
  shipped: 'bg-purple-100 text-purple-800 border-purple-200',
  delivered: 'bg-green-100 text-green-800 border-green-200',
  cancelled: 'bg-red-100 text-red-800 border-red-200',
  refunded: 'bg-gray-100 text-gray-800 border-gray-200',
  // Return statuses
  requested: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  approved: 'bg-green-100 text-green-800 border-green-200',
  rejected: 'bg-red-100 text-red-800 border-red-200',
  shipped_back: 'bg-purple-100 text-purple-800 border-purple-200',
  received: 'bg-blue-100 text-blue-800 border-blue-200',
  // Generic
  active: 'bg-green-100 text-green-800 border-green-200',
  inactive: 'bg-gray-100 text-gray-800 border-gray-200',
  expired: 'bg-red-100 text-red-800 border-red-200',
  draft: 'bg-gray-100 text-gray-800 border-gray-200',
  published: 'bg-green-100 text-green-800 border-green-200',
  resolved: 'bg-green-100 text-green-800 border-green-200',
  open: 'bg-yellow-100 text-yellow-800 border-yellow-200',
  paid: 'bg-green-100 text-green-800 border-green-200',
};

interface StatusBadgeProps {
  status: string;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const colorClass = statusColorMap[status] || 'bg-gray-100 text-gray-800 border-gray-200';
  return (
    <Badge variant="outline" className={cn(colorClass, 'capitalize', className)}>
      {status.replace(/_/g, ' ')}
    </Badge>
  );
}
