import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/shared/components/ui/button';
import { Zap } from 'lucide-react';
import type { FlashSale } from '../services/flash-sale.api';

interface FlashSaleBannerProps {
  sale: FlashSale;
}

export function FlashSaleBanner({ sale }: FlashSaleBannerProps) {
  const [timeLeft, setTimeLeft] = useState('');

  useEffect(() => {
    const updateTimer = () => {
      const now = new Date().getTime();
      const end = new Date(sale.ends_at).getTime();
      const diff = end - now;

      if (diff <= 0) {
        setTimeLeft('Ended');
        return;
      }

      const hours = Math.floor(diff / (1000 * 60 * 60));
      const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
      const seconds = Math.floor((diff % (1000 * 60)) / 1000);
      setTimeLeft(`${hours}h ${minutes}m ${seconds}s`);
    };

    updateTimer();
    const interval = setInterval(updateTimer, 1000);
    return () => clearInterval(interval);
  }, [sale.ends_at]);

  return (
    <div className="rounded-lg bg-gradient-to-r from-orange-500 to-red-500 px-6 py-4 text-white">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Zap className="h-6 w-6" />
          <div>
            <h3 className="text-lg font-bold">{sale.title}</h3>
            <p className="text-sm text-white/80">{sale.description}</p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <div className="text-right">
            <p className="text-xs text-white/80">Ends in</p>
            <p className="font-mono text-lg font-bold">{timeLeft}</p>
          </div>
          <Button asChild variant="secondary">
            <Link to={`/flash-sales/${sale.id}`}>Shop Now</Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
