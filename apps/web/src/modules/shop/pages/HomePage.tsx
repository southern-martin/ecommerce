import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import {
  ArrowRight,
  Truck,
  ShieldCheck,
  RefreshCcw,
  Headphones,
  ChevronRight,
  Sparkles,
  TrendingUp,
  Zap,
} from 'lucide-react';
import { ProductGrid } from '../components/ProductGrid';
import { productApi } from '../services/product.api';
import { categoryApi } from '../services/category.api';

const CATEGORY_IMAGES: Record<string, string> = {
  electronics: 'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=400&h=400&fit=crop',
  clothing: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=400&h=400&fit=crop',
  apparel: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=400&h=400&fit=crop',
  gadgets: 'https://images.unsplash.com/photo-1519389950473-47ba0277781c?w=400&h=400&fit=crop',
  peripherals: 'https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=400&h=400&fit=crop',
  default: 'https://images.unsplash.com/photo-1472851294608-062f824d29cc?w=400&h=400&fit=crop',
};

function getCategoryImage(slug: string): string {
  for (const [key, url] of Object.entries(CATEGORY_IMAGES)) {
    if (slug.toLowerCase().includes(key)) return url;
  }
  return CATEGORY_IMAGES.default;
}

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
    <div className="space-y-0">
      {/* ── Hero Banner ── */}
      <section className="relative overflow-hidden bg-gradient-to-br from-violet-600 via-purple-600 to-indigo-700">
        <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=1920&q=80')] bg-cover bg-center opacity-20" />
        <div className="absolute inset-0 bg-gradient-to-r from-violet-900/80 to-transparent" />
        <div className="relative mx-auto max-w-7xl px-4 py-20 sm:px-6 sm:py-28 lg:px-8 lg:py-36">
          <div className="max-w-2xl">
            <span className="mb-4 inline-flex items-center gap-2 rounded-full bg-white/15 px-4 py-1.5 text-sm font-medium text-white backdrop-blur-sm">
              <Sparkles className="h-4 w-4" />
              New Season Collection
            </span>
            <h1 className="text-4xl font-extrabold tracking-tight text-white sm:text-5xl lg:text-6xl">
              Discover Products
              <br />
              <span className="bg-gradient-to-r from-amber-200 to-yellow-400 bg-clip-text text-transparent">
                You'll Love
              </span>
            </h1>
            <p className="mt-6 max-w-lg text-lg text-violet-100">
              Premium quality products at unbeatable prices. Free shipping on orders over $50.
            </p>
            <div className="mt-8 flex flex-wrap gap-4">
              <Button asChild size="lg" className="bg-white text-violet-700 hover:bg-violet-50 shadow-xl shadow-violet-900/30 font-semibold px-8">
                <Link to="/products">
                  Shop Now <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button asChild size="lg" variant="outline" className="border-white/30 text-white hover:bg-white/10 backdrop-blur-sm font-semibold px-8">
                <Link to="/promotions">
                  <Zap className="mr-2 h-5 w-5" />
                  View Deals
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* ── Trust Bar ── */}
      <section className="border-b bg-muted/50">
        <div className="container mx-auto grid grid-cols-2 gap-4 px-4 py-6 md:grid-cols-4">
          {[
            { icon: Truck, title: 'Free Shipping', desc: 'On orders over $50' },
            { icon: ShieldCheck, title: 'Secure Payment', desc: '100% protected' },
            { icon: RefreshCcw, title: 'Easy Returns', desc: '30-day guarantee' },
            { icon: Headphones, title: '24/7 Support', desc: 'Expert help' },
          ].map(({ icon: Icon, title, desc }) => (
            <div key={title} className="flex items-center gap-3 px-2">
              <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary">
                <Icon className="h-5 w-5" />
              </div>
              <div>
                <p className="text-sm font-semibold">{title}</p>
                <p className="text-xs text-muted-foreground">{desc}</p>
              </div>
            </div>
          ))}
        </div>
      </section>

      <div className="container mx-auto space-y-16 px-4 py-12">
        {/* ── Shop by Category ── */}
        {categories && categories.length > 0 && (
          <section>
            <div className="mb-8 flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold tracking-tight sm:text-3xl">Shop by Category</h2>
                <p className="mt-1 text-muted-foreground">Browse our curated collections</p>
              </div>
              <Link
                to="/categories"
                className="hidden items-center gap-1 text-sm font-medium text-primary hover:underline sm:flex"
              >
                View All <ChevronRight className="h-4 w-4" />
              </Link>
            </div>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6">
              {categories.slice(0, 6).map((category) => (
                <Link
                  key={category.id}
                  to={`/products?category=${category.slug}`}
                  className="group relative overflow-hidden rounded-2xl"
                >
                  <div className="aspect-[4/5] w-full overflow-hidden">
                    <img
                      src={category.image_url || getCategoryImage(category.slug)}
                      alt={category.name}
                      className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
                    />
                  </div>
                  <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent" />
                  <div className="absolute bottom-0 left-0 right-0 p-4">
                    <h3 className="text-lg font-bold text-white drop-shadow-lg">{category.name}</h3>
                    <span className="mt-1 inline-flex items-center gap-1 text-sm text-white/80 transition-colors group-hover:text-white">
                      Shop now <ArrowRight className="h-3.5 w-3.5 transition-transform group-hover:translate-x-1" />
                    </span>
                  </div>
                </Link>
              ))}
            </div>
          </section>
        )}

        {/* ── Trending Products ── */}
        <section>
          <div className="mb-8 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-orange-100 text-orange-600">
                <TrendingUp className="h-5 w-5" />
              </div>
              <div>
                <h2 className="text-2xl font-bold tracking-tight sm:text-3xl">Trending Now</h2>
                <p className="text-muted-foreground">Most popular picks this week</p>
              </div>
            </div>
            <Link
              to="/products?sort=popular"
              className="hidden items-center gap-1 text-sm font-medium text-primary hover:underline sm:flex"
            >
              View All <ChevronRight className="h-4 w-4" />
            </Link>
          </div>
          <ProductGrid products={trending ?? []} isLoading={trendingLoading} />
        </section>

        {/* ── Promo Banner ── */}
        <section className="relative overflow-hidden rounded-3xl bg-gradient-to-r from-amber-500 to-orange-600">
          <div className="absolute inset-0 bg-[url('https://images.unsplash.com/photo-1607083206869-4c7672e72a8a?w=1200&q=80')] bg-cover bg-center opacity-15" />
          <div className="relative flex flex-col items-center gap-6 px-8 py-14 text-center md:flex-row md:justify-between md:text-left">
            <div>
              <span className="mb-2 inline-flex items-center gap-2 rounded-full bg-white/20 px-3 py-1 text-sm font-medium text-white">
                <Zap className="h-4 w-4" />
                Limited Time Offer
              </span>
              <h2 className="mt-3 text-3xl font-extrabold text-white sm:text-4xl">
                Up to 50% Off Flash Sale
              </h2>
              <p className="mt-2 text-lg text-white/80">
                Don't miss out on incredible deals. Hurry, offer ends soon!
              </p>
            </div>
            <Button asChild size="lg" className="bg-white text-orange-600 hover:bg-orange-50 font-semibold shadow-xl shadow-orange-900/20 px-8 shrink-0">
              <Link to="/promotions">
                Shop the Sale <ArrowRight className="ml-2 h-5 w-5" />
              </Link>
            </Button>
          </div>
        </section>

        {/* ── Featured Products ── */}
        <section>
          <div className="mb-8 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-violet-100 text-violet-600">
                <Sparkles className="h-5 w-5" />
              </div>
              <div>
                <h2 className="text-2xl font-bold tracking-tight sm:text-3xl">Featured Products</h2>
                <p className="text-muted-foreground">Hand-picked just for you</p>
              </div>
            </div>
            <Link
              to="/products"
              className="hidden items-center gap-1 text-sm font-medium text-primary hover:underline sm:flex"
            >
              View All <ChevronRight className="h-4 w-4" />
            </Link>
          </div>
          <ProductGrid products={featured ?? []} isLoading={featuredLoading} />
        </section>
      </div>
    </div>
  );
}
