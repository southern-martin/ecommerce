import { Button } from '@/shared/components/ui/button';
import { Input } from '@/shared/components/ui/input';
import { Label } from '@/shared/components/ui/label';
import { Separator } from '@/shared/components/ui/separator';
import { Star, X } from 'lucide-react';
import type { FilterState, Category } from '../types/shop.types';

interface FilterPanelProps {
  filters: Partial<FilterState>;
  categories: Category[];
  onFilterChange: (filters: Partial<FilterState>) => void;
  onReset: () => void;
}

export function FilterPanel({
  filters,
  categories,
  onFilterChange,
  onReset,
}: FilterPanelProps) {
  return (
    <aside className="w-64 space-y-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold">Filters</h3>
        <Button variant="ghost" size="sm" onClick={onReset}>
          <X className="mr-1 h-4 w-4" />
          Reset
        </Button>
      </div>

      <Separator />

      <div className="space-y-3">
        <Label className="text-sm font-medium">Category</Label>
        <div className="space-y-2">
          {categories.map((cat) => (
            <button
              key={cat.id}
              onClick={() => onFilterChange({ category: cat.slug })}
              className={`block w-full text-left text-sm hover:text-primary ${
                filters.category === cat.slug
                  ? 'font-medium text-primary'
                  : 'text-muted-foreground'
              }`}
            >
              {cat.name}
            </button>
          ))}
        </div>
      </div>

      <Separator />

      <div className="space-y-3">
        <Label className="text-sm font-medium">Price Range</Label>
        <div className="flex items-center gap-2">
          <Input
            type="number"
            placeholder="Min"
            value={filters.min_price ?? ''}
            onChange={(e) =>
              onFilterChange({ min_price: e.target.value ? Number(e.target.value) : undefined })
            }
            className="h-9"
          />
          <span className="text-muted-foreground">-</span>
          <Input
            type="number"
            placeholder="Max"
            value={filters.max_price ?? ''}
            onChange={(e) =>
              onFilterChange({ max_price: e.target.value ? Number(e.target.value) : undefined })
            }
            className="h-9"
          />
        </div>
      </div>

      <Separator />

      <div className="space-y-3">
        <Label className="text-sm font-medium">Minimum Rating</Label>
        <div className="space-y-1">
          {[4, 3, 2, 1].map((rating) => (
            <button
              key={rating}
              onClick={() => onFilterChange({ rating })}
              className={`flex w-full items-center gap-1 rounded px-2 py-1 text-sm hover:bg-muted ${
                filters.rating === rating ? 'bg-muted font-medium' : ''
              }`}
            >
              {Array.from({ length: 5 }).map((_, i) => (
                <Star
                  key={i}
                  className={`h-3.5 w-3.5 ${
                    i < rating ? 'fill-yellow-400 text-yellow-400' : 'text-muted-foreground/30'
                  }`}
                />
              ))}
              <span className="ml-1">& up</span>
            </button>
          ))}
        </div>
      </div>

      <Separator />

      <div className="flex items-center gap-2">
        <input
          type="checkbox"
          id="in_stock"
          checked={filters.in_stock ?? false}
          onChange={(e) => onFilterChange({ in_stock: e.target.checked || undefined })}
          className="rounded border-input"
        />
        <Label htmlFor="in_stock" className="text-sm">
          In Stock Only
        </Label>
      </div>
    </aside>
  );
}
