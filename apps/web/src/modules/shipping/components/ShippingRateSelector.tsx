import { Truck } from 'lucide-react';
import { formatPrice } from '@/shared/lib/utils';
import { cn } from '@/shared/lib/utils';
import type { ShippingRate } from '../services/shipping.api';

interface ShippingRateSelectorProps {
  rates: ShippingRate[];
  selectedId?: string;
  onSelect: (rate: ShippingRate) => void;
}

export function ShippingRateSelector({ rates, selectedId, onSelect }: ShippingRateSelectorProps) {
  return (
    <div className="space-y-2">
      {rates.map((rate) => (
        <button
          key={rate.id}
          onClick={() => onSelect(rate)}
          className={cn(
            'flex w-full items-center gap-4 rounded-lg border p-4 text-left transition-colors',
            selectedId === rate.id ? 'border-primary bg-primary/5' : 'hover:bg-muted'
          )}
        >
          <Truck className="h-5 w-5 text-muted-foreground" />
          <div className="flex-1">
            <p className="text-sm font-medium">
              {rate.carrier} - {rate.service}
            </p>
            <p className="text-xs text-muted-foreground">
              Estimated {rate.estimated_days} business day{rate.estimated_days !== 1 ? 's' : ''}
            </p>
          </div>
          <span className="font-semibold">{formatPrice(rate.price)}</span>
        </button>
      ))}
    </div>
  );
}
