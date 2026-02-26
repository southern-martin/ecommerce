import { useState } from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Skeleton } from '@/shared/components/ui/skeleton';
import { Plus, ChevronLeft, ChevronRight } from 'lucide-react';
import { SellerProductTable } from '../components/SellerProductTable';
import { useSellerProducts, useDeleteProduct } from '../hooks/useSellerProducts';

export default function SellerProductsPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useSellerProducts(page);
  const deleteProduct = useDeleteProduct();
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
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold">My Products</h1>
        <Button asChild>
          <Link to="/seller/products/new">
            <Plus className="mr-2 h-4 w-4" />
            Add Product
          </Link>
        </Button>
      </div>

      {data && data.data.length > 0 ? (
        <>
          <SellerProductTable
            products={data.data}
            onDelete={(id) => deleteProduct.mutate(id)}
          />
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
        <p className="py-8 text-center text-muted-foreground">No products yet. Create your first product.</p>
      )}
    </div>
  );
}
