import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Card, CardContent } from '@/shared/components/ui/card';
import { ArrowRight } from 'lucide-react';
import { ProductGrid } from '../components/ProductGrid';
import { productApi } from '../services/product.api';
import { categoryApi } from '../services/category.api';

export default function HomePage() {
  const { data: featured, isLoading: featuredLoading } = useQuery({
    queryKey: ['products', 'featured'],
    queryFn: () => productApi.getFeaturedProducts(),
  });

  const { data: trending, isLoading: trendingLoading } = useQuery({
    queryKey: ['products', 'trending'],
    queryFn: () => productApi.getTrendingProducts(),
  });

  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.getCategories(),
  });

  return (
    <div className="space-y-12">
      {/* Hero Banner */}
      <section className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-primary to-primary/80 px-8 py-16 text-primary-foreground md:px-16 md:py-24">
        <div className="relative z-10 max-w-xl">
          <h1 className="text-4xl font-bold tracking-tight md:text-5xl">
            Discover Amazing Products
          </h1>
          <p className="mt-4 text-lg text-primary-foreground/80">
            Shop the latest trends with unbeatable prices and fast shipping.
          </p>
          <Button asChild size="lg" variant="secondary" className="mt-8">
            <Link to="/products">
              Shop Now <ArrowRight className="ml-2 h-5 w-5" />
            </Link>
          </Button>
        </div>
      </section>

      {/* Featured Categories */}
      {categories && categories.length > 0 && (
        <section>
          <div className="mb-6 flex items-center justify-between">
            <h2 className="text-2xl font-bold">Shop by Category</h2>
            <Link to="/categories" className="text-sm text-primary hover:underline">
              View All
            </Link>
          </div>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6">
            {categories.slice(0, 6).map((category) => (
              <Link key={category.id} to={`/products?category=${category.slug}`}>
                <Card className="transition-shadow hover:shadow-md">
                  <CardContent className="flex flex-col items-center p-4">
                    {category.image_url && (
                      <img
                        src={category.image_url}
                        alt={category.name}
                        className="mb-3 h-16 w-16 rounded-full object-cover"
                      />
                    )}
                    <span className="text-center text-sm font-medium">{category.name}</span>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Trending Products */}
      <section>
        <div className="mb-6 flex items-center justify-between">
          <h2 className="text-2xl font-bold">Trending Now</h2>
          <Link to="/products?sort=popular" className="text-sm text-primary hover:underline">
            View All
          </Link>
        </div>
        <ProductGrid products={trending ?? []} isLoading={trendingLoading} />
      </section>

      {/* Featured Products */}
      <section>
        <div className="mb-6 flex items-center justify-between">
          <h2 className="text-2xl font-bold">Featured Products</h2>
          <Link to="/products" className="text-sm text-primary hover:underline">
            View All
          </Link>
        </div>
        <ProductGrid products={featured ?? []} isLoading={featuredLoading} />
      </section>
    </div>
  );
}
