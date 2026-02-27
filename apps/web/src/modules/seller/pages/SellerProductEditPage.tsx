import { useParams } from 'react-router-dom';
import { Loader2, RefreshCw } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import { Button } from '@/shared/components/ui/button';
import { ProductForm } from '../components/ProductForm';
import { ProductOptionManager } from '../components/ProductOptionManager';
import { VariantTable } from '../components/VariantTable';
import { ProductAttributeForm } from '../components/ProductAttributeForm';
import { useUpdateProduct } from '../hooks/useSellerProducts';
import { useGenerateVariants } from '../hooks/useSellerVariants';
import { useQuery } from '@tanstack/react-query';
import { sellerProductApi } from '../services/seller-product.api';

export default function SellerProductEditPage() {
  const { id } = useParams<{ id: string }>();
  const updateProduct = useUpdateProduct();
  const generateVariants = useGenerateVariants();

  const { data: product, isLoading } = useQuery({
    queryKey: ['seller-products', id],
    queryFn: () => sellerProductApi.getProductById(id!),
    enabled: !!id,
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <Loader2 className="h-8 w-8 animate-spin" />
      </div>
    );
  }

  if (!product) {
    return <p className="text-center text-muted-foreground">Product not found</p>;
  }

  const isConfigurable = product.product_type === 'configurable';

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Edit Product: {product.name}</h1>

      <Tabs defaultValue="basic">
        <TabsList>
          <TabsTrigger value="basic">Basic Info</TabsTrigger>
          {isConfigurable && (
            <TabsTrigger value="variants">Options & Variants</TabsTrigger>
          )}
          <TabsTrigger value="attributes">Attributes</TabsTrigger>
        </TabsList>

        <TabsContent value="basic" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Product Details</CardTitle>
            </CardHeader>
            <CardContent>
              <ProductForm
                defaultValues={{
                  name: product.name,
                  description: product.description || '',
                  price: product.base_price_cents ? product.base_price_cents / 100 : 0,
                  compare_at_price: 0,
                  category_id: product.category_id || '',
                  product_type: product.product_type || 'simple',
                  stock_quantity: product.stock_quantity ?? 0,
                }}
                onSubmit={(data) =>
                  updateProduct.mutate({
                    id: product.id,
                    data: {
                      name: data.name,
                      description: data.description,
                      category_id: data.category_id,
                      base_price_cents: Math.round(data.price * 100),
                      stock_quantity: product.product_type === 'simple' ? (data.stock_quantity ?? 0) : undefined,
                    },
                  })
                }
                isPending={updateProduct.isPending}
                submitLabel="Update Product"
              />
            </CardContent>
          </Card>
        </TabsContent>

        {isConfigurable && (
          <TabsContent value="variants" className="mt-4 space-y-4">
            <ProductOptionManager
              productId={product.id}
              options={product.options || []}
            />

            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="text-base">Variants</CardTitle>
                <Button
                  size="sm"
                  onClick={() => generateVariants.mutate(product.id)}
                  disabled={generateVariants.isPending}
                >
                  <RefreshCw className="mr-2 h-4 w-4" />
                  {generateVariants.isPending ? 'Generating...' : 'Generate Variants'}
                </Button>
              </CardHeader>
              <CardContent>
                <VariantTable
                  productId={product.id}
                  variants={product.variants || []}
                />
              </CardContent>
            </Card>
          </TabsContent>
        )}

        <TabsContent value="attributes" className="mt-4">
          <ProductAttributeForm
            productId={product.id}
            categoryId={product.category_id}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
