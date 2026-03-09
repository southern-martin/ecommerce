import { Star } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/components/ui/card';

interface Product {
  name: string;
  unitsSold: number;
  revenue: number;
  rating: number;
}

const products: Product[] = [
  { name: 'Wireless Noise-Cancelling Headphones', unitsSold: 1243, revenue: 186450, rating: 4.8 },
  { name: 'Ultra-Slim Laptop Stand', unitsSold: 987, revenue: 49350, rating: 4.6 },
  { name: 'Mechanical Keyboard RGB', unitsSold: 856, revenue: 102720, rating: 4.7 },
  { name: 'Portable Bluetooth Speaker', unitsSold: 734, revenue: 51380, rating: 4.5 },
  { name: 'Smart Fitness Tracker', unitsSold: 698, revenue: 83760, rating: 4.4 },
];

function formatCurrency(value: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(value);
}

export default function TopProductsTable() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Top Selling Products</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b text-left text-muted-foreground">
                <th className="pb-3 font-medium">Product Name</th>
                <th className="pb-3 font-medium text-right">Units Sold</th>
                <th className="pb-3 font-medium text-right">Revenue</th>
                <th className="pb-3 font-medium text-right">Rating</th>
              </tr>
            </thead>
            <tbody>
              {products.map((product, index) => (
                <tr
                  key={product.name}
                  className={index % 2 === 0 ? 'bg-muted/50' : ''}
                >
                  <td className="py-3 pr-4 font-medium">{product.name}</td>
                  <td className="py-3 text-right tabular-nums">
                    {product.unitsSold.toLocaleString()}
                  </td>
                  <td className="py-3 text-right tabular-nums">
                    {formatCurrency(product.revenue)}
                  </td>
                  <td className="py-3 text-right">
                    <span className="inline-flex items-center gap-1">
                      <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                      {product.rating}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}
