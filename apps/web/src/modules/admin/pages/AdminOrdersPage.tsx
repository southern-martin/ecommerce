import { useState } from 'react';
import { useQuery, keepPreviousData } from '@tanstack/react-query';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import apiClient from '@/shared/lib/api-client';
import type { PaginatedResponse } from '@/shared/types/api.types';
import type { Order } from '@/modules/checkout/services/order.api';
import { SellerOrderTable } from '@/modules/seller/components/SellerOrderTable';

export default function AdminOrdersPage() {
  const [page, setPage] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ['admin-orders', page],
    queryFn: async () => {
      const response = await apiClient.get<PaginatedResponse<Order>>('/admin/orders', {
        params: { page, page_size: 20 },
      });
      return response.data;
    },
    placeholderData: keepPreviousData,
  });

  const totalPages = data ? Math.ceil(data.total / data.page_size) : 0;

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">All Orders</h1>

      {data && data.data.length > 0 ? (
        <>
          <SellerOrderTable orders={data.data} />
          {totalPages > 1 && (
            <div className="mt-6 flex items-center justify-center gap-2">
              <Button variant="outline" size="sm" disabled={page === 1} onClick={() => setPage((p) => p - 1)}>
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <span className="text-sm text-muted-foreground">Page {page} of {totalPages}</span>
              <Button variant="outline" size="sm" disabled={page === totalPages} onClick={() => setPage((p) => p + 1)}>
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          )}
        </>
      ) : (
        <p className="py-8 text-center text-muted-foreground">No orders found.</p>
      )}
    </div>
  );
}
