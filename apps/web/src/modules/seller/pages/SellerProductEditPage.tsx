import { useParams } from 'react-router-dom';
import { Loader2, RefreshCw, Package, Settings2, Tags } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/shared/components/ui/tabs';
import { Badge } from '@/shared/components/ui/badge';
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
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold">Edit Product: {product.name}</h1>
        <Badge variant={isConfigurable ? 'default' : 'outline'}>
          {isConfigurable ? 'Configurable' : 'Simple'}
        </Badge>
      </div>

      <Tabs defaultValue="basic">
        <TabsList className={`grid w-full max-w-lg ${isConfigurable ? 'grid-cols-3' : 'grid-cols-2'}`}>
          <TabsTrigger value="basic" className="flex items-center gap-2">
            <Package className="h-4 w-4" />
            Basic Info
          </TabsTrigger>
          {isConfigurable && (
            <TabsTrigger value="variants" className="flex items-center gap-2">
              <Settings2 className="h-4 w-4" />
              Options & Variants
            </TabsTrigger>
          )}
          <TabsTrigger value="attributes" className="flex items-center gap-2">
            <Tags className="h-4 w-4" />
            Attributes
          </TabsTrigger>
        </TabsList>

        <TabsContent value="basic" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Product Details</CardTitle>
              <CardDescription>
                Update the basic information for this product.
                {isConfigurable
                  ? ' Price and stock are managed per variant in the Options & Variants tab.'
                  : ' Set the price and stock quantity directly.'}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ProductForm
                defaultValues={{
                  name: product.name,
                  description: product.description || '',
                  price: product.base_price_cents ? product.base_price_cents / 100 : 0,
                  compare_at_price: 0,
                  category_id: product.category_id || '',
                  attribute_group_id: product.attribute_group_id || '',
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
                      attribute_group_id: data.attribute_group_id || undefined,
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
            <Card>
              <CardHeader>
                <CardTitle>Product Options</CardTitle>
                <CardDescription>
                  Define options like Size, Color, or Material. Each option can have multiple values.
                  After adding options, generate variants to create all possible combinations.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <ProductOptionManager
                  productId={product.id}
                  options={product.options || []}
                />
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle className="text-base">Variants</CardTitle>
                  <CardDescription className="mt-1">
                    Each variant has its own price, stock, and SKU. Generate variants from options above.
                  </CardDescription>
                </div>
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
          <Card>
            <CardHeader>
              <CardTitle>Product Attributes</CardTitle>
              <CardDescription>
                Set specification attributes for this product (e.g., Brand, Material, Weight).
                Available attributes depend on the product&apos;s attribute group.
                {!product.attribute_group_id && ' Select an attribute group first in the Basic Info tab.'}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ProductAttributeForm
                productId={product.id}
                attributeGroupId={product.attribute_group_id}
              />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
