import { useParams, Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ArrowLeft } from 'lucide-react';
import { OrderDetail } from '../components/OrderDetail';
import { useOrder } from '../hooks/useOrders';

export default function OrderDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { data: order, isLoading } = useOrder(id!);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!order) {
    return <p className="py-8 text-center text-muted-foreground">Order not found.</p>;
  }

  return (
    <div>
      <Button asChild variant="ghost" size="sm" className="mb-4">
        <Link to="/account/orders">
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Orders
        </Link>
      </Button>
      <OrderDetail order={order} />
    </div>
  );
}
