import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { ProductForm } from '../components/ProductForm';
import { useCreateProduct } from '../hooks/useSellerProducts';

export default function SellerProductNewPage() {
  const navigate = useNavigate();
  const createProduct = useCreateProduct();

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Create New Product</h1>
      <Card>
        <CardHeader>
          <CardTitle>Product Details</CardTitle>
        </CardHeader>
        <CardContent>
          <ProductForm
            onSubmit={(data) =>
              createProduct.mutate({ ...data, images: [] }, {
                onSuccess: (product) => navigate(`/seller/products/${product.id}/edit`),
              })
            }
            isPending={createProduct.isPending}
          />
        </CardContent>
      </Card>
    </div>
  );
}
