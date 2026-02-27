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
            showProductTypeSelector
            onSubmit={(data) =>
              createProduct.mutate(
                {
                  name: data.name,
                  description: data.description,
                  category_id: data.category_id,
                  base_price_cents: Math.round(data.price * 100),
                  product_type: data.product_type || 'simple',
                  stock_quantity: data.product_type === 'simple' ? (data.stock_quantity ?? 0) : 0,
                  image_urls: [],
                },
                {
                  onSuccess: (product) => navigate(`/seller/products/${product.id}/edit`),
                }
              )
            }
            isPending={createProduct.isPending}
          />
        </CardContent>
      </Card>
    </div>
  );
}
